// Archivos/infrastructure/controllers/grant_change_permission_controller.go
package controllers

import (
	"VaultDoc-VD/Archivos/application"
	entities "VaultDoc-VD/Archivos/domain/entities"
	history "VaultDoc-VD/Historial/application"
	historyEntities "VaultDoc-VD/Historial/domain/entities"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type GrantChangePermissionController struct {
	useCase *application.GrantChangePermissionUseCase
	historyUseCase *history.SaveActionUseCase
	uc * application.GetFileByIdUseCase
}

func NewGrantChangePermissionController(
	useCase *application.GrantChangePermissionUseCase,
	historyUseCase *history.SaveActionUseCase,
	uc * application.GetFileByIdUseCase,
	) *GrantChangePermissionController {
	return &GrantChangePermissionController{useCase: useCase, historyUseCase: historyUseCase, uc: uc}
}

func (c *GrantChangePermissionController) Execute(ctx *gin.Context) {
	idUser := ctx.Param("id_user")
	if idUser == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID del usuario requerido",
		})
		return
	}
	id_user, err := strconv.Atoi(idUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID inválido",
			"error":   "El ID debe ser un número entero válido",
		})
		return
	}

	var input struct {
		Id_File int `json:"id_file" binding:"required"`
		Id_User int `json:"id_user" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Entrada de datos no válida",
			"error":   err.Error(),
		})
		return
	}

	// Validaciones básicas
	if input.Id_File <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID de archivo no válido",
			"error":   "El ID del archivo debe ser un número positivo",
		})
		return
	}

	if input.Id_User <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID de usuario no válido",
			"error":   "El ID del usuario debe ser un número positivo",
		})
		return
	}

	changeFile := entities.ChangeFile{
		Id_File: input.Id_File,
		Id_User: input.Id_User,
	}

	if err := c.useCase.Execute(changeFile); err != nil {
		statusCode := http.StatusInternalServerError
		// Determinar código de estado basado en el error
		if strings.Contains(err.Error(), "no existe") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(err.Error(), "ya existe") {
			statusCode = http.StatusConflict
		}

		ctx.JSON(statusCode, gin.H{
			"message": "Error al otorgar permiso de edición",
			"error":   err.Error(),
		})
		return
	}

	file, err := c.uc.Execute(input.Id_File)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Archivo no encontrado",
			"error":   err.Error(),
		})
		return
	}

	var record historyEntities.ReceiveHistory
	record.Departamento = file.Departamento
	record.Id_user = id_user
	record.Id_folder = file.Id_Folder
	record.Id_file = file.Id
	record.Movimiento = "Concedió acceso al archivo"

	err = c.historyUseCase.Execute(record)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error interno al crear registro en el historial",
			"details": err.Error(),
		})
		return
	}

	record.Departamento = file.Departamento
	record.Id_user = input.Id_User
	record.Id_folder = file.Id_Folder
	record.Id_file = file.Id
	record.Movimiento = "Se le otorgó acceso al archivo"

	err = c.historyUseCase.Execute(record)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error interno al crear registro en el historial",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Permiso de edición otorgado exitosamente",
		"id_file": input.Id_File,
		"id_user": input.Id_User,
	})
}