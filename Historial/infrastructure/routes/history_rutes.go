package routes

import (
	"VaultDoc-VD/Historial/infrastructure/controllers"

	"github.com/gin-gonic/gin"
)

func SetupHistoryRoutes(
	r *gin.Engine,
	saveActions *controllers.SaveActionController,
	getHistory *controllers.GetHistoryController,
){
	r.Group("history")
	r.POST("/", saveActions.Execute)
	r.GET("/:departament", getHistory.Execute)
}