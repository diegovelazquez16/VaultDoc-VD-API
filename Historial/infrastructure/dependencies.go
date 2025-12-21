// VaultDoc-VD-API/Historial/infrastructure/dependencies.go
package infrastructure

import (
	"VaultDoc-VD/Historial/application"
	"VaultDoc-VD/Historial/infrastructure/controllers"
	"VaultDoc-VD/Historial/infrastructure/repository"
	"VaultDoc-VD/Historial/infrastructure/routes"
	"VaultDoc-VD/core"

	"github.com/gin-gonic/gin"
)

func SetupDependencies(r *gin.Engine, dbPool *core.Conn_PostgreSQL){
	historyRepo := repository.NewHistoryPostgreSQLRepository(dbPool)

	saveActionsUseCase := application.NewSaveActionsUseCase(historyRepo)
	getHistoryUseCase := application.NewGetHistoryUseCase(historyRepo)
	getHistoryByIDUseCase := application.NewGetHistoryByIDUseCase(historyRepo)
	getAllHistoryUseCase := application.NewGetAllHistoryUseCase(historyRepo)
	getHistoryByUserUseCase := application.NewGetHistoryByUserUseCase(historyRepo)
	getHistoryByFolderUseCase := application.NewGetHistoryByFolderUseCase(historyRepo)
	getHistoryByFileUseCase := application.NewGetHistoryByFileUseCase(historyRepo)
	updateHistoryUseCase := application.NewUpdateHistoryUseCase(historyRepo)
	deleteHistoryUseCase := application.NewDeleteHistoryUseCase(historyRepo)

	saveActionsController := controllers.NewSaveActionController(*saveActionsUseCase)
	getHistoryController := controllers.NewGetHistoryController(*getHistoryUseCase)
	getHistoryByIDController := controllers.NewGetHistoryByIDController(*getHistoryByIDUseCase)
	getAllHistoryController := controllers.NewGetAllHistoryController(*getAllHistoryUseCase)
	getHistoryByUserController := controllers.NewGetHistoryByUserController(*getHistoryByUserUseCase)
	getHistoryByFolderController := controllers.NewGetHistoryByFolderController(*getHistoryByFolderUseCase)
	getHistoryByFileController := controllers.NewGetHistoryByFileController(*getHistoryByFileUseCase)
	updateHistoryController := controllers.NewUpdateHistoryController(*updateHistoryUseCase)
	deleteHistoryController := controllers.NewDeleteHistoryController(*deleteHistoryUseCase)
	

	routes.SetupHistoryRoutes(r, saveActionsController, getHistoryController, getHistoryByIDController, getAllHistoryController, getHistoryByUserController, getHistoryByFolderController, getHistoryByFileController, updateHistoryController, deleteHistoryController)
}