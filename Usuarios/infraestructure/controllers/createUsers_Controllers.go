// Usuarios/controllers/create_user_controller.go
package controllers

import (
	"VaultDoc-VD/Usuarios/application"
	"VaultDoc-VD/Usuarios/domain/entities"
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

	// Respuesta exitosa con información del usuario creado
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

	// Validar departamento si se proporciona
	if user.Departamento != "" && !c.isValidDepartamento(user.Departamento) {
		return fmt.Errorf("el departamento debe ser: Finanzaz, Operativo o General")
	}

	// Validar id_rol si se proporciona
	if user.Id_Rol != 0 && user.Id_Rol < 1 {
		return fmt.Errorf("el id_rol debe ser un número positivo")
	}

	return nil
}

func (c *CreateUserController) isValidDepartamento(departamento string) bool {
	validDepartamentos := []string{"Finanzaz", "Operativo", "General"}
	for _, validDept := range validDepartamentos {
		if strings.EqualFold(departamento, validDept) {
			return true
		}
	}
	return false
}