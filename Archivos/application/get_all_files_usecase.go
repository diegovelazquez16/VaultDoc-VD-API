// Archivos/application/get_all_files_usecase.go
package application

import (
	entities "VaultDoc-VD/Archivos/domain/entities"
	"VaultDoc-VD/Archivos/domain/repository"
)

type GetAllFilesUseCase struct {
	repo repository.FilesRepository
}

func NewGetAllFilesUseCase(repo repository.FilesRepository) *GetAllFilesUseCase {
	return &GetAllFilesUseCase{repo: repo}
}

func (uc *GetAllFilesUseCase) Execute() ([]entities.Files, error) {
	return uc.repo.GetAll()
}