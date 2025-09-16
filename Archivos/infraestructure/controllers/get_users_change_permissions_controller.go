package controllers

import (
	"VaultDoc-VD/Archivos/application"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetUsersChangePermissionsController struct {
	uc application.GetUsersChangePermissionsUseCase
}

func NewGetUsersChangePermissionsController(uc application.GetUsersChangePermissionsUseCase)*GetUsersChangePermissionsController{
	return&GetUsersChangePermissionsController{uc: uc}
}

func(c *GetUsersChangePermissionsController)Execute(ctx *gin.Context){
	fileId := ctx.Param("file_id")
	if fileId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID del archivo requerido",
		})
		return
	}

	id, err := strconv.Atoi(fileId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID inválido",
			"error":   "El ID debe ser un número entero válido",
		})
		return
	}

	users_whitp, users_whitoutp, err1, err2 := c.uc.Execute(id); 
	if err1 != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al obtener usuarios con permiso de edición",
			"error":   err1.Error(),
		})
		return
	}
	if err2 != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al obtener usuarios sin permiso de edición",
			"error":   err2.Error(),
		})
		return
	}
fmt.Println(users_whitp, users_whitoutp)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Usuarios con y sin permisos de edición obtenidos",
		"users_p": users_whitp,
		"users_wp": users_whitoutp,

	})
}