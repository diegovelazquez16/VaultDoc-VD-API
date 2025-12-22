// VaultDoc-VD-API/Historial/infrastructure/controllers/GetHistoryByFolder.go
package controllers

import (
	"VaultDoc-VD/Historial/application"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetHistoryByFolderController struct {
	uc application.GetHistoryByFolderUseCase
}

func NewGetHistoryByFolderController(uc application.GetHistoryByFolderUseCase) *GetHistoryByFolderController {
	return &GetHistoryByFolderController{uc: uc}
}

func (c *GetHistoryByFolderController) Execute(ctx *gin.Context) {
	folderIDStr := ctx.Param("folderId")
	folderID, err := strconv.Atoi(folderIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID de carpeta inválido",
			"details": err.Error(),
		})
		return
	}

	history, err := c.uc.Execute(folderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener historial por carpeta",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, history)
}