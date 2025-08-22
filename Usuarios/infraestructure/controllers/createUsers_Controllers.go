// Usuarios/controllers/create_user_controller.go
package controllers

import (
	"VaultDoc-VD/Usuarios/application"
	"VaultDoc-VD/Usuarios/domain/entities"
	"VaultDoc-VD/validators"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type CreateUserController struct {
	useCase *application.CreateUserUseCase
}

func NewCreateUserController(useCase *application.CreateUserUseCase) *CreateUserController {
	return &CreateUserController{useCase: useCase}
}

func (c *CreateUserController) Execute(ctx *gin.Context) {
	// Verificar que el usuario autenticado sea admin (doble verificación)
	roleID, exists := ctx.Get("roleID")
	if !exists || roleID != 3 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Acceso denegado. Solo los administradores pueden crear usuarios",
		})
		return
	}

	// Obtener información del admin que está creando el usuario
	adminEmail, _ := ctx.Get("email")

	var user entities.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error al leer los datos del usuario",
			"details": err.Error(),
		})
		return
	}

	// Validaciones básicas
	if err := c.validateUserInput(user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos de usuario inválidos",
			"details": err.Error(),
		})
		return
	}

	// Si no se especifica rol, asignar rol de usuario regular (por ejemplo, ID 2)
	if user.Id_Rol == 0 {
		user.Id_Rol = 2 // O el ID que corresponda a "usuario regular"
	}

	// Ejecutar caso de uso
	createdUser, err := c.useCase.Execute(user)
	if err != nil {
		// Clasificar errores para respuestas más específicas
		if strings.Contains(err.Error(), "ya está registrado") ||
			strings.Contains(err.Error(), "ya existe") {
			ctx.JSON(http.StatusConflict, gin.H{
				"error":   "Email ya registrado",
				"details": err.Error(),
			})
			return
		}

		if strings.Contains(err.Error(), "validación") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Error de validación",
				"details": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error interno al crear usuario",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Usuario creado exitosamente",
		"user": gin.H{
			"id":           createdUser.Id,
			"email":        createdUser.Email,
			"nombre":       createdUser.Nombre,
			"apellidos":    createdUser.Apellidos,
			"id_rol":       createdUser.Id_Rol,
			"departamento": createdUser.Departamento,
		},
		"created_by":  adminEmail,
		"sync_status": "local_saved",
	})
}

func (c *CreateUserController) validateUserInput(user entities.User) error {
	if strings.TrimSpace(user.Email) == "" {
		return fmt.Errorf("el email es requerido")
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

	if user.Departamento != "" && !validators.IsValidDepartamento(user.Departamento) {
		return fmt.Errorf("el departamento debe ser uno de los ya existentes")
	}

	// Validar que el rol sea válido
	if user.Id_Rol != 0 {
		if user.Id_Rol < 1 {
			return fmt.Errorf("el id_rol debe ser un número positivo")
		}
		// Opcional: restringir la creación de otros admins
		if user.Id_Rol == 3 {
			return fmt.Errorf("no se puede asignar rol de administrador directamente")
		}
	}

	return nil
}