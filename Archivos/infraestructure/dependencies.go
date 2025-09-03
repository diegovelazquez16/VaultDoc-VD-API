// Archivos/infrastructure/dependencies.go
package infraestructure

import (
	"VaultDoc-VD/Archivos/application"
	historyApplication "VaultDoc-VD/Historial/application"
	"VaultDoc-VD/Archivos/infraestructure/adapters"
	"VaultDoc-VD/Archivos/infraestructure/controllers"
	"VaultDoc-VD/Archivos/infraestructure/repository"
	historyRepository "VaultDoc-VD/Historial/infrastructure/repository"
	"VaultDoc-VD/Archivos/infraestructure/routes"
	folderRepository "VaultDoc-VD/Carpetas/infrastructure/repository"

	"VaultDoc-VD/core"

	"github.com/gin-gonic/gin"
)

func SetupDependencies(r *gin.Engine, dbPool *core.Conn_PostgreSQL) {
	// Inicializar repositories y adaptadores
	filesStorageService := adapters.NewNextcloudFileAdapter()
	userService := adapters.NewUserPostgreSQLAdapter(dbPool)
	filesRepo := repository.NewFilesPostgreSQLRepository(dbPool)
	changeFileRepo := repository.NewChangeFilePostgreSQLRepository(dbPool)
	viewFileRepo := repository.NewViewFilePostgreSQLRepository(dbPool)
	historyRepo := historyRepository.NewHistoryPostgreSQLRepository(dbPool)
	folderRepo := folderRepository.NewFoldersPostgreSQLRepository(dbPool)

	// Inicializar use cases
	createFileUseCase := application.NewCreateFileUseCase(filesRepo, filesStorageService, changeFileRepo, viewFileRepo, userService, folderRepo)
	getFileByIdUseCase := application.NewGetFileByIdUseCase(filesRepo)
	getAllFilesUseCase := application.NewGetAllFilesUseCase(filesRepo)
	getFilesByFolderUseCase := application.NewGetFilesByFolderUseCase(filesRepo)
	getFilesByNameUseCase := application.NewGetFileByNameUseCase(filesRepo)
	updateFileUseCase := application.NewUpdateFileUseCase(filesRepo, filesStorageService, folderRepo)
	deleteFileUseCase := application.NewDeleteFileUseCase(filesRepo, filesStorageService)
	downloadFileUseCase := application.NewDownloadFileUseCase(filesRepo, filesStorageService)
	grantChangePermissionUseCase := application.NewGrantChangePermissionUseCase(changeFileRepo)
	removeChangePermissionUseCase := application.NewRemoveChangePermissionUseCase(changeFileRepo)
	grantViewPermissionUseCase := application.NewGrantViewPermissionUseCase(viewFileRepo)
	removeViewPermissionUseCase := application.NewRemoveViewPermissionUseCase(viewFileRepo)
	checkPermissionsUseCase := application.NewCheckPermissionsUseCase(changeFileRepo, viewFileRepo, filesRepo)
	saveRecordUseCase := historyApplication.NewSaveActionsUseCase(historyRepo)
	searchFileUseCase := application.NewSearchFileUseCase(filesRepo)

	// Inicializar controllers
	createFileController := controllers.NewCreateFileController(createFileUseCase, saveRecordUseCase, getFilesByNameUseCase)
	getFileByIdController := controllers.NewGetFileByIdController(getFileByIdUseCase)
	getAllFilesController := controllers.NewGetAllFilesController(getAllFilesUseCase)
	getFilesByFolderController := controllers.NewGetFilesByFolderController(getFilesByFolderUseCase)
	updateFileController := controllers.NewUpdateFileController(updateFileUseCase, saveRecordUseCase, getFileByIdUseCase)
	deleteFileController := controllers.NewDeleteFileController(deleteFileUseCase, saveRecordUseCase, getFileByIdUseCase)
	downloadFileController := controllers.NewDownloadFileController(downloadFileUseCase, saveRecordUseCase, getFileByIdUseCase)
	grantChangePermissionController := controllers.NewGrantChangePermissionController(grantChangePermissionUseCase, saveRecordUseCase, getFileByIdUseCase)
	removeChangePermissionController := controllers.NewRemoveChangePermissionController(removeChangePermissionUseCase)
	grantViewPermissionController := controllers.NewGrantViewPermissionController(grantViewPermissionUseCase, saveRecordUseCase, getFileByIdUseCase)
	removeViewPermissionController := controllers.NewRemoveViewPermissionController(removeViewPermissionUseCase)
	checkPermissionsController := controllers.NewCheckPermissionsController(checkPermissionsUseCase)
	searchFileController := controllers.NewSearchFileController(*searchFileUseCase)

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
		searchFileController,
	)
}
