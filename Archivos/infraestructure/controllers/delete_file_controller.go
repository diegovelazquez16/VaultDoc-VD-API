// Archivos/infrastructure/controllers/delete_file_controller.go (Actualizado)
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
	// 1. Obtener ID del archivo
	idParam := ctx.Param("id")
	if idParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID del archivo requerido",
		})
		return
	}

	// 2. Convertir ID a entero
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID inválido",
			"error":   "El ID debe ser un número entero válido",
		})
		return
	}

	// 3. Ejecutar caso de uso (elimina tanto de BD como de Nextcloud)
	if err := c.useCase.Execute(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al eliminar archivo",
			"error":   err.Error(),
		})
		return
	}

	// 4. Respuesta exitosa
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Archivo eliminado exitosamente de BD y Nextcloud",
		"id":      id,
	})
}

