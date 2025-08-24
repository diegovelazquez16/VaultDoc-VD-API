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
	getUsersByDepartmentUseCase := application.NewGetUsersByDepartmentUseCase(userRepo)
	getProfileUseCase := application.NewGetProfileUseCase(userRepo)
	updateUserUseCase := application.NewUpdateUserUseCase(userRepo, bcryptService)
	updateMyProfileUseCase := application.NewUpdateProfileUseCase(userRepo, bcryptService)
	deleteUserUseCase := application.NewDeleteUserUseCase(userRepo)
	loginUseCase := application.NewLoginUseCase(userRepo, tokenManager, bcryptService)

	// Inicializar controllers
	createUserController := controllers.NewCreateUserController(createUserUseCase)
	getAllUsersController := controllers.NewGetAllUsersController(getAllUsersUseCase)
	getUserByIdController := controllers.NewGetUserByIdUseController(getUserByIdUseCase)
	getUsersByDepartmentController := controllers.NewGetUsersByDepartmentController(getUsersByDepartmentUseCase)
	getProfileController := controllers.NewGetProfileController(getProfileUseCase)
	updateUserController := controllers.NewUpdateUserController(updateUserUseCase)
	updateMyProfileController := controllers.NewUpdateProfileController(updateMyProfileUseCase)
	deleteUserController := controllers.NewDeleteUserController(deleteUserUseCase)
	loginUserController := controllers.NewLoginUserController(loginUseCase)

	// Configurar rutas
	routes.SetupUserRoutes(
		r,
		createUserController,
		getAllUsersController,
		getUserByIdController,
		getUsersByDepartmentController,
		getProfileController,
		updateUserController,
		updateMyProfileController,
		deleteUserController,
		loginUserController,
	)
}