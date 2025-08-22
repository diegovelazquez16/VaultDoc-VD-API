package entities

import (
	folderEntities "VaultDoc-VD/Carpetas/domain/entities"
	userEntities "VaultDoc-VD/Usuarios/domain/entities"
	fileEntities "VaultDoc-VD/Archivos/domain/entities"
)

type SendHistory struct {
	Id             int
	Movimiento     string
	Departamento   string
	Id_folder      folderEntities.Folders
	Id_file        fileEntities.Files
	Id_user        userEntities.User
	Fecha_registro string
}