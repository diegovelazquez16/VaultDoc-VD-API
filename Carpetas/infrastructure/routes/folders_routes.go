package routes

import (
	"VaultDoc-VD/Carpetas/infrastructure/controllers"
	"VaultDoc-VD/Middlewares"
	"os"

	"github.com/gin-gonic/gin"
)

func SetUpFoldersRoutes(
	r *gin.Engine, createFolderController *controllers.CreateFolderController,
	getFoldersByDepartamentController *controllers.GetFoldersByDepartamentController,
	getFolderByNameController *controllers.GetFolderByNameController,
	){

	jwtSecret := os.Getenv("JWT_SECRET")

	g := r.Group("folders")
	{
		// solo el jefe de departamento:
		g.POST("/", service.BossMiddleware(jwtSecret), createFolderController.Execute)
		// cualquier usuario autenticado
		g.GET("/:departament", getFoldersByDepartamentController.Execute)
		g.GET("/n/:folder", service.AuthMiddleware(jwtSecret), getFolderByNameController.Execute)
	}
}