// Archivos/domain/repository/changeFile_repository.go
package repository

import (
	entities "VaultDoc-VD/Archivos/domain/entities"
	userEntities "VaultDoc-VD/Usuarios/domain/entities"
)

type ChangeFileRepository interface {
	GrantPermission(changeFile entities.ChangeFile) error
	RemovePermission(changeFile entities.ChangeFile) error
	HasPermission(fileId, userId int) (bool, error)
	GetUsersWithChangePermission(fileId int) ([]userEntities.User, error)
	GetUsersWithoutChangePermission(fileId int) ([]userEntities.User, error)
	GetChangePermissionsOfAFolder(folderId int, userId int) ([]int, error)
}