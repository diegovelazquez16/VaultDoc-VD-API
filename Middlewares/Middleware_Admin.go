// Middlewares/Middleware_Admin.go
package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AdminMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Primero ejecutar la validación del token (similar al AuthMiddleware)
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
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

		// Verificar que el rol sea admin (ID 3)
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

		// Verificar que sea el rol admin (ID 3)
		if roleID != 3 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Acceso denegado. Solo los administradores pueden realizar esta acción",
			})
			return
		}

		// Guardar información del usuario en el contexto
		c.Set("userID", claims["userId"])
		c.Set("roleID", roleID)
		c.Set("email", claims["email"])
		c.Set("department", claims["department"])

		c.Next()
	}
}