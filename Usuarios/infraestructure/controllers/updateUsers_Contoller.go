// Usuarios/controllers/update_user_controller.go
package controllers

import (
	"VaultDoc-VD/Usuarios/application"
	"VaultDoc-VD/Usuarios/domain/entities"
	"VaultDoc-VD/validators"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UpdateUserController struct {
	useCase *application.UpdateUserUseCase
}

func NewUpdateUserController(useCase *application.UpdateUserUseCase) *UpdateUserController {
	return &UpdateUserController{useCase: useCase}
}

func (c *UpdateUserController) Execute(ctx *gin.Context) {
	// Validar y obtener ID del parámetro
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID de usuario inválido",
			"details": "El ID debe ser un número entero positivo",
		})
		return
	}

	// Bindear datos JSON
	var userUpdate entities.User
	if err := ctx.ShouldBindJSON(&userUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error al leer los datos del usuario",
			"details": err.Error(),
		})
		return
	}

	// Asignar ID del parámetro URL
	userUpdate.Id = id

	// Validaciones básicas
	if err := c.validateUpdateInput(userUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos de actualización inválidos",
			"details": err.Error(),
		})
		return
	}

	// Ejecutar caso de uso de actualización
	updatedUser, err := c.useCase.Execute(userUpdate)
	if err != nil {
		if strings.Contains(err.Error(), "no existe") ||
			strings.Contains(err.Error(), "no encontrado") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Usuario no encontrado",
				"details": err.Error(),
			})
			return
		}

		if strings.Contains(err.Error(), "ya está siendo usado") ||
			strings.Contains(err.Error(), "email duplicado") {
			ctx.JSON(http.StatusConflict, gin.H{
				"error":   "Email ya está en uso",
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
			"error":   "Error interno al actualizar usuario",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Usuario actualizado correctamente",
		"user": gin.H{
			"id":           updatedUser.Id,
			"email":        updatedUser.Email,
			"nombre":       updatedUser.Nombre,
			"apellidos":    updatedUser.Apellidos,
			"id_rol":       updatedUser.Id_Rol,
			"departamento": updatedUser.Departamento,
		},
		"sync_status": "local_updated",
	})
}

func (c *UpdateUserController) validateUpdateInput(user entities.User) error {
	if strings.TrimSpace(user.Email) == "" {
		return fmt.Errorf("el email es requerido")
	}

	if strings.TrimSpace(user.Nombre) == "" {
		return fmt.Errorf("el nombre es requerido")
	}

	if strings.TrimSpace(user.Apellidos) == "" {
		return fmt.Errorf("los apellidos son requeridos")
	}

	if user.Password != "" && len(user.Password) < 6 {
		return fmt.Errorf("la contraseña debe tener al menos 6 caracteres")
	}

	if user.Departamento != "" && !validators.IsValidDepartamento(user.Departamento) {
		return fmt.Errorf("el departamento debe ser uno de los ya existentes")
	}

	if user.Id_Rol != 0 && user.Id_Rol < 1 {
		return fmt.Errorf("el id_rol debe ser un número positivo")
	}

	return nil
}