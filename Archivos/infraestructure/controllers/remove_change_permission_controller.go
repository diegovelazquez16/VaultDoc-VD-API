// Archivos/infrastructure/controllers/remove_change_permission_controller.go
package controllers

import (
	"VaultDoc-VD/Archivos/application"
	entities "VaultDoc-VD/Archivos/domain/entities"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RemoveChangePermissionController struct {
	useCase *application.RemoveChangePermissionUseCase
}

func NewRemoveChangePermissionController(useCase *application.RemoveChangePermissionUseCase) *RemoveChangePermissionController {
	return &RemoveChangePermissionController{useCase: useCase}
}

func (c *RemoveChangePermissionController) Execute(ctx *gin.Context) {
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

	changeFile := entities.ChangeFile{
		Id_File: id_file,
		Id_User: id_user,
	}

	if err := c.useCase.Execute(changeFile); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al revocar permiso de edición",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Permiso de edición revocado exitosamente",
	})
}
