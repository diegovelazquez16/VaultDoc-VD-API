package controllers

import (
	"VaultDoc-VD/Carpetas/application"

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

	}
//	folders, err := c.uc.Execute(departament)
}