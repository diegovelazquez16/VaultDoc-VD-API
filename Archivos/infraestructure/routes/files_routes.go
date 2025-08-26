// Archivos/infrastructure/routes/files_routes.go
package routes

import (
	"VaultDoc-VD/Archivos/infraestructure/controllers"
	"VaultDoc-VD/Middlewares"
	"os"

	"github.com/gin-gonic/gin"
)

func SetupFilesRoutes(
	r *gin.Engine,
	createFileController *controllers.CreateFileController,
	getFileByIdController *controllers.GetFileByIdController,
	getAllFilesController *controllers.GetAllFilesController,
	getFilesByFolderController *controllers.GetFilesByFolderController,
	updateFileController *controllers.UpdateFileController,
	deleteFileController *controllers.DeleteFileController,
	downloadFileController *controllers.DownloadFileController,
	grantChangePermissionController *controllers.GrantChangePermissionController,
	removeChangePermissionController *controllers.RemoveChangePermissionController,
	grantViewPermissionController *controllers.GrantViewPermissionController,
	removeViewPermissionController *controllers.RemoveViewPermissionController,
	checkPermissionsController *controllers.CheckPermissionsController,
) {

    jwtSecret := os.Getenv("JWT_SECRET")

	filesGroup := r.Group("files")
	{
		// cualquiera con autenticación
		filesGroup.POST("/", service.AuthMiddleware(jwtSecret), createFileController.Execute)
		filesGroup.GET("/:id", service.AuthMiddleware(jwtSecret), getFileByIdController.Execute)
		filesGroup.GET("/folder/:folderId", service.AuthMiddleware(jwtSecret), getFilesByFolderController.Execute)
		filesGroup.PUT("/:id/:id_user", service.AuthMiddleware(jwtSecret), updateFileController.Execute)
		filesGroup.DELETE("/:id/:id_user", service.AuthMiddleware(jwtSecret), deleteFileController.Execute)
		filesGroup.GET("/download/:id/:id_user", service.AuthMiddleware(jwtSecret), downloadFileController.Execute)
        
		// solo el jefe de departamento:
		// Permisos de edición
		filesGroup.POST("/permissions/change/:id_user", service.BossMiddleware(jwtSecret), grantChangePermissionController.Execute)
		filesGroup.DELETE("/permissions/change", service.BossMiddleware(jwtSecret), removeChangePermissionController.Execute)
		// Permisos de visualización
		filesGroup.POST("/permissions/view/:id_user", service.BossMiddleware(jwtSecret), grantViewPermissionController.Execute)
		filesGroup.DELETE("/permissions/view", service.BossMiddleware(jwtSecret), removeViewPermissionController.Execute)
		// Verificar permisos
		filesGroup.GET("/permissions/:fileId/:userId", service.BossMiddleware(jwtSecret), checkPermissionsController.Execute)

		filesGroup.GET("/", service.AdminMiddleware(jwtSecret), getAllFilesController.Execute)

	}
}