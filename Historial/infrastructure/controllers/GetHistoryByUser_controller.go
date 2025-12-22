// VaultDoc-VD-API/Historial/infrastructure/controllers/GetHistoryByUser.go
package controllers

import (
	"VaultDoc-VD/Historial/application"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetHistoryByUserController struct {
	uc application.GetHistoryByUserUseCase
}

func NewGetHistoryByUserController(uc application.GetHistoryByUserUseCase) *GetHistoryByUserController {
	return &GetHistoryByUserController{uc: uc}
}

func (c *GetHistoryByUserController) Execute(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID de usuario inválido",
			"details": err.Error(),
		})
		return
	}

	history, err := c.uc.Execute(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener historial por usuario",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, history)
}