// Middlewares/Middleware_AdminBoss.go
package service

import (
	"net/http"
	"strings"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AdminBossMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// Usar el mismo método que AuthMiddleware
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer token not found"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Verificar que el rol esté presente
		roleIDInterface, exists := claims["roleId"]
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Rol no encontrado en el token"})
			return
		}

		// Convertir roleId a int (puede venir como float64 desde JWT)
		var roleID int
		switch v := roleIDInterface.(type) {
		case float64:
			roleID = int(v)
		case int:
			roleID = v
		case string:
			if id, err := strconv.Atoi(v); err == nil {
				roleID = id
			} else {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Formato de rol inválido"})
				return
			}
		default:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Tipo de rol no soportado"})
			return
		}

		if roleID != 2 && roleID != 3 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Acceso denegado. Solo los administradores y jefes de departamento pueden realizar esta acción",
			})
			return
		}

		// Convertir los claims a los tipos correctos antes de almacenarlos (igual que AuthMiddleware)
		c.Set("userID", claims["userId"])
		c.Set("email", claims["email"])
		c.Set("roleID", claims["roleId"])
		
		// Asegurar que department se almacene como string
		if dept, exists := claims["department"]; exists {
			c.Set("department", dept)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Department not found in token"})
			return
		}

		// Información adicional sobre el tipo de acceso
		if roleID == 3 {
			c.Set("accessType", "admin")
		} else if roleID == 2 {
			c.Set("accessType", "boss")
		}

		c.Next()
	}
}