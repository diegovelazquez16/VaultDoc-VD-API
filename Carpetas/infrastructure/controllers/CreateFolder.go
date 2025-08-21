package controllers

import (
	"VaultDoc-VD/Carpetas/application"
	"VaultDoc-VD/Carpetas/domain/entities"
	"VaultDoc-VD/validators"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type CreateFolderController struct {
	uc application.CreateFolderUseCase
}

func NewCreateFolderController(uc *application.CreateFolderUseCase)*CreateFolderController{
	return&CreateFolderController{uc: *uc}
}

func(c *CreateFolderController)Execute(ctx *gin.Context){
	var folder entities.Folders
	if err := ctx.ShouldBindJSON(&folder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error al leer los datos de la carpeta",
			"details": err.Error(),
		})
		return
	}

	if err := c.validateFolderInput(folder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos de usuario inválidos",
			"details": err.Error(),
		})
		return
	}

	// Ejecutar caso de uso
	createdFolder, err := c.uc.Execute(folder.Name, folder.Departamento, folder.Id_uploader)
	if err != nil {
		// Clasificar errores para respuestas más específicas
		if strings.Contains(err.Error(), "ya está registrado") ||
			strings.Contains(err.Error(), "ya existe") {
			ctx.JSON(http.StatusConflict, gin.H{
				"error":   "Folder ya registrado",
				"details": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error interno al crear carpeta",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Carpeta creada exitosamente",
		"folder": gin.H{
			"id":          	createdFolder.Id,
			"name":        	createdFolder.Name,
			"departament": 	createdFolder.Departamento,
			"id_uploader":  createdFolder.Id_uploader,
		},
		"sync_status": "local_saved",
	})
}

func (c *CreateFolderController) validateFolderInput(folder entities.Folders) error {
	if strings.TrimSpace(folder.Name) == "" {
		return fmt.Errorf("el email es requerido")
	}

	// Validar departamento si se proporciona
	if folder.Departamento != "" && !validators.IsValidDepartamento(folder.Departamento) {
		return fmt.Errorf("el departamento debe ser: Finanzaz, Operativo o General")
	}

	// Validar id del usuario quien lo subió
	if folder.Id_uploader < 0 {
		return fmt.Errorf("el id_uploader debe ser un número positivo")
	}

	return nil
}