//Carpetas/infraestructure/controllers/GetFolderByName.go
package controllers

import (
	"VaultDoc-VD/Carpetas/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetFolderByNameController struct {
	uc application.GetFolderByNameUseCase
}

func NewGetFolderByName(uc *application.GetFolderByNameUseCase)*GetFolderByNameController{
	return&GetFolderByNameController{uc: *uc}
}

func(c *GetFolderByNameController)Execute(ctx *gin.Context){
	name := ctx.Param("folder")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al obtener carpetas: Campo del nombre del folder no incluido"})
		return
	}
	folders, err := c.uc.Execute(name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener carpetas: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"folders": folders,
	})
}