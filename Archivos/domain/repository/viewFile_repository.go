// Archivos/domain/repository/viewFile_repository.go
package repository

import (
	entities "VaultDoc-VD/Archivos/domain/entities"
	userEntities "VaultDoc-VD/Usuarios/domain/entities"
)

type ViewFileRepository interface {
	GrantPermission(viewFile entities.ViewFile) error
	RemovePermission(viewFile entities.ViewFile) error
	HasPermission(fileId, userId int) (bool, error)
	GetUsersWithViewPermission(fileId int) ([]userEntities.User, error)
	GetUsersWithoutViewPermission(fileId int) ([]userEntities.User, error)
}