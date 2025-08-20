// Usuarios/domain/services/bcrypt.service.go
package services

type IBcryptService interface {
	HashPassword(password string) (string, error)
	ComparePasswords(hashedPassword string, providedPassword string) bool
}