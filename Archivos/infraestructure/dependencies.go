// Archivos/infrastructure/dependencies.go
package infraestructure

import (
	"VaultDoc-VD/Archivos/application"
	"VaultDoc-VD/Archivos/infraestructure/controllers"
	"VaultDoc-VD/Archivos/infraestructure/repository"
	"VaultDoc-VD/Archivos/infraestructure/routes"
	"VaultDoc-VD/core"

	"github.com/gin-gonic/gin"
)

func SetupDependencies(r *gin.Engine, dbPool *core.Conn_PostgreSQL) {
	// Inicializar repositories
	filesRepo := repository.NewFilesPostgreSQLRepository(dbPool)
	changeFileRepo := repository.NewChangeFilePostgreSQLRepository(dbPool)
	viewFileRepo := repository.NewViewFilePostgreSQLRepository(dbPool)

	// Inicializar use cases - Files CRUD
	createFileUseCase := application.NewCreateFileUseCase(filesRepo)
	getFileByIdUseCase := application.NewGetFileByIdUseCase(filesRepo)
	getAllFilesUseCase := application.NewGetAllFilesUseCase(filesRepo)
	updateFileUseCase := application.NewUpdateFileUseCase(filesRepo)
	deleteFileUseCase := application.NewDeleteFileUseCase(filesRepo)
	downloadFileUseCase := application.NewDownloadFileUseCase()

	// Inicializar use cases - Change Permissions
	grantChangePermissionUseCase := application.NewGrantChangePermissionUseCase(changeFileRepo)
	removeChangePermissionUseCase := application.NewRemoveChangePermissionUseCase(changeFileRepo)

	// Inicializar use cases - View Permissions
	grantViewPermissionUseCase := application.NewGrantViewPermissionUseCase(viewFileRepo)
	removeViewPermissionUseCase := application.NewRemoveViewPermissionUseCase(viewFileRepo)

	// Inicializar controllers - Files CRUD
	createFileController := controllers.NewCreateFileController(createFileUseCase)
	getFileByIdController := controllers.NewGetFileByIdController(getFileByIdUseCase)
	getAllFilesController := controllers.NewGetAllFilesController(getAllFilesUseCase)
	updateFileController := controllers.NewUpdateFileController(updateFileUseCase)
	deleteFileController := controllers.NewDeleteFileController(deleteFileUseCase)
	downloadFileController := controllers.NewDownloadFileController(downloadFileUseCase)

	// Inicializar controllers - Change Permissions
	grantChangePermissionController := controllers.NewGrantChangePermissionController(grantChangePermissionUseCase)
	removeChangePermissionController := controllers.NewRemoveChangePermissionController(removeChangePermissionUseCase)

	// Inicializar controllers - View Permissions
	grantViewPermissionController := controllers.NewGrantViewPermissionController(grantViewPermissionUseCase)
	removeViewPermissionController := controllers.NewRemoveViewPermissionController(removeViewPermissionUseCase)

	// Configurar rutas
	routes.SetupFilesRoutes(
		r,
		createFileController,
		getFileByIdController,
		getAllFilesController,
		updateFileController,
		deleteFileController,
		downloadFileController,
		grantChangePermissionController,
		removeChangePermissionController,
		grantViewPermissionController,
		removeViewPermissionController,
	)
} 