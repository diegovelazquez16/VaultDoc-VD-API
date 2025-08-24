// Usuarios/domain/services/token_manager.go
package services

type TokenManager interface {
	GenerateToken(userId int) (string, error)
	ValidateToken(token string) (bool, map[string]interface{}, error)
}