// Usuarios/controllers/update_myProfile_Controller.go
package controllers

import (
	"VaultDoc-VD/Usuarios/application"
	"VaultDoc-VD/Usuarios/domain/entities"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UpdateProfileController struct {
	useCase *application.UpdateProfileUseCase
}

func NewUpdateProfileController(useCase *application.UpdateProfileUseCase) *UpdateProfileController {
	return &UpdateProfileController{useCase: useCase}
}

func (c *UpdateProfileController) Execute(ctx *gin.Context) {
	// Obtener el ID del usuario autenticado desde el token JWT
	userIDInterface, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Usuario no autenticado",
		})
		return
	}

	// Convertir userID a int (puede venir como float64 desde JWT)
	var authenticatedUserID int
	switch v := userIDInterface.(type) {
	case float64:
		authenticatedUserID = int(v)
	case int:
		authenticatedUserID = v
	case string:
		if id, err := strconv.Atoi(v); err == nil {
			authenticatedUserID = id
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "ID de usuario inválido en token",
			})
			return
		}
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Tipo de ID de usuario no soportado",
		})
		return
	}

	// Bindear datos JSON
	var userUpdate entities.User
	if err := ctx.ShouldBindJSON(&userUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error al leer los datos del perfil",
			"details": err.Error(),
		})
		return
	}

	// El ID del usuario a actualizar debe ser el mismo que el autenticado
	userUpdate.Id = authenticatedUserID

	// Validaciones básicas
	if err := c.validateProfileInput(userUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos de perfil inválidos",
			"details": err.Error(),
		})
		return
	}

	// Ejecutar caso de uso de actualización de perfil
	updatedUser, err := c.useCase.Execute(userUpdate, authenticatedUserID)
	if err != nil {
		if strings.Contains(err.Error(), "no tienes permisos") {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error":   "Acceso denegado",
				"details": err.Error(),
			})
			return
		}

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
			"error":   "Error interno al actualizar perfil",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Perfil actualizado correctamente",
		"user": gin.H{
			"id":           updatedUser.Id,
			"email":        updatedUser.Email,
			"nombre":       updatedUser.Nombre,
			"apellidos":    updatedUser.Apellidos,
			"id_rol":       updatedUser.Id_Rol,
			"departamento": updatedUser.Departamento,
		},
		"sync_status": "profile_updated",
	})
}

func (c *UpdateProfileController) validateProfileInput(user entities.User) error {
	if strings.TrimSpace(user.Email) == "" {
		return fmt.Errorf("el email es requerido")
	}

	if strings.TrimSpace(user.Nombre) == "" {
		return fmt.Errorf("el nombre es requerido")
	}

	if strings.TrimSpace(user.Apellidos) == "" {
		return fmt.Errorf("los apellidos son requeridos")
	}

	// Validar contraseña solo si se proporciona
	if user.Password != "" && len(user.Password) < 6 {
		return fmt.Errorf("la contraseña debe tener al menos 6 caracteres")
	}

	return nil
}