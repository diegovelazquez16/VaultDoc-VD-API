// Usuarios/controllers/get_profile_controller.go
package controllers

import (
	"VaultDoc-VD/Usuarios/application"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type GetProfileController struct {
	useCase *application.GetProfileUseCase
}

func NewGetProfileController(useCase *application.GetProfileUseCase) *GetProfileController {
	return &GetProfileController{useCase: useCase}
}

func (c *GetProfileController) Execute(ctx *gin.Context) {
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

	userProfile, err := c.useCase.Execute(authenticatedUserID)
	if err != nil {
		if strings.Contains(err.Error(), "no encontrado") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Perfil no encontrado",
				"details": err.Error(),
			})
			return
		}

		if strings.Contains(err.Error(), "inválido") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Solicitud inválida",
				"details": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error interno al obtener perfil",
			"details": err.Error(),
		})
		return
	}

	roleName := c.getRoleName(userProfile.Id_Rol)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Perfil obtenido exitosamente",
		"user": gin.H{
			"id":           userProfile.Id,
			"email":        userProfile.Email,
			"nombre":       userProfile.Nombre,
			"apellidos":    userProfile.Apellidos,
			"id_rol":       userProfile.Id_Rol,
			"rol_nombre":   roleName,
			"departamento": userProfile.Departamento,
		},
		"editable_fields": []string{"nombre", "apellidos", "email", "password"},
		"readonly_fields": []string{"id", "id_rol", "departamento"},
	})
}

func (c *GetProfileController) getRoleName(roleID int) string {
	switch roleID {
	case 1:
		return "usuario"
	case 2:
		return "empleado"
	case 3:
		return "admin"
	case 4:
		return "jefe_departamento"
	case 5:
		return "supervisor"
	default:
		return "desconocido"
	}
}