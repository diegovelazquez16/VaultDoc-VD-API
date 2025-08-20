package repository

import (
	entity "VaultDoc-VD/Archivos/domain/entities"
)

type FilesRepository interface {
	Create(file entity.Files) error
	GetByID(id int) (entity.Files, error)
	Update(file entity.Files) error
	Delete(id int) error
	GetAll() ([]entity.Files, error)
}
