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
	// 1. Obtener y validar el ID del archivo
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID de archivo no válido",
			"error":   err.Error(),
		})
		return
	}

	// Validar que el ID sea positivo
	if id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID de archivo no válido",
			"error":   "El ID debe ser un número positivo",
		})
		return
	}

	// 2. Ejecutar el caso de uso para eliminar archivo
	if err := c.useCase.Execute(id); err != nil {
		// Determinar el código de error basado en el tipo
		statusCode := http.StatusInternalServerError
		if err.Error() == "archivo no encontrado en base de datos" {
			statusCode = http.StatusNotFound
		}
		
		ctx.JSON(statusCode, gin.H{
			"message": "Error al eliminar archivo",
			"error":   err.Error(),
		})
		return
	}

	// 3. Respuesta exitosa
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Archivo eliminado exitosamente",
		"id":      id,
	})
}