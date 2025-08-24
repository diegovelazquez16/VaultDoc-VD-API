// Usuarios/application/update_user_usecase.go
package application

import (
	"fmt"
	"strings"

	"VaultDoc-VD/Usuarios/domain/entities"
	"VaultDoc-VD/Usuarios/domain/repository"
	"VaultDoc-VD/Usuarios/domain/services"
)

type UpdateUserUseCase struct {
	repo   repository.UserRepository
	bcrypt services.IBcryptService
}

func NewUpdateUserUseCase(repo repository.UserRepository, bcrypt services.IBcryptService) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		repo:   repo,
		bcrypt: bcrypt,
	}
}

func (uc *UpdateUserUseCase) Execute(userUpdate entities.User) (*entities.User, error) {

	existingUser, err := uc.repo.FindById(userUpdate.Id)
	if err != nil {
		return nil, fmt.Errorf("usuario con ID %d no encontrado: %w", userUpdate.Id, err)
	}

	if err := uc.validateUpdateData(userUpdate); err != nil {
		return nil, fmt.Errorf("validación fallida: %w", err)
	}

	// Verificar que el email no esté siendo usado por otro usuario
	if userUpdate.Email != existingUser.Email {
		userWithEmail, _ := uc.repo.FindByEmail(userUpdate.Email)
		if userWithEmail != nil && userWithEmail.Id != userUpdate.Id {
			return nil, fmt.Errorf("el email %s ya está siendo usado por otro usuario", userUpdate.Email)
		}
	}

	updatedUser := *existingUser
	updatedUser.Nombre = strings.TrimSpace(userUpdate.Nombre)
	updatedUser.Apellidos = strings.TrimSpace(userUpdate.Apellidos)
	updatedUser.Email = strings.ToLower(strings.TrimSpace(userUpdate.Email))

	if userUpdate.Departamento != "" {
		updatedUser.Departamento = strings.TrimSpace(userUpdate.Departamento)
	}

	if userUpdate.Id_Rol > 0 {
		updatedUser.Id_Rol = userUpdate.Id_Rol
	}

	// Actualizar contraseña solo si se proporciona una nueva
	if userUpdate.Password != "" {
		hashedPassword, err := uc.bcrypt.HashPassword(userUpdate.Password)
		if err != nil {
			return nil, fmt.Errorf("error al procesar la nueva contraseña: %w", err)
		}
		updatedUser.Password = hashedPassword
	}

	if err := uc.repo.Update(updatedUser); err != nil {
		return nil, fmt.Errorf("error al actualizar usuario: %w", err)
	}

	// Obtener el usuario actualizado
	finalUser, err := uc.repo.FindById(updatedUser.Id)
	if err != nil {
		// Si no se puede recuperar, usar el usuario actualizado
		finalUser = &updatedUser
	}

	// Limpiar password antes de retornar
	finalUser.Password = ""
	return finalUser, nil
}

func (uc *UpdateUserUseCase) validateUpdateData(user entities.User) error {
	if user.Id <= 0 {
		return fmt.Errorf("ID de usuario inválido")
	}

	if strings.TrimSpace(user.Email) == "" {
		return fmt.Errorf("el email es requerido")
	}

	if !uc.isValidEmail(user.Email) {
		return fmt.Errorf("el formato del email no es válido")
	}

	if strings.TrimSpace(user.Nombre) == "" {
		return fmt.Errorf("el nombre es requerido")
	}

	if strings.TrimSpace(user.Apellidos) == "" {
		return fmt.Errorf("los apellidos son requeridos")
	}

	if user.Password != "" && len(user.Password) < 6 {
		return fmt.Errorf("la nueva contraseña debe tener al menos 6 caracteres")
	}

	if user.Id_Rol != 0 && user.Id_Rol < 1 {
		return fmt.Errorf("el id_rol debe ser un número positivo")
	}

	return nil
}

func (uc *UpdateUserUseCase) isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}