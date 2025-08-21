// Archivos/domain/repository/changeFile_repository.go
package repository

import entities"VaultDoc-VD/Archivos/domain/entities"

type ChangeFileRepository interface {
	GrantPermission(changeFile entities.ChangeFile) error
	RemovePermission(changeFile entities.ChangeFile) error
}