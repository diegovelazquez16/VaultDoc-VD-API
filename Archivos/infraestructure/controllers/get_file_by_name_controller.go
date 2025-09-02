package controllers

import (
	"VaultDoc-VD/Archivos/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SearchFileController struct {
	uc application.SearchFileUseCase
}

func NewSearchFileController(uc application.SearchFileUseCase)*SearchFileController{
	return&SearchFileController{uc: uc}
}

func(c *SearchFileController)Execute(ctx *gin.Context){
	name := ctx.Param("filename")

	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error al obtener nombre de archivo a buscar",
			"error": "No se encontró el parámetro de nombre",
		})
	}
	
	// Ejecutar el caso de uso
	files, err := c.uc.Execute(name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al buscar archivos",
			"error":   err.Error(),
		})
		return
	}

	// Responder con los archivos encontrados
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Archivos del folder obtenidos exitosamente",
		"data":    files,
	})
}