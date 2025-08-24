// Archivos/infrastructure/controllers/grant_change_permission_controller.go
package controllers

import (
	"net/http"
	"strings"
	"VaultDoc-VD/Archivos/application"
	entities "VaultDoc-VD/Archivos/domain/entities"
	"github.com/gin-gonic/gin"
)

type GrantChangePermissionController struct {
	useCase *application.GrantChangePermissionUseCase
}

func NewGrantChangePermissionController(useCase *application.GrantChangePermissionUseCase) *GrantChangePermissionController {
	return &GrantChangePermissionController{useCase: useCase}
}

func (c *GrantChangePermissionController) Execute(ctx *gin.Context) {
	var input struct {
		Id_File int `json:"id_file" binding:"required"`
		Id_User int `json:"id_user" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Entrada de datos no válida",
			"error":   err.Error(),
		})
		return
	}

	// Validaciones básicas
	if input.Id_File <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID de archivo no válido",
			"error":   "El ID del archivo debe ser un número positivo",
		})
		return
	}

	if input.Id_User <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID de usuario no válido",
			"error":   "El ID del usuario debe ser un número positivo",
		})
		return
	}

	changeFile := entities.ChangeFile{
		Id_File: input.Id_File,
		Id_User: input.Id_User,
	}

	if err := c.useCase.Execute(changeFile); err != nil {
		statusCode := http.StatusInternalServerError
		// Determinar código de estado basado en el error
		if strings.Contains(err.Error(), "no existe") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(err.Error(), "ya existe") {
			statusCode = http.StatusConflict
		}

		ctx.JSON(statusCode, gin.H{
			"message": "Error al otorgar permiso de edición",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Permiso de edición otorgado exitosamente",
		"id_file": input.Id_File,
		"id_user": input.Id_User,
	})
}