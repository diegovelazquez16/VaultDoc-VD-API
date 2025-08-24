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

func NewSaveActionController(uc application.SaveActionUseCase)*SaveActionController{
	return&SaveActionController{uc: uc}
}

func(c *SaveActionController)Execute(ctx *gin.Context){
	var record entities.ReceiveHistory
	if err := ctx.ShouldBindJSON(&record); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error al leer los datos del registro de historial",
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
		"message": "Registro en el historial creado",
	})

	/* Respuesta exitosa con información del usuario creado
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Usuario creado exitosamente",
		"record": gin.H{
			"id":           	createdRecord.Id,
			"movimiento":   	createdRecord.Movimiento,
			"departamento": 	createdRecord.Departamento,
			"id_folder":    	createdRecord.Id_folder,
			"id_file":      	createdRecord.Id_file,
			"id_user": 			createdRecord.Id_user,
			"fecha_registro":	createdRecord.Fecha_registro,
		},
		"sync_status": "local_saved",
	})
		*/
}

func (c *SaveActionController) validateRecordInput(record entities.ReceiveHistory) error {
	if strings.TrimSpace(record.Movimiento) == "" {
		return fmt.Errorf("El movimiento realizado es requerido")
	}

	if strings.TrimSpace(record.Departamento) == "" || !validators.IsValidDepartamento(record.Departamento) {
		return fmt.Errorf("El departamento es requerido")
	}

	if record.Id_file < 0 || record.Id_folder < 0 || record.Id_user < 0 {
		return fmt.Errorf("Alguna de las Id recibidas son inválidas inválido")
	}

	return nil
}
