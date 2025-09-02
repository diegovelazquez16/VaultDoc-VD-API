package application

import (
		entities "VaultDoc-VD/Archivos/domain/entities"
		"VaultDoc-VD/Archivos/domain/repository"
	)

type SearchFileUseCase struct {
	repo repository.FilesRepository
}

func NewSearchFileUseCase(repo repository.FilesRepository)*SearchFileUseCase{
	return&SearchFileUseCase{repo: repo}
}

func (uc *SearchFileUseCase)Execute(name string)([]entities.Files, error) {
	return uc.repo.SearchFile(name);
}