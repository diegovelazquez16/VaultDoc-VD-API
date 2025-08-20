// Usuarios/controllers/login_controller.go
package controllers

import (
	"net/http"

	"VaultDoc-VD/Usuarios/application"

	"github.com/gin-gonic/gin"
)

type LoginUserController struct {
	useCase *application.LoginUseCase
}

func NewLoginUserController(useCase *application.LoginUseCase) *LoginUserController {
	return &LoginUserController{useCase: useCase}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (c *LoginUserController) Execute(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	user, token, err := c.useCase.Execute(req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login exitoso",
		"token":   token,
		"user": gin.H{
			"id":        user.Id,
			"nombre":    user.Nombre,
			"apellidos": user.Apellidos,
			"email":     user.Email,
		},
	})
}
