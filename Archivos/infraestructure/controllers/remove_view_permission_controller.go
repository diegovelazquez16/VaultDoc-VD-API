// Archivos/infrastructure/controllers/remove_view_permission_controller.go
package controllers

import (
	"VaultDoc-VD/Archivos/application"
	entities "VaultDoc-VD/Archivos/domain/entities"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RemoveViewPermissionController struct {
	useCase *application.RemoveViewPermissionUseCase
}

func NewRemoveViewPermissionController(useCase *application.RemoveViewPermissionUseCase) *RemoveViewPermissionController {
	return &RemoveViewPermissionController{useCase: useCase}
}

func (c *RemoveViewPermissionController) Execute(ctx *gin.Context) {
	idUser := ctx.Param("id_user")
	if idUser == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID del usuario requerido",
		})
		return
	}
	id_user, err := strconv.Atoi(idUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID inválido",
			"error":   "El ID debe ser un número entero válido",
		})
		return
	}

	idFile := ctx.Param("id_file")
	if idFile == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID del archivo requerido",
		})
		return
	}
	id_file, err := strconv.Atoi(idUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID inválido",
			"error":   "El ID debe ser un número entero válido",
		})
		return
	}


	viewFile := entities.ViewFile{
		Id_File: id_file,
		Id_User: id_user,
	}

	if err := c.useCase.Execute(viewFile); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al revocar permiso de visualización",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Permiso de visualización revocado exitosamente",
	})
}