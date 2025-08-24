package routes

import (
	"VaultDoc-VD/Historial/infrastructure/controllers"
	"VaultDoc-VD/Middlewares"
	"os"

	"github.com/gin-gonic/gin"
)

func SetupHistoryRoutes(
	r *gin.Engine,
	saveActions *controllers.SaveActionController,
	getHistory *controllers.GetHistoryController,
){

    jwtSecret := os.Getenv("JWT_SECRET")

	g := r.Group("history")
	{
		g.POST("/", service.BossMiddleware(jwtSecret), saveActions.Execute)
		g.GET("/:departament", service.BossMiddleware(jwtSecret), getHistory.Execute)
	}
}