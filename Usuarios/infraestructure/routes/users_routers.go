// Usuarios/routes/users_routers.go
package routes

import (
	"VaultDoc-VD/Usuarios/infraestructure/controllers"
	"VaultDoc-VD/Middlewares"
	"os"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine, createUserController *controllers.CreateUserController,
	getUsersController *controllers.GetAllUsersController,
	getUsersControllerById *controllers.GetUserByIdController,
	updateUserController *controllers.UpdateUserController,
	updateProfileController *controllers.UpdateProfileController,
	deleteUserController *controllers.DeleteUserController,
	loginUserController *controllers.LoginUserController,
) {
	jwtSecret := os.Getenv("JWT_SECRET")

	// solo el admin
	r.POST("/users", service.AdminMiddleware(jwtSecret), createUserController.Execute)
	r.DELETE("/users/:id", service.AdminMiddleware(jwtSecret), deleteUserController.Execute)
	r.GET("/users", service.AdminMiddleware(jwtSecret), getUsersController.Execute)
	r.GET("/users/:id", service.AdminMiddleware(jwtSecret), getUsersControllerById.Execute)
	r.PUT("/users/:id", service.AdminMiddleware(jwtSecret), updateUserController.Execute)

    // solo el jefe de departamento

	// cualquier usuario autenticado
	r.PUT("/users/profile", service.AuthMiddleware(jwtSecret), updateProfileController.Execute)

	// sin autenticación
	r.POST("/users/login", loginUserController.Execute)

}
