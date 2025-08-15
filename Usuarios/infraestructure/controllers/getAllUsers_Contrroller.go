package controllers

import (
	"net/http"

	"VaultDoc-VD/Usuarios/application"

	"github.com/gin-gonic/gin"
)

type GetAllUsersController struct {
	useCase *application.GetUsers
}

func NewGetAllUsersController(useCase *application.GetUsers) *GetAllUsersController {
	return &GetAllUsersController{useCase: useCase}
}

func (c *GetAllUsersController) Execute(ctx *gin.Context) {
	users, err := c.useCase.Execute()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los usuarios: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}
