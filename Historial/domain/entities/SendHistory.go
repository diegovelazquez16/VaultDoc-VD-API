// VaultDoc-VD-API/Historial/domain/entities/SendHistory.go
package entities

import (
	folderEntities "VaultDoc-VD/Carpetas/domain/entities"
	userEntities "VaultDoc-VD/Usuarios/domain/entities"
	fileEntities "VaultDoc-VD/Archivos/domain/entities"
)

type SendHistory struct {
	Id             int                      `json:"id"`
	Movimiento     string                   `json:"movimiento"`
	Departamento   string                   `json:"departamento"`
	Id_folder      folderEntities.Folders   `json:"folder"`
	Id_file        fileEntities.Files       `json:"file"`
	Id_user        userEntities.User        `json:"user"`
	
	FolderName     string                   `json:"folder_name_backup,omitempty"`
	FileName       string                   `json:"file_name_backup,omitempty"`
	UserName       string                   `json:"user_name_backup,omitempty"`
	
	Fecha_registro string                   `json:"fecha_registro"`
}