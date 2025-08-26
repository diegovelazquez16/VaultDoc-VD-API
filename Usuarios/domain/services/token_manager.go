// Usuarios/domain/services/token_manager.go
package services

import "VaultDoc-VD/Usuarios/domain/entities"

type TokenManager interface {
	GenerateToken(user *entities.User) (string, error)  // Cambiar parámetro
	ValidateToken(token string) (bool, map[string]interface{}, error)
}