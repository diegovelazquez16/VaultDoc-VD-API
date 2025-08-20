// Usuarios/routes/users_routers.go
package routes

import (
	"VaultDoc-VD/Usuarios/infraestructure/controllers"
	_ "VaultDoc-VD/Usuarios/infraestructure/services"
	_ "os"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine, createUserController *controllers.CreateUserController,
	getUsersController *controllers.GetAllUsersController,
	getUsersControllerById *controllers.GetUserByIdController,
	updateUserController *controllers.UpdateUserController,
	deleteUserController *controllers.DeleteUserController,
	loginUserController *controllers.LoginUserController,
) {
	//jwtSecret := os.Getenv("JWT_SECRET")

	r.POST("/users", createUserController.Execute)
	r.GET("/users", getUsersController.Execute)
	r.GET("/users/:id", getUsersControllerById.Execute)
	r.PUT("/users/:id", updateUserController.Execute)
	r.DELETE("/users/:id", deleteUserController.Execute)
	r.POST("/users/login", loginUserController.Execute)

}
