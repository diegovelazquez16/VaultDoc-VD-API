package application

import (
	"VaultDoc-VD/Usuarios/domain/entities"
	"VaultDoc-VD/Usuarios/domain/repository"
	"VaultDoc-VD/Usuarios/domain/services"
	"fmt"
)

type LoginUseCase struct {
	db     repository.UserRepository
	jwt    services.TokenManager
	bcrypt services.IBcryptService
}

func NewLoginUseCase(db repository.UserRepository, jwt services.TokenManager, bcrypt services.IBcryptService) *LoginUseCase {
	return &LoginUseCase{
		db:     db,
		jwt:    jwt,
		bcrypt: bcrypt,
	}
}

func (lu *LoginUseCase) Execute(email string, password string) (*entities.User, string, error) {
	user, err := lu.db.FindByEmail(email)
	if err != nil {
		return nil, "", fmt.Errorf("credenciales inválidas")
	}

	if !lu.bcrypt.ComparePasswords(user.Password, password) {
		return nil, "", fmt.Errorf("credenciales inválidas")
	}

	token, err := lu.jwt.GenerateToken(user.Id)
	if err != nil {
		return nil, "", fmt.Errorf("error generando token: %w", err)
	}

	return user, token, nil
}
