//Carpetas/infraestructure/dependencies.go
package infrastructure

import (
	"VaultDoc-VD/Carpetas/application"
	"VaultDoc-VD/Carpetas/infrastructure/controllers"
	"VaultDoc-VD/Carpetas/infrastructure/repository"
	"VaultDoc-VD/Carpetas/infrastructure/routes"
	"VaultDoc-VD/Carpetas/infrastructure/services/adapters"
	"VaultDoc-VD/core"

	"github.com/gin-gonic/gin"
)

func SetupDependenciesFolders(r *gin.Engine, dbPool *core.Conn_PostgreSQL){
	folderRepo := repository.NewFoldersPostgreSQLRepository(dbPool)
	nextcloudAdapter := adapters.NewNextcloudAdapter()

	createFolderUseCase := application.NewCreateFolderUseCase(folderRepo, nextcloudAdapter)
	getFoldersByDepartamentUseCase := application.NewGetFoldersByDepartamentUseCase(folderRepo, nextcloudAdapter)
	getFolderByNameUseCase := application.NewGetFolderByName(folderRepo, nextcloudAdapter)

	createFolderController := controllers.NewCreateFolderController(createFolderUseCase)
	getFoldersByDepartamentController := controllers.NewGetFoldersByDepartamentController(getFoldersByDepartamentUseCase)
	getFolderByNameController := controllers.NewGetFolderByName(getFolderByNameUseCase)

	routes.SetUpFoldersRoutes(r, createFolderController, getFoldersByDepartamentController, getFolderByNameController)
}