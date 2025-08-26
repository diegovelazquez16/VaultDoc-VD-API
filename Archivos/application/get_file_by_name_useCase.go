package application

import (
	"VaultDoc-VD/Archivos/domain/repository"
	"VaultDoc-VD/Archivos/domain/entities"
)

type GetFileByNameUseCase struct {
	repo repository.FilesRepository
}

func NewGetFileByNameUseCase(repo repository.FilesRepository)*GetFileByNameUseCase{
	return&GetFileByNameUseCase{repo: repo}
}

func (uc *GetFileByNameUseCase) Execute(name string) (domain.Files, error) {
	return uc.repo.GetByName(name)
}