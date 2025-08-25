// Archivos/infrastructure/dependencies.go
package infraestructure

import (
	"VaultDoc-VD/Archivos/application"
	"VaultDoc-VD/Archivos/infraestructure/adapters"
	"VaultDoc-VD/Archivos/infraestructure/controllers"
	"VaultDoc-VD/Archivos/infraestructure/repository"
	"VaultDoc-VD/Archivos/infraestructure/routes"

	"VaultDoc-VD/core"

	"github.com/gin-gonic/gin"
)

func SetupDependencies(r *gin.Engine, dbPool *core.Conn_PostgreSQL) {
	// Inicializar repositories
	filesStorageService := adapters.NewNextcloudFileAdapter()
	filesRepo := repository.NewFilesPostgreSQLRepository(dbPool)
	changeFileRepo := repository.NewChangeFilePostgreSQLRepository(dbPool)
	viewFileRepo := repository.NewViewFilePostgreSQLRepository(dbPool)

	// Inicializar use cases
	createFileUseCase := application.NewCreateFileUseCase(filesRepo, filesStorageService)
	getFileByIdUseCase := application.NewGetFileByIdUseCase(filesRepo)
	getAllFilesUseCase := application.NewGetAllFilesUseCase(filesRepo)
	getFilesByFolderUseCase := application.NewGetFilesByFolderUseCase(filesRepo)
	updateFileUseCase := application.NewUpdateFileUseCase(filesRepo, filesStorageService)
	deleteFileUseCase := application.NewDeleteFileUseCase(filesRepo, filesStorageService)
	downloadFileUseCase := application.NewDownloadFileUseCase(filesRepo, filesStorageService)
	grantChangePermissionUseCase := application.NewGrantChangePermissionUseCase(changeFileRepo)
	removeChangePermissionUseCase := application.NewRemoveChangePermissionUseCase(changeFileRepo)
	grantViewPermissionUseCase := application.NewGrantViewPermissionUseCase(viewFileRepo)
	removeViewPermissionUseCase := application.NewRemoveViewPermissionUseCase(viewFileRepo)
	checkPermissionsUseCase := application.NewCheckPermissionsUseCase(changeFileRepo, viewFileRepo, filesRepo)

	// Inicializar controllers
	createFileController := controllers.NewCreateFileController(createFileUseCase)
	getFileByIdController := controllers.NewGetFileByIdController(getFileByIdUseCase)
	getAllFilesController := controllers.NewGetAllFilesController(getAllFilesUseCase)
	getFilesByFolderController := controllers.NewGetFilesByFolderController(getFilesByFolderUseCase)
	updateFileController := controllers.NewUpdateFileController(updateFileUseCase)
	deleteFileController := controllers.NewDeleteFileController(deleteFileUseCase)
	downloadFileController := controllers.NewDownloadFileController(downloadFileUseCase)
	grantChangePermissionController := controllers.NewGrantChangePermissionController(grantChangePermissionUseCase)
	removeChangePermissionController := controllers.NewRemoveChangePermissionController(removeChangePermissionUseCase)
	grantViewPermissionController := controllers.NewGrantViewPermissionController(grantViewPermissionUseCase)
	removeViewPermissionController := controllers.NewRemoveViewPermissionController(removeViewPermissionUseCase)
	checkPermissionsController := controllers.NewCheckPermissionsController(checkPermissionsUseCase)

	// Configurar rutas
	routes.SetupFilesRoutes(
		r,
		createFileController,
		getFileByIdController,
		getAllFilesController,
		getFilesByFolderController,
		updateFileController,
		deleteFileController,
		downloadFileController,
		grantChangePermissionController,
		removeChangePermissionController,
		grantViewPermissionController,
		removeViewPermissionController,
		checkPermissionsController,
	)
}