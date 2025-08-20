package routes

import (
	"VaultDoc-VD/Carpetas/infrastructure/controllers"

	"github.com/gin-gonic/gin"
)

func SetUpFoldersRoutes(
	r *gin.Engine, createFolderController *controllers.CreateFolderController,
	getFoldersByDepartamentController *controllers.GetFoldersByDepartamentController,
	){
	g := r.Group("folders")
	{
		g.POST("/", createFolderController.Execute)
		g.GET("/:departament", getFoldersByDepartamentController.Execute)
	}
}