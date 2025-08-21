// Archivos/infrastructure/controllers/grant_view_permission_controller.go
package controllers

import (
	"net/http"
	"VaultDoc-VD/Archivos/application"
	entities "VaultDoc-VD/Archivos/domain/entities"
	"github.com/gin-gonic/gin"
)

type GrantViewPermissionController struct {
	useCase *application.GrantViewPermissionUseCase
}

func NewGrantViewPermissionController(useCase *application.GrantViewPermissionUseCase) *GrantViewPermissionController {
	return &GrantViewPermissionController{useCase: useCase}
}

func (c *GrantViewPermissionController) Execute(ctx *gin.Context) {
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

	viewFile := entities.ViewFile{
		Id_File: input.Id_File,
		Id_User: input.Id_User,
	}

	if err := c.useCase.Execute(viewFile); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al otorgar permiso de visualización",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Permiso de visualización otorgado exitosamente",
	})
}