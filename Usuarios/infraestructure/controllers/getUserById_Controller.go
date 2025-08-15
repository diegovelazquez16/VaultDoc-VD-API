package controllers

import (
	"VaultDoc-VD/Usuarios/application"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetUserByIdController struct {
	useCase *application.GetUserById
}

func NewGetUserByIdUseController(usecase *application.GetUserById) *GetUserByIdController {
	return &GetUserByIdController{useCase: usecase}
}

func (c *GetUserByIdController) Execute(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID no valido"})
		return
	}

	user, err := c.useCase.Execute(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Usuario inexistente"})
		return
	}
	ctx.JSON(http.StatusOK, user)
}
