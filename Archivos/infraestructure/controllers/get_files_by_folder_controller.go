// Archivos/infrastructure/controllers/get_files_by_folder_controller.go
package controllers

import (
	"net/http"
	"strconv"
	"VaultDoc-VD/Archivos/application"
	"github.com/gin-gonic/gin"
)

type GetFilesByFolderController struct {
	useCase *application.GetFilesByFolderUseCase
}

func NewGetFilesByFolderController(useCase *application.GetFilesByFolderUseCase) *GetFilesByFolderController {
	return &GetFilesByFolderController{useCase: useCase}
}

func (c *GetFilesByFolderController) Execute(ctx *gin.Context) {
	// Obtener el folderId de los parámetros de la URL
	folderIdParam := ctx.Param("folderId")
	folderId, err := strconv.Atoi(folderIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID de folder inválido",
			"error":   "El ID debe ser un número entero",
		})
		return
	}

	// Ejecutar el caso de uso
	files, err := c.useCase.Execute(folderId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al obtener archivos del folder",
			"error":   err.Error(),
		})
		return
	}

	// Responder con los archivos encontrados
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Archivos del folder obtenidos exitosamente",
		"data":    files,
		"count":   len(files),
	})
}