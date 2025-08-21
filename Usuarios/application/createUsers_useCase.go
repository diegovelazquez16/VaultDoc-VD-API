// Usuarios/application/createUsers_useCase.go
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

	// Limpiar y normalizar datos
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.Nombre = strings.TrimSpace(user.Nombre)
	user.Apellidos = strings.TrimSpace(user.Apellidos)
	user.Departamento = strings.TrimSpace(user.Departamento)

	// Validar que id_rol sea válido si no se proporciona, usar valor por defecto
	if user.Id_Rol <= 0 {
		user.Id_Rol = 2 // Valor por defecto
	}

	// Si departamento está vacío, usar valor por defecto
	if user.Departamento == "" {
		user.Departamento = "Gerencia Operativa" // Valor por defecto del enum
	}

	if err := uc.repo.Save(user); err != nil {
		return nil, fmt.Errorf("error al guardar usuario: %w", err)
	}

	// Obtener el usuario creado
	createdUser, err := uc.repo.FindByEmail(user.Email)
	if err != nil {
		// Si no se puede recuperar, retornar el usuario sin ID
		return &user, nil
	}

	// Limpiar password antes de retornar
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

	// Validar departamento si se proporciona
	if user.Departamento != "" {
		if !uc.isValidDepartamento(user.Departamento) {
			return fmt.Errorf("el departamento debe ser: Finanzaz, Gerencia Operativa o General")
		}
	}

	// Validar id_rol si se proporciona
	if user.Id_Rol != 0 && user.Id_Rol < 1 {
		return fmt.Errorf("el id_rol debe ser un número positivo")
	}

	return nil
}

func (uc *CreateUserUseCase) isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func (uc *CreateUserUseCase) isValidDepartamento(departamento string) bool {
	validDepartamentos := []string{"Finanzaz", "Gerencia Operativa", "General"}
	for _, validDept := range validDepartamentos {
		if strings.EqualFold(departamento, validDept) {
			return true
		}
	}
	return false
}