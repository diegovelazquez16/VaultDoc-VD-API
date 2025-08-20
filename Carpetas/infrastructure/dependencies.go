package infrastructure

import (
	"VaultDoc-VD/Carpetas/application"
	"VaultDoc-VD/Carpetas/infrastructure/controllers"
	"VaultDoc-VD/Carpetas/infrastructure/repository"
	"VaultDoc-VD/Carpetas/infrastructure/routes"
	"VaultDoc-VD/core"

	"github.com/gin-gonic/gin"
)

func SetupDependenciesFolders(r *gin.Engine, dbPool *core.Conn_PostgreSQL){
	folderRepo := repository.NewFoldersPostgreSQLRepository(dbPool)

	createFolderUseCase := application.NewCreateFolderUseCase(folderRepo)

	createFolderController := controllers.NewCreateFolderController(createFolderUseCase)

	routes.SetUpFoldersRoutes(r, createFolderController)
}