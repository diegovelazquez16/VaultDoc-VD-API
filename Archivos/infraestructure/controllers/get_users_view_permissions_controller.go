package controllers

import (
	"VaultDoc-VD/Archivos/application"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetUsersViewPermissionsController struct {
	uc application.GetUsersViewPermissionsUseCase
}

func NewGetUsersViewPermissionsController(uc application.GetUsersViewPermissionsUseCase)*GetUsersViewPermissionsController{
	return&GetUsersViewPermissionsController{uc: uc}
}

func(c *GetUsersViewPermissionsController)Execute(ctx *gin.Context){
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
			"message": "Error al obtener usuarios con permiso de lectura",
			"error":   err1.Error(),
		})
		return
	}
	if err2 != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al obtener usuarios sin permiso de lectura",
			"error":   err2.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Usuarios con y sin permisos de lactura obtenidos",
		"users_p": users_whitp,
		"users_wp": users_whitoutp,
	})
}