// Archivos/infrastructure/controllers/get_file_by_id_controller.go
package controllers

import (
	"net/http"
	"strconv"
	"VaultDoc-VD/Archivos/application"
	"github.com/gin-gonic/gin"
)

type GetFileByIdController struct {
	useCase *application.GetFileByIdUseCase
}

func NewGetFileByIdController(useCase *application.GetFileByIdUseCase) *GetFileByIdController {
	return &GetFileByIdController{useCase: useCase}
}

func (c *GetFileByIdController) Execute(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID de archivo no válido",
			"error":   err.Error(),
		})
		return
	}

	file, err := c.useCase.Execute(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Archivo no encontrado",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Archivo encontrado",
		"data":    file,
	})
}
