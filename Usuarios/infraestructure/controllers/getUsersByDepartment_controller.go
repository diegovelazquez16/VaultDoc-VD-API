// Usuarios/infraestructure/controllers/getUsersByDepartment_controller.go
package controllers

import (
	"net/http"

	"VaultDoc-VD/Usuarios/application"

	"github.com/gin-gonic/gin"
)

type GetUsersByDepartmentController struct {
	useCase *application.GetUsersByDepartment
}

func NewGetUsersByDepartmentController(useCase *application.GetUsersByDepartment) *GetUsersByDepartmentController {
	return &GetUsersByDepartmentController{useCase: useCase}
}

func (c *GetUsersByDepartmentController) Execute(ctx *gin.Context) {
	departamento := ctx.Param("departamento")
	if departamento == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro departamento es requerido"})
		return
	}

	users, err := c.useCase.Execute(departamento)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuarios por departamento: " + err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, users)
}