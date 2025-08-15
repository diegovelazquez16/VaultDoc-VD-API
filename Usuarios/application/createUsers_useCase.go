
package application

import (
	"fmt"
	"strings"

	"VaultDoc-VD/Usuarios/domain/entities"
	"VaultDoc-VD/Usuarios/domain/repository"
	"VaultDoc-VD/Usuarios/domain/services"
)

type CreateUserUseCase struct {
	repo   repository.UserRepository
	bcrypt services.IBcryptService
}

func NewCreateUserUseCase(repo repository.UserRepository, bcrypt services.IBcryptService) *CreateUserUseCase {
	return &CreateUserUseCase{
		repo:   repo,
		bcrypt: bcrypt,
	}
}

func (uc *CreateUserUseCase) Execute(user entities.User) (*entities.User, error) {
	
	if err := uc.validateUser(user); err != nil {
		return nil, fmt.Errorf("validación fallida: %w", err)
	}

	existingUser, _ := uc.repo.FindByEmail(user.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("el email %s ya está registrado", user.Email)
	}

	hashedPassword, err := uc.bcrypt.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("error al procesar la contraseña: %w", err)
	}
	user.Password = hashedPassword

	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.Nombre = strings.TrimSpace(user.Nombre)
	user.Apellidos = strings.TrimSpace(user.Apellidos)
	user.Id_Rol = 2
	user.Departamento = "General"

	if err := uc.repo.Save(user); err != nil {
		return nil, fmt.Errorf("error al guardar usuario: %w", err)
	}

	
	createdUser, err := uc.repo.FindByEmail(user.Email)
	if err != nil {
		
		return &user, nil
	}

	createdUser.Password = ""
	
	return createdUser, nil
}

func (uc *CreateUserUseCase) validateUser(user entities.User) error {
	if strings.TrimSpace(user.Email) == "" {
		return fmt.Errorf("el email es requerido")
	}

	if !uc.isValidEmail(user.Email) {
		return fmt.Errorf("el formato del email no es válido")
	}

	if strings.TrimSpace(user.Password) == "" {
		return fmt.Errorf("la contraseña es requerida")
	}

	if len(user.Password) < 6 {
		return fmt.Errorf("la contraseña debe tener al menos 6 caracteres")
	}

	if strings.TrimSpace(user.Nombre) == "" {
		return fmt.Errorf("el nombre es requerido")
	}

	if strings.TrimSpace(user.Apellidos) == "" {
		return fmt.Errorf("los apellidos son requeridos")
	}

	if strings.TrimSpace(user.Departamento) == "" {
		return fmt.Errorf("el departamento es requerido")
	}

	return nil
}

func (uc *CreateUserUseCase) isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}