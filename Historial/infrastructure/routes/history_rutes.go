//VaultDoc-VD-API/Historial/infrastructure/routes/history_rutes.go
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
	getHistoryByID *controllers.GetHistoryByIDController,
	getAllHistory *controllers.GetAllHistoryController,
	getHistoryByUser *controllers.GetHistoryByUserController,
	getHistoryByFolder *controllers.GetHistoryByFolderController,
	getHistoryByFile *controllers.GetHistoryByFileController,
	updateHistory *controllers.UpdateHistoryController,
	deleteHistory *controllers.DeleteHistoryController,
){

    jwtSecret := os.Getenv("JWT_SECRET")

	g := r.Group("history")
	{
		g.POST("/", service.BossMiddleware(jwtSecret), saveActions.Execute)
		g.GET("/:departament", /*service.BossMiddleware(jwtSecret),*/ getHistory.Execute)
		g.GET("/id/:id", service.BossMiddleware(jwtSecret), getHistoryByID.Execute)
		g.GET("/all", service.BossMiddleware(jwtSecret), getAllHistory.Execute)
		g.GET("/user/:userID", service.BossMiddleware(jwtSecret), getHistoryByUser.Execute)
		g.GET("/folder/:folderID", service.BossMiddleware(jwtSecret), getHistoryByFolder.Execute)
		g.GET("/file/:fileID", service.BossMiddleware(jwtSecret), getHistoryByFile.Execute)
		g.PUT("/:id", service.BossMiddleware(jwtSecret), updateHistory.Execute)
		g.DELETE("/:id", service.BossMiddleware(jwtSecret), deleteHistory.Execute)
	}
}