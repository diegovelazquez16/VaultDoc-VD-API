// Archivos/infrastructure/controllers/get_all_files_controller.go
package controllers

import (
	"net/http"
	"VaultDoc-VD/Archivos/application"
	"github.com/gin-gonic/gin"
)

type GetAllFilesController struct {
	useCase *application.GetAllFilesUseCase
}

func NewGetAllFilesController(useCase *application.GetAllFilesUseCase) *GetAllFilesController {
	return &GetAllFilesController{useCase: useCase}
}

func (c *GetAllFilesController) Execute(ctx *gin.Context) {
	files, err := c.useCase.Execute()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al obtener archivos",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Archivos obtenidos exitosamente",
		"data":    files,
	})
}