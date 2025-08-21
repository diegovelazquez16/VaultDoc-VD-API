// Archivos/infrastructure/controllers/check_permissions_controller.go
package controllers

import (
	"net/http"
	"strconv"
	"VaultDoc-VD/Archivos/application"
	"github.com/gin-gonic/gin"
)

type CheckPermissionsController struct {
	useCase *application.CheckPermissionsUseCase
}

func NewCheckPermissionsController(useCase *application.CheckPermissionsUseCase) *CheckPermissionsController {
	return &CheckPermissionsController{useCase: useCase}
}

func (c *CheckPermissionsController) Execute(ctx *gin.Context) {
	// Obtener parámetros de la URL
	fileIdStr := ctx.Param("fileId")
	userIdStr := ctx.Param("userId")

	fileId, err := strconv.Atoi(fileIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID de archivo no válido",
			"error":   err.Error(),
		})
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID de usuario no válido",
			"error":   err.Error(),
		})
		return
	}

	permissions, err := c.useCase.Execute(fileId, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al verificar permisos",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":     "Permisos verificados exitosamente",
		"file_id":     fileId,
		"user_id":     userId,
		"permissions": permissions,
	})
}