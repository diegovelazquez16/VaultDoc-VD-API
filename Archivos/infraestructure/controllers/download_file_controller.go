// Archivos/infrastructure/controllers/download_file_controller.go
package controllers

import (
	"net/http"
	"path/filepath"
	
	"VaultDoc-VD/Archivos/application"
	"github.com/gin-gonic/gin"
)

type DownloadFileController struct {
	useCase *application.DownloadFileUseCase
}

func NewDownloadFileController(useCase *application.DownloadFileUseCase) *DownloadFileController {
	return &DownloadFileController{useCase: useCase}
}

func (c *DownloadFileController) Execute(ctx *gin.Context) {
	// 1. Obtener la ruta del parámetro wildcard (*dir)
	dir := ctx.Param("dir")
	if dir == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Ruta de directorio requerida",
			"error":   "Especifica la ruta del archivo a descargar",
		})
		return
	}

	// 2. Extraer el nombre del archivo de la ruta
	filename := filepath.Base(dir)
	dirPath := filepath.Dir(dir)
	
	// Si dir es solo el nombre del archivo, dirPath será "."
	if dirPath == "." {
		dirPath = ""
	}

	// 3. Ejecutar caso de uso para verificar que el archivo existe
	fullPath, err := c.useCase.Execute(dirPath, filename)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Archivo no encontrado",
			"error":   err.Error(),
		})
		return
	}

	// 4. Configurar headers para la descarga
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Header("Content-Type", "application/octet-stream")

	// 5. Enviar el archivo
	ctx.File(fullPath)
}