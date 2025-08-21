// Archivos/domain/repository/files_repository.go
package repository

import (
	entities "VaultDoc-VD/Archivos/domain/entities"
)

type FilesRepository interface {
	Create(file entities.Files) error
	GetByID(id int) (entities.Files, error)
	Update(file entities.Files) error
	Delete(id int) error
	GetAll() ([]entities.Files, error)
}
