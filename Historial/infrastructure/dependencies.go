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

	saveActionsController := controllers.NewSaveActionController(*saveActionsUseCase)
	getHistoryController := controllers.NewGetHistoryController(*getHistoryUseCase)

	routes.SetupHistoryRoutes(r, saveActionsController, getHistoryController)
}