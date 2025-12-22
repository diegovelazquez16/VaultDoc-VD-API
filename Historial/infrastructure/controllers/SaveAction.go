// VaultDoc-VD-API/Historial/infrastructure/controllers/SaveAction.go
package controllers

import (
	"VaultDoc-VD/Historial/application"
	"VaultDoc-VD/Historial/domain/entities"
	"VaultDoc-VD/validators"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type SaveActionController struct {
	uc application.SaveActionUseCase
}

func NewSaveActionController(uc application.SaveActionUseCase) *SaveActionController {
	return &SaveActionController{uc: uc}
}

func (c *SaveActionController) Execute(ctx *gin.Context) {
	var record entities.ReceiveHistory

	if err := ctx.ShouldBindJSON(&record); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error al leer los datos del registro de historial",
			"details": err.Error(),
		})
		return
	}

	// Validar el input
	if err := c.validateRecordInput(record); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos inválidos",
			"details": err.Error(),
		})
		return
	}

	err := c.uc.Execute(record)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error interno al crear registro en el historial",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Registro en el historial creado exitosamente",
	})
}

func (c *SaveActionController) validateRecordInput(record entities.ReceiveHistory) error {
	if strings.TrimSpace(record.Movimiento) == "" {
		return fmt.Errorf("el movimiento realizado es requerido")
	}

	// Validar que el movimiento sea uno de los valores permitidos del enum
	validMovimientos := map[string]bool{
		"Subió archivo":                     true,
		"Bajó archivo":                      true,
		"Modificó información de un archivo": true,
		"Eliminó archivo":                   true,
		"Concedió acceso al archivo":        true,
		"Se le otorgó acceso al archivo":    true,
	}
	if !validMovimientos[record.Movimiento] {
		return fmt.Errorf("el tipo de movimiento no es válido")
	}

	if strings.TrimSpace(record.Departamento) == "" || !validators.IsValidDepartamento(record.Departamento) {
		return fmt.Errorf("el departamento es requerido y debe ser válido")
	}

	if record.Id_user <= 0 {
		return fmt.Errorf("el ID de usuario es inválido")
	}

	// Validar que al menos uno de los IDs (folder o file) sea válido cuando corresponda
	// Algunos movimientos no requieren ambos
	if record.Id_file < 0 || record.Id_folder < 0 {
		return fmt.Errorf("los IDs no pueden ser negativos")
	}

	return nil
}