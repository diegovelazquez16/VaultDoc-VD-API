// Archivos/application/create_file_usecase.go
package application

import (
	entities "VaultDoc-VD/Archivos/domain/entities"
	"VaultDoc-VD/Archivos/domain/repository"
)

type CreateFileUseCase struct {
	repo repository.FilesRepository
}

func NewCreateFileUseCase(repo repository.FilesRepository) *CreateFileUseCase {
	return &CreateFileUseCase{repo: repo}
}

func (uc *CreateFileUseCase) Execute(file entities.Files) error {
	return uc.repo.Create(file)
}