// Archivos/infrastructure/routes/files_routes.go
package routes

import (
	"VaultDoc-VD/Archivos/infraestructure/controllers"
	"github.com/gin-gonic/gin"
)

func SetupFilesRoutes(
	r *gin.Engine,
	createFileController *controllers.CreateFileController,
	getFileByIdController *controllers.GetFileByIdController,
	getAllFilesController *controllers.GetAllFilesController,
	updateFileController *controllers.UpdateFileController,
	deleteFileController *controllers.DeleteFileController,
	downloadFileController *controllers.DownloadFileController,
	grantChangePermissionController *controllers.GrantChangePermissionController,
	removeChangePermissionController *controllers.RemoveChangePermissionController,
	grantViewPermissionController *controllers.GrantViewPermissionController,
	removeViewPermissionController *controllers.RemoveViewPermissionController,
	checkPermissionsController *controllers.CheckPermissionsController,
) {
	filesGroup := r.Group("files")
	{
		// CRUD de archivos
		filesGroup.POST("/", createFileController.Execute)
		filesGroup.GET("/:id", getFileByIdController.Execute)
		filesGroup.GET("/", getAllFilesController.Execute)
		filesGroup.PUT("/:id", updateFileController.Execute)
		filesGroup.DELETE("/:id", deleteFileController.Execute)
		
		// Descarga de archivos por ruta
		filesGroup.GET("/download/*dir", downloadFileController.Execute)
		
		// Permisos de edición
		filesGroup.POST("/permissions/change", grantChangePermissionController.Execute)
		filesGroup.DELETE("/permissions/change", removeChangePermissionController.Execute)
		
		// Permisos de visualización
		filesGroup.POST("/permissions/view", grantViewPermissionController.Execute)
		filesGroup.DELETE("/permissions/view", removeViewPermissionController.Execute)
		
		// Verificar permisos
		filesGroup.GET("/permissions/:fileId/:userId", checkPermissionsController.Execute)
	}
}