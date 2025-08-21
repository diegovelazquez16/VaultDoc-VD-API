package main

import (
	"log"

	usuariosInfra "VaultDoc-VD/Usuarios/infraestructure"
	archivosInfra "VaultDoc-VD/Archivos/infraestructure"
	"VaultDoc-VD/Carpetas/infrastructure"

	"VaultDoc-VD/core"

	"github.com/gin-gonic/gin"
)

func main() {
	dbPool := core.GetDBPool()
	if dbPool.Err != "" {
		log.Fatalf("Error al conectar con la base de datos: %s", dbPool.Err)
	}
	defer dbPool.DB.Close()

	r := gin.Default()

	r.Use(core.SetupCORS())

	usuariosInfra.SetupDependencies(r, dbPool)
	archivosInfra.SetupDependencies(r, dbPool)
	infrastructure.SetupDependenciesFolders(r, dbPool)

	log.Println("Servidor iniciado en puerto 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}