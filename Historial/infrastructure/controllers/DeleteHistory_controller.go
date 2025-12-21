// VaultDoc-VD-API/Historial/infrastructure/controllers/DeleteHistory.go
package controllers

import (
	"VaultDoc-VD/Historial/application"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeleteHistoryController struct {
	uc application.DeleteHistoryUseCase
}

func NewDeleteHistoryController(uc application.DeleteHistoryUseCase) *DeleteHistoryController {
	return &DeleteHistoryController{uc: uc}
}

func (c *DeleteHistoryController) Execute(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID inválido",
			"details": err.Error(),
		})
		return
	}

	err = c.uc.Execute(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al eliminar el registro del historial",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Registro eliminado exitosamente",
	})
}