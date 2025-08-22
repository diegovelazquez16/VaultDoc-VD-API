// Archivos/domain/repository/viewFile_repository.go
package repository

import entities "VaultDoc-VD/Archivos/domain/entities"

type ViewFileRepository interface {
	GrantPermission(viewFile entities.ViewFile) error
	RemovePermission(viewFile entities.ViewFile) error
	HasPermission(fileId, userId int) (bool, error)
}