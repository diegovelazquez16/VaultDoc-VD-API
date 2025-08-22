package controllers

import (
	"VaultDoc-VD/Historial/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetHistoryController struct {
	uc application.GetHistoryUseCase
}

func NewGetHistoryController(uc application.GetHistoryUseCase)*GetHistoryController{
	return&GetHistoryController{uc: uc}
}

func(c *GetHistoryController)Execute(ctx *gin.Context){
	departament := ctx.Param("departament")
	history, err := c.uc.Execute(departament)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el historial: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, history)
}