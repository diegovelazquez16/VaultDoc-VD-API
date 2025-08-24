// Archivos/application/update_file_usecase.go
package application

import (
	entities "VaultDoc-VD/Archivos/domain/entities"
	"VaultDoc-VD/Archivos/domain/repository"
)

type UpdateFileUseCase struct {
	repo repository.FilesRepository
}

func NewUpdateFileUseCase(repo repository.FilesRepository) *UpdateFileUseCase {
	return &UpdateFileUseCase{repo: repo}
}

func (uc *UpdateFileUseCase) Execute(file entities.Files) error {
	return uc.repo.Update(file)
}