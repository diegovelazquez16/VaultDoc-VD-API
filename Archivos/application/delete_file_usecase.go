// Archivos/application/delete_file_usecase.go
package application

import (
	"VaultDoc-VD/Archivos/domain/repository"
)

type DeleteFileUseCase struct {
	repo repository.FilesRepository
}

func NewDeleteFileUseCase(repo repository.FilesRepository) *DeleteFileUseCase {
	return &DeleteFileUseCase{repo: repo}
}

func (uc *DeleteFileUseCase) Execute(id int) error {
	return uc.repo.Delete(id)
}