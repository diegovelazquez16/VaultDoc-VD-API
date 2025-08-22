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
	g := r.Group("history")
	{
		g.POST("/", saveActions.Execute)
		g.GET("/:departament", getHistory.Execute)
	}
}