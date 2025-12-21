// VaultDoc-VD-API/Historial/infrastructure/controllers/GetAllHistory.go
package controllers

import (
	"VaultDoc-VD/Historial/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetAllHistoryController struct {
	uc application.GetAllHistoryUseCase
}

func NewGetAllHistoryController(uc application.GetAllHistoryUseCase) *GetAllHistoryController {
	return &GetAllHistoryController{uc: uc}
}

func (c *GetAllHistoryController) Execute(ctx *gin.Context) {
	history, err := c.uc.Execute()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener el historial completo",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, history)
}
