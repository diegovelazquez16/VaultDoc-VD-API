package controllers

import (
	"VaultDoc-VD/Carpetas/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetFoldersByDepartamentController struct {
	uc application.GetFoldersByDepartamentUseCase
}

func NewGetFoldersByDepartamentController(uc *application.GetFoldersByDepartamentUseCase)*GetFoldersByDepartamentController{
	return&GetFoldersByDepartamentController{uc: *uc}
}

func(c *GetFoldersByDepartamentController)Execute(ctx *gin.Context){
	departament := ctx.Param("departament")
	if departament == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al obtener carpetas: Campo de departamento no incluido"})
		return
	}
	folders, err := c.uc.Execute(departament)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener carpetas: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"folders": folders,
	})
}