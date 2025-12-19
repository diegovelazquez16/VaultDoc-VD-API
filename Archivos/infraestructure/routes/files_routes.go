// Archivos/infrastructure/routes/files_routes.go
package routes

import (
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"VaultDoc-VD/Archivos/infraestructure/controllers"
	"VaultDoc-VD/Middlewares"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// limiterEntry almacena un rate limiter con su timestamp de último uso
type limiterEntry struct {
	limiter  *rate.Limiter
	lastSeen atomic.Value // thread-safe timestamp
}

// RateLimiter gestiona los limitadores por IP/usuario con expiración automática
type RateLimiter struct {
	ips             map[string]*limiterEntry
	mu              sync.RWMutex
	r               rate.Limit
	b               int
	ttl             time.Duration
	cleanupInterval time.Duration
	globalLimiter   *rate.Limiter // Mejora #5: Rate limiter global
}

// RateLimiterConfig contiene la configuración del rate limiter
type RateLimiterConfig struct {
	RequestsPerSecond float64
	Burst             int
	TTL               time.Duration
	CleanupInterval   time.Duration
	GlobalRPS         float64 // Mejora #5: Límite global de requests
	GlobalBurst       int
}

// NewRateLimiter crea un nuevo rate limiter con limpieza automática y límite global
func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	rl := &RateLimiter{
		ips:             make(map[string]*limiterEntry),
		r:               rate.Limit(config.RequestsPerSecond),
		b:               config.Burst,
		ttl:             config.TTL,
		cleanupInterval: config.CleanupInterval,
		globalLimiter:   rate.NewLimiter(rate.Limit(config.GlobalRPS), config.GlobalBurst),
	}

	// Iniciar goroutine de limpieza
	go rl.cleanupExpiredEntries()

	return rl
}

// GetLimiter obtiene o crea un limiter para una IP/usuario
func (rl *RateLimiter) GetLimiter(key string) *rate.Limiter {
	now := time.Now()

	// Intento 1: lectura rápida con RLock
	rl.mu.RLock()
	entry, exists := rl.ips[key]
	if exists {
		entry.lastSeen.Store(now)
		rl.mu.RUnlock()
		return entry.limiter
	}
	rl.mu.RUnlock()

	// Intento 2: no existe, necesitamos crear (Lock completo)
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Double-check: otra goroutine pudo haberlo creado
	entry, exists = rl.ips[key]
	if exists {
		entry.lastSeen.Store(now)
		return entry.limiter
	}

	// Crear nueva entrada
	limiter := rate.NewLimiter(rl.r, rl.b)
	newEntry := &limiterEntry{
		limiter: limiter,
	}
	newEntry.lastSeen.Store(now)
	rl.ips[key] = newEntry

	return limiter
}

// cleanupExpiredEntries elimina periódicamente las entradas antiguas
func (rl *RateLimiter) cleanupExpiredEntries() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, entry := range rl.ips {
			lastSeen := entry.lastSeen.Load().(time.Time)
			if now.Sub(lastSeen) > rl.ttl {
				delete(rl.ips, key)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimitMiddleware devuelve el middleware de rate limiting con mejoras
func (rl *RateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rl.globalLimiter.Allow() {
			setRateLimitHeaders(c, 0, rl.globalLimiter)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Servicio temporalmente saturado. Intenta más tarde.",
				"message":     "Global rate limit exceeded",
				"retry_after": 60,
			})
			c.Abort()
			return
		}

		// Obtener clave de rate limiting (IP + UserID si está autenticado)
		key := getRateLimitKey(c)
		limiter := rl.GetLimiter(key)

		// Verificar límite individual
		if !limiter.Allow() {
			setRateLimitHeaders(c, 0, limiter)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Demasiadas peticiones. Intenta más tarde.",
				"message":     "Rate limit exceeded",
				"retry_after": 60,
			})
			c.Abort()
			return
		}

		remaining := limiter.Tokens()
		setRateLimitHeaders(c, int(remaining), limiter)

		c.Next()
	}
}

// getRateLimitKey obtiene la clave de rate limiting (IP o UserID si está autenticado)
func getRateLimitKey(c *gin.Context) string {
	// Priorizar UserID si está autenticado (más preciso que IP)
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(string); ok && uid != "" {
			return "user:" + uid
		}
	}

	// Fallback a IP
	return "ip:" + c.ClientIP()
}

func setRateLimitHeaders(c *gin.Context, remaining int, limiter *rate.Limiter) {
	c.Header("X-RateLimit-Limit", strconv.Itoa(limiter.Burst()))
	c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
	c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Minute).Unix(), 10))
	
	if remaining == 0 {
		c.Header("Retry-After", "60")
	}
}

func getEnvFloat(key string, defaultVal float64) float64 {
	if val := os.Getenv(key); val != "" {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultVal
}

func SetupFilesRoutes(
	r *gin.Engine,
	createFileController *controllers.CreateFileController,
	getFileByIdController *controllers.GetFileByIdController,
	getAllFilesController *controllers.GetAllFilesController,
	getFilesByFolderController *controllers.GetFilesByFolderController,
	updateFileController *controllers.UpdateFileController,
	deleteFileController *controllers.DeleteFileController,
	downloadFileController *controllers.DownloadFileController,
	grantChangePermissionController *controllers.GrantChangePermissionController,
	removeChangePermissionController *controllers.RemoveChangePermissionController,
	grantViewPermissionController *controllers.GrantViewPermissionController,
	removeViewPermissionController *controllers.RemoveViewPermissionController,
	checkPermissionsController *controllers.CheckPermissionsController,
	searchFileController *controllers.SearchFileController,
	getUsersViewPermissionsControlleer *controllers.GetUsersViewPermissionsController,
	getUsersChangePermissionsController *controllers.GetUsersChangePermissionsController,
	getChangePermissionsOfAFolder *controllers.GetChangePermissionsOfAFolderController,
) {
	// Mejora #1: Configurar trusted proxies (CRÍTICO para seguridad)
	// Ajustar según tu infraestructura:
	// - Nginx local: []string{"127.0.0.1", "::1"}
	// - AWS ALB: rangos CIDR de AWS
	// - Cloudflare: IPs de Cloudflare
	trustedProxies := os.Getenv("TRUSTED_PROXIES")
	if trustedProxies != "" {
		// Ejemplo: "127.0.0.1,10.0.0.0/8"
		r.SetTrustedProxies(parseProxies(trustedProxies))
	} else {
		// Default para desarrollo local
		r.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	
	writeLimiter := NewRateLimiter(RateLimiterConfig{
		RequestsPerSecond: getEnvFloat("FILE_WRITE_RATE_LIMIT", 2),    
		Burst:             getEnvInt("FILE_WRITE_BURST_LIMIT", 5),     
		TTL:               getEnvDuration("RATE_LIMIT_TTL", 15*time.Minute),
		CleanupInterval:   getEnvDuration("RATE_LIMIT_CLEANUP", 5*time.Minute),
		GlobalRPS:         getEnvFloat("FILE_WRITE_GLOBAL_LIMIT", 50), 
		GlobalBurst:       getEnvInt("FILE_WRITE_GLOBAL_BURST", 100),
	})

	readLimiter := NewRateLimiter(RateLimiterConfig{
		RequestsPerSecond: getEnvFloat("FILE_READ_RATE_LIMIT", 10),
		Burst:             getEnvInt("FILE_READ_BURST_LIMIT", 20),
		TTL:               getEnvDuration("RATE_LIMIT_TTL", 15*time.Minute),
		CleanupInterval:   getEnvDuration("RATE_LIMIT_CLEANUP", 5*time.Minute),
		GlobalRPS:         getEnvFloat("FILE_READ_GLOBAL_LIMIT", 200), 
		GlobalBurst:       getEnvInt("FILE_READ_GLOBAL_BURST", 400),
	})

	downloadLimiter := NewRateLimiter(RateLimiterConfig{
		RequestsPerSecond: getEnvFloat("FILE_DOWNLOAD_RATE_LIMIT", 1),
		Burst:             getEnvInt("FILE_DOWNLOAD_BURST_LIMIT", 3),  
		TTL:               getEnvDuration("RATE_LIMIT_TTL", 15*time.Minute),
		CleanupInterval:   getEnvDuration("RATE_LIMIT_CLEANUP", 5*time.Minute),
		GlobalRPS:         getEnvFloat("FILE_DOWNLOAD_GLOBAL_LIMIT", 30), 
		GlobalBurst:       getEnvInt("FILE_DOWNLOAD_GLOBAL_BURST", 60),
	})

	permissionLimiter := NewRateLimiter(RateLimiterConfig{
		RequestsPerSecond: getEnvFloat("FILE_PERMISSION_RATE_LIMIT", 5),
		Burst:             getEnvInt("FILE_PERMISSION_BURST_LIMIT", 10),
		TTL:               getEnvDuration("RATE_LIMIT_TTL", 15*time.Minute),
		CleanupInterval:   getEnvDuration("RATE_LIMIT_CLEANUP", 5*time.Minute),
		GlobalRPS:         getEnvFloat("FILE_PERMISSION_GLOBAL_LIMIT", 100),
		GlobalBurst:       getEnvInt("FILE_PERMISSION_GLOBAL_BURST", 200),
	})

	filesGroup := r.Group("files")
	{
		writeGroup := filesGroup.Group("")
		writeGroup.Use(writeLimiter.RateLimitMiddleware())
		{
			writeGroup.POST("/", service.AuthMiddleware(jwtSecret), createFileController.Execute)
			writeGroup.PUT("/:id/:id_user", service.AuthMiddleware(jwtSecret), updateFileController.Execute)
			writeGroup.DELETE("/:id/:id_user", service.AuthMiddleware(jwtSecret), deleteFileController.Execute)
		}

		readGroup := filesGroup.Group("")
		readGroup.Use(readLimiter.RateLimitMiddleware())
		{
			readGroup.GET("/:id", service.AuthMiddleware(jwtSecret), getFileByIdController.Execute)
			readGroup.GET("/folder/:folderId", service.AuthMiddleware(jwtSecret), getFilesByFolderController.Execute)
			readGroup.GET("/search/:filename", service.AuthMiddleware(jwtSecret), searchFileController.Execute)
			readGroup.GET("/permissions/:fileId/:userId", service.AuthMiddleware(jwtSecret), checkPermissionsController.Execute)
			readGroup.GET("/permissions/change/:id_folder/:id_user", service.AuthMiddleware(jwtSecret), getChangePermissionsOfAFolder.Execute)
			readGroup.GET("/", service.AdminMiddleware(jwtSecret), getAllFilesController.Execute)
		}

		downloadGroup := filesGroup.Group("")
		downloadGroup.Use(downloadLimiter.RateLimitMiddleware())
		{
			downloadGroup.GET("/download/:id/:id_user", service.AuthMiddleware(jwtSecret), downloadFileController.Execute)
		}

		permissionGroup := filesGroup.Group("/permissions")
		permissionGroup.Use(permissionLimiter.RateLimitMiddleware())
		{
			permissionGroup.GET("/view/g/:file_id", service.BossMiddleware(jwtSecret), getUsersViewPermissionsControlleer.Execute)
			permissionGroup.GET("/change/g/:file_id", service.AdminBossMiddleware(jwtSecret), getUsersChangePermissionsController.Execute)

			permissionGroup.POST("/change/:id_user", service.BossMiddleware(jwtSecret), grantChangePermissionController.Execute)
			permissionGroup.DELETE("/change/:id_user/:id_file", service.BossMiddleware(jwtSecret), removeChangePermissionController.Execute)

			permissionGroup.POST("/view/:id_user", service.BossMiddleware(jwtSecret), grantViewPermissionController.Execute)
			permissionGroup.DELETE("/view/:id_user/:id_file", service.BossMiddleware(jwtSecret), removeViewPermissionController.Execute)
		}
	}
}

func parseProxies(proxies string) []string {
	if proxies == "" {
		return nil
	}
	
	result := []string{}
	for _, proxy := range splitAndTrim(proxies, ",") {
		if proxy != "" {
			result = append(result, proxy)
		}
	}
	return result
}

func splitAndTrim(s, sep string) []string {
	parts := []string{}
	for _, part := range splitString(s, sep) {
		trimmed := trimSpace(part)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

func splitString(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	result := []string{}
	current := ""
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, current)
			current = ""
			i += len(sep) - 1
		} else {
			current += string(s[i])
		}
	}
	result = append(result, current)
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}