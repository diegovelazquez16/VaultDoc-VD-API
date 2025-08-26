// Usuarios/application/get_profile_usecase.go
package application

import (
	"fmt"

	"VaultDoc-VD/Usuarios/domain/entities"
	"VaultDoc-VD/Usuarios/domain/repository"
)

type GetProfileUseCase struct {
	repo repository.UserRepository
}

func NewGetProfileUseCase(repo repository.UserRepository) *GetProfileUseCase {
	return &GetProfileUseCase{
		repo: repo,
	}
}

func (uc *GetProfileUseCase) Execute(authenticatedUserID int) (*entities.User, error) {
	// Validar que el ID sea válido
	if authenticatedUserID <= 0 {
		return nil, fmt.Errorf("ID de usuario inválido")
	}

	// Obtener el perfil del usuario
	user, err := uc.repo.GetProfile(authenticatedUserID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener perfil del usuario: %w", err)
	}

	// Asegurar que no se incluya información sensible
	user.Password = ""
	
	return user, nil
}