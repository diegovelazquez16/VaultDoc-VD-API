package application

import (
	"VaultDoc-VD/Carpetas/domain/entities"
	"VaultDoc-VD/Carpetas/domain/repository"
)

type GetFolderByNameUseCase struct {
	repo repository.FoldersRepository
}

func NewGetFolderByName(repo repository.FoldersRepository)*GetFolderByNameUseCase{
	return&GetFolderByNameUseCase{repo: repo}
}

func(uc *GetFolderByNameUseCase)Execute(name string)([]entities.Folders, error){
	return uc.repo.GetFolderByName(name)
}