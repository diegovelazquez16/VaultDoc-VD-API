// Usuarios/application/update_myProfile_usecase.go
package application

import (
	"fmt"
	"strings"

	"VaultDoc-VD/Usuarios/domain/entities"
	"VaultDoc-VD/Usuarios/domain/repository"
	"VaultDoc-VD/Usuarios/domain/services"
)

type UpdateProfileUseCase struct {
	repo   repository.UserRepository
	bcrypt services.IBcryptService
}

func NewUpdateProfileUseCase(repo repository.UserRepository, bcrypt services.IBcryptService) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{
		repo:   repo,
		bcrypt: bcrypt,
	}
}

func (uc *UpdateProfileUseCase) Execute(userUpdate entities.User, authenticatedUserID int) (*entities.User, error) {
	// Verificar que el usuario solo pueda actualizar su propio perfil
	if userUpdate.Id != authenticatedUserID {
		return nil, fmt.Errorf("no tienes permisos para actualizar el perfil de otro usuario")
	}

	// Verificar que el usuario existe
	existingUser, err := uc.repo.FindById(userUpdate.Id)
	if err != nil {
		return nil, fmt.Errorf("usuario con ID %d no encontrado: %w", userUpdate.Id, err)
	}

	// Validar datos de entrada
	if err := uc.validateProfileData(userUpdate); err != nil {
		return nil, fmt.Errorf("validación fallida: %w", err)
	}

	// Verificar que el email no esté siendo usado por otro usuario (si se está cambiando)
	if userUpdate.Email != existingUser.Email {
		userWithEmail, _ := uc.repo.FindByEmail(userUpdate.Email)
		if userWithEmail != nil && userWithEmail.Id != userUpdate.Id {
			return nil, fmt.Errorf("el email %s ya está siendo usado por otro usuario", userUpdate.Email)
		}
	}

	// Preparar datos para actualización (solo campos permitidos)
	updatedUser := entities.User{
		Id:        userUpdate.Id,
		Nombre:    strings.TrimSpace(userUpdate.Nombre),
		Apellidos: strings.TrimSpace(userUpdate.Apellidos),
		Email:     strings.ToLower(strings.TrimSpace(userUpdate.Email)),
		Password:  existingUser.Password, // Mantener la contraseña actual por defecto
	}

	// Actualizar contraseña solo si se proporciona una nueva
	if userUpdate.Password != "" {
		hashedPassword, err := uc.bcrypt.HashPassword(userUpdate.Password)
		if err != nil {
			return nil, fmt.Errorf("error al procesar la nueva contraseña: %w", err)
		}
		updatedUser.Password = hashedPassword
	}

	// Actualizar usando el método UpdateProfile que solo actualiza campos básicos
	if err := uc.repo.UpdateProfile(updatedUser); err != nil {
		return nil, fmt.Errorf("error al actualizar perfil: %w", err)
	}

	// Obtener el usuario actualizado
	finalUser, err := uc.repo.FindById(updatedUser.Id)
	if err != nil {
		// Si no se puede recuperar, usar el usuario actualizado pero con los datos que no cambian
		finalUser = &entities.User{
			Id:           updatedUser.Id,
			Nombre:       updatedUser.Nombre,
			Apellidos:    updatedUser.Apellidos,
			Email:        updatedUser.Email,
			Id_Rol:       existingUser.Id_Rol,       // Mantener rol original
			Departamento: existingUser.Departamento, // Mantener departamento original
		}
	}

	// Limpiar password antes de retornar
	finalUser.Password = ""
	return finalUser, nil
}

func (uc *UpdateProfileUseCase) validateProfileData(user entities.User) error {
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

	// Validar contraseña solo si se proporciona
	if user.Password != "" && len(user.Password) < 6 {
		return fmt.Errorf("la nueva contraseña debe tener al menos 6 caracteres")
	}

	return nil
}

func (uc *UpdateProfileUseCase) isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}