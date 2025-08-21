// Archivos/infrastructure/controllers/update_file_controller.go
package controllers

import (
	"net/http"
	"strconv"
	"VaultDoc-VD/Archivos/application"
	entities "VaultDoc-VD/Archivos/domain/entities"
	"github.com/gin-gonic/gin"
)

type UpdateFileController struct {
	useCase *application.UpdateFileUseCase
}

func NewUpdateFileController(useCase *application.UpdateFileUseCase) *UpdateFileController {
	return &UpdateFileController{useCase: useCase}
}

func (c *UpdateFileController) Execute(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID de archivo no válido",
			"error":   err.Error(),
		})
		return
	}

	var input struct {
		Departamento string `json:"departamento" binding:"required"`
		Nombre       string `json:"nombre" binding:"required"`
		Tamano       int    `json:"tamano" binding:"required"`
		Fecha        string `json:"fecha" binding:"required"`
		Folio        string `json:"folio" binding:"required"`
		Extension    string `json:"extension" binding:"required"`
		Id_Folder    int    `json:"id_folder" binding:"required"`
		Id_Uploader  int    `json:"id_uploader" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Entrada de datos no válida",
			"error":   err.Error(),
		})
		return
	}

	file := entities.Files{
		Id:           id,
		Departamento: input.Departamento,
		Nombre:       input.Nombre,
		Tamano:       input.Tamano,
		Fecha:        input.Fecha,
		Folio:        input.Folio,
		Extension:    input.Extension,
		Id_Folder:    input.Id_Folder,
		Id_Uploader:  input.Id_Uploader,
	}

	if err := c.useCase.Execute(file); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al actualizar archivo",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Archivo actualizado exitosamente",
	})
}
