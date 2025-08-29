// Archivos/infrastructure/controllers/download_file_controller.go (Actualizado)
package controllers

import (
	"VaultDoc-VD/Archivos/application"
	history "VaultDoc-VD/Historial/application"
	"VaultDoc-VD/Historial/domain/entities"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DownloadFileController struct {
	useCase *application.DownloadFileUseCase
	historyUseCase *history.SaveActionUseCase
	uc * application.GetFileByIdUseCase
}

func NewDownloadFileController(
	useCase *application.DownloadFileUseCase,
	historyUseCase *history.SaveActionUseCase,
	uc * application.GetFileByIdUseCase,
	) *DownloadFileController {
	return &DownloadFileController{useCase: useCase, historyUseCase: historyUseCase, uc: uc}
}

func (c *DownloadFileController) Execute(ctx *gin.Context) {
	// 1. Obtener el ID del archivo de los parámetros de la URL
	idParam := ctx.Param("id")
	if idParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID del archivo requerido",
			"error":   "Especifica el ID del archivo a descargar",
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

	// 3. Ejecutar caso de uso para descargar desde Nextcloud
	content, fileName, err := c.useCase.Execute(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Archivo no encontrado",
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
	record.Movimiento = "Bajó archivo"

	err = c.historyUseCase.Execute(record)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error interno al crear registro en el historial",
			"details": err.Error(),
		})
		return
	}

	// 4. Configurar headers para la descarga
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", "attachment; filename="+fileName)
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Length", strconv.Itoa(len(content)))

	// 5. Enviar el contenido del archivo
	ctx.Data(http.StatusOK, "application/octet-stream", content)
}