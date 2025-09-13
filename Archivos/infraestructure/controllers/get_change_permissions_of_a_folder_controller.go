package controllers

import (
	"VaultDoc-VD/Archivos/application"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetChangePermissionsOfAFolderController struct {
	uc application.GetChangePermissionsOfAFolderUseCase
}

func NewGetChangePermissionsOfAFolderController(uc application.GetChangePermissionsOfAFolderUseCase)*GetChangePermissionsOfAFolderController{
	return&GetChangePermissionsOfAFolderController{uc: uc}
}

func(c *GetChangePermissionsOfAFolderController)Execute(ctx *gin.Context){
	idFolder := ctx.Param("id_folder")
	if idFolder == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID del folder requerido",
		})
		return
	}

	idUser := ctx.Param("id_user")
	if idUser == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID del usuario requerido",
		})
		return
	}

	// 2. Convertir ID a entero
	id_folder, err := strconv.Atoi(idFolder)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID inválido",
			"error":   "El ID debe ser un número entero válido",
		})
		return
	}
	id_user, err := strconv.Atoi(idUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID inválido",
			"error":   "El ID debe ser un número entero válido",
		})
		return
	}

	permissions, err := c.uc.Execute(id_folder, id_user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al obtener permisos de modificación en los archivos de un folder",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Permisos obtenidos",
		"permissions": permissions,
	})
}