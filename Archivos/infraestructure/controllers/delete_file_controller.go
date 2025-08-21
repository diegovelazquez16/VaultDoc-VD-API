// Archivos/infrastructure/controllers/delete_file_controller.go
package controllers

import (
	"net/http"
	"strconv"
	"VaultDoc-VD/Archivos/application"
	"github.com/gin-gonic/gin"
)

type DeleteFileController struct {
	useCase *application.DeleteFileUseCase
}

func NewDeleteFileController(useCase *application.DeleteFileUseCase) *DeleteFileController {
	return &DeleteFileController{useCase: useCase}
}

func (c *DeleteFileController) Execute(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID de archivo no válido",
			"error":   err.Error(),
		})
		return
	}

	if err := c.useCase.Execute(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al eliminar archivo",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Archivo eliminado exitosamente",
	})
}