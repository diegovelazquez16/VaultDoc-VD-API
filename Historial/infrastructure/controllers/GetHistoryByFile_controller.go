// VaultDoc-VD-API/Historial/infrastructure/controllers/GetHistoryByFile.go
package controllers

import (
	"VaultDoc-VD/Historial/application"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetHistoryByFileController struct {
	uc application.GetHistoryByFileUseCase
}

func NewGetHistoryByFileController(uc application.GetHistoryByFileUseCase) *GetHistoryByFileController {
	return &GetHistoryByFileController{uc: uc}
}

func (c *GetHistoryByFileController) Execute(ctx *gin.Context) {
	fileIDStr := ctx.Param("fileId")
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID de archivo inválido",
			"details": err.Error(),
		})
		return
	}

	history, err := c.uc.Execute(fileID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener historial por archivo",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, history)
}