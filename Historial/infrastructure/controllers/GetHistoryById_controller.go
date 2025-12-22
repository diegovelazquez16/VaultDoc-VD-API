// VaultDoc-VD-API/Historial/infrastructure/controllers/GetHistoryByID.go
package controllers

import (
	"VaultDoc-VD/Historial/application"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetHistoryByIDController struct {
	uc application.GetHistoryByIDUseCase
}

func NewGetHistoryByIDController(uc application.GetHistoryByIDUseCase) *GetHistoryByIDController {
	return &GetHistoryByIDController{uc: uc}
}

func (c *GetHistoryByIDController) Execute(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID inválido",
			"details": err.Error(),
		})
		return
	}

	record, err := c.uc.Execute(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Registro no encontrado",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, record)
}