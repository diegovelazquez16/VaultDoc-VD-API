// Archivos/infrastructure/controllers/delete_file_controller.go (Actualizado)
package controllers

import (
	"VaultDoc-VD/Archivos/application"
	history "VaultDoc-VD/Historial/application"
	"VaultDoc-VD/Historial/domain/entities"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeleteFileController struct {
	useCase *application.DeleteFileUseCase
	historyUseCase *history.SaveActionUseCase
	uc * application.GetFileByIdUseCase
}

func NewDeleteFileController(
	useCase *application.DeleteFileUseCase, 
	historyUseCase *history.SaveActionUseCase,
	uc *application.GetFileByIdUseCase) *DeleteFileController {
	return &DeleteFileController{useCase: useCase, historyUseCase: historyUseCase, uc: uc}
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

	idUser := ctx.Param("id_user")
	if idUser == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID del usuario requerido",
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
	id_user, err := strconv.Atoi(idUser)
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

	file, err := c.uc.Execute(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Archivo no encontrado",
			"error":   err.Error(),
		})
		return
	}

	var record entities.ReceiveHistory
	record.Departamento = file.Departamento
	record.Id_user = id_user
	record.Id_folder = file.Id_Folder
	record.Id_file = file.Id
	record.Movimiento = "Eliminó archivo"

	err = c.historyUseCase.Execute(record)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error interno al crear registro en el historial",
			"details": err.Error(),
		})
		return
	}

	// 4. Respuesta exitosa
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Archivo eliminado exitosamente de BD y Nextcloud",
		"id":      id,
	})
}

