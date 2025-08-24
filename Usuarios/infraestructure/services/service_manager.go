// Usuarios/infraestructure/services/service_manager.go
package service

import (
	"os"

	"VaultDoc-VD/Usuarios/domain/services"
	adapters "VaultDoc-VD/Usuarios/infraestructure/adapters"
)

// Inicializar el servicio de BCrypt
func InitBcryptService() services.IBcryptService {
	return adapters.NewBcrypt()
}

// Inicializar el Token Manager
func InitTokenManager() services.TokenManager {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET no está configurado en las variables de entorno")
	}
	return &adapters.JWTManager{SecretKey: jwtSecret}
}
