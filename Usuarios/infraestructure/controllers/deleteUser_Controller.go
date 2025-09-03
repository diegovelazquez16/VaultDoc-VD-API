// Usuarios/controllers/delete_user_controller.go
package controllers

import (
	"net/http"
	"strconv"
	"VaultDoc-VD/Usuarios/application"

	"github.com/gin-gonic/gin"
)

type DeleteUserController struct {
	useCase *application.DeleteUserUseCase
}

func NewDeleteUserController(useCase *application.DeleteUserUseCase) *DeleteUserController {
	return &DeleteUserController{useCase: useCase}
}

func (c *DeleteUserController) Execute(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	myId, _ := ctx.Get("userID")


	// Convertir myId a int
	var myIdInt int
	switch v := myId.(type) {
	case float64:
		myIdInt = int(v)
	case int:
		myIdInt = v
	case string:
		if val, err := strconv.Atoi(v); err == nil {
			myIdInt = val
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Formato de userID inválido en el token"})
			return
		}
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Tipo de userID no soportado"})
		return
	}

	// Evitar que un usuario se elimine a sí mismo
	if myIdInt == id {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "No puedes eliminar tu propio usuario"})
		return
	}

	if err := c.useCase.Execute(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al eliminar usuario"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado correctamente"})
}
