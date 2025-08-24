// Archivos/application/get_file_by_id_usecase.go
package application

import (
	entities "VaultDoc-VD/Archivos/domain/entities"
	"VaultDoc-VD/Archivos/domain/repository"
)

type GetFileByIdUseCase struct {
	repo repository.FilesRepository
}

func NewGetFileByIdUseCase(repo repository.FilesRepository) *GetFileByIdUseCase {
	return &GetFileByIdUseCase{repo: repo}
}

func (uc *GetFileByIdUseCase) Execute(id int) (entities.Files, error) {
	return uc.repo.GetByID(id)
}