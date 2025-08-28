package controllers

import (
	"VaultDoc-VD/Carpetas/application"
	"net/http"
	"github.com/gin-gonic/gin"
	"fmt"
)

type GetFolderByMyDepartamentController struct { 
	uc application.GetFolderByMyDepartamentUseCase
}

func NewGetFoldersByMyDepartamentController(uc *application.GetFolderByMyDepartamentUseCase) *GetFolderByMyDepartamentController {
	return &GetFolderByMyDepartamentController{uc: *uc}
}

func(c *GetFolderByMyDepartamentController) Execute(ctx *gin.Context) {
	userDepartmentInterface, exists := ctx.Get("department")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
            "message": "Error del sistema",
			"error":   "No se pudo obtener el departamento del usuario del token",
		})
		return
	}

	var userDepartment string
	switch v := userDepartmentInterface.(type) {
	case string:
		userDepartment = v
	case interface{}:
		if str, ok := v. (string); ok {
			userDepartment = str
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error del sistema",
				"error":   "Formato de departamento inválido en el token",
			})
			return
		}
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error del sistema",
			"error": "Tipo de departamento no reconocido",
		})
		return
	}

	validDepartments := []string{
		"Dirección General", 
		"Área Técnica", 
		"Comisaria", 
		"Coordinación Juridica", 
		"Gerencia Administrativa", 
		"Gerencia Operativa", 
		"Departamento de Finanzas", 
		"Departamento de Planeación", 
		"Departamento de Sistema Eléctrico", 
		"Departamento de Sistema Hidrosánitario y Aire Acondicionado", 
		"Departamento de Mantenimiento General", 
		"Departamento de Voz y Datos", 
		"Departamento de Seguridad e Higiene",
	}
	
	isValidDepartment := false
	for _, dept := range validDepartments {
		if userDepartment == dept { 
			isValidDepartment = true
			break
		}
	}
	
	if !isValidDepartment {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Departamento no válido",
			"error":   fmt.Sprintf("El departamento del usuario '%s' no está permitido", userDepartment),
		})
		return
	}

	folders, err := c.uc.Execute(userDepartment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error en el sistema",
			"error": "Error al obtener carpetas: " + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"folders": folders,
	})


}