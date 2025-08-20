// Usuarios/infraestructure/dependencies.go
package infraestructure

import (
	"VaultDoc-VD/Usuarios/application"
	"VaultDoc-VD/Usuarios/infraestructure/controllers"
	"VaultDoc-VD/Usuarios/infraestructure/repository"
	"VaultDoc-VD/Usuarios/infraestructure/services"
	"VaultDoc-VD/Usuarios/infraestructure/routes"
	"VaultDoc-VD/core"

	"github.com/gin-gonic/gin"
)

func SetupDependencies(r *gin.Engine, dbPool *core.Conn_PostgreSQL) {
	// Inicializar servicios
	bcryptService := service.InitBcryptService()
	tokenManager := service.InitTokenManager()

	// Inicializar repository
	userRepo := repository.NewUserPostgreSQLRepository(dbPool)

	// Inicializar use cases
	createUserUseCase := application.NewCreateUserUseCase(userRepo, bcryptService)
	getAllUsersUseCase := application.NewGetUsersUseCase(userRepo)
	getUserByIdUseCase := application.NewGetUserByIdUseCase(userRepo)
	updateUserUseCase := application.NewUpdateUserUseCase(userRepo, bcryptService)
	deleteUserUseCase := application.NewDeleteUserUseCase(userRepo)
	loginUseCase := application.NewLoginUseCase(userRepo, tokenManager, bcryptService)

	// Inicializar controllers
	createUserController := controllers.NewCreateUserController(createUserUseCase)
	getAllUsersController := controllers.NewGetAllUsersController(getAllUsersUseCase)
	getUserByIdController := controllers.NewGetUserByIdUseController(getUserByIdUseCase)
	updateUserController := controllers.NewUpdateUserController(updateUserUseCase)
	deleteUserController := controllers.NewDeleteUserController(deleteUserUseCase)
	loginUserController := controllers.NewLoginUserController(loginUseCase)

	// Configurar rutas
	routes.SetupUserRoutes(
		r,
		createUserController,
		getAllUsersController,
		getUserByIdController,
		updateUserController,
		deleteUserController,
		loginUserController,
	)
}