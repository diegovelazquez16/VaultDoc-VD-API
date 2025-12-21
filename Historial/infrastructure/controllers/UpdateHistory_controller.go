// VaultDoc-VD-API/Historial/infrastructure/controllers/UpdateHistory.go
package controllers

import (
	"VaultDoc-VD/Historial/application"
	"VaultDoc-VD/Historial/domain/entities"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UpdateHistoryController struct {
	uc application.UpdateHistoryUseCase
}

func NewUpdateHistoryController(uc application.UpdateHistoryUseCase) *UpdateHistoryController {
	return &UpdateHistoryController{uc: uc}
}

func (c *UpdateHistoryController) Execute(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID inválido",
			"details": err.Error(),
		})
		return
	}

	var record entities.ReceiveHistory
	if err := ctx.ShouldBindJSON(&record); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error al leer los datos del registro",
			"details": err.Error(),
		})
		return
	}

	err = c.uc.Execute(id, record)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al actualizar el registro en el historial",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Registro actualizado exitosamente",
	})
}