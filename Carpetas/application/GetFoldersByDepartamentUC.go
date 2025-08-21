package application

import (
	"VaultDoc-VD/Carpetas/domain/entities"
	"VaultDoc-VD/Carpetas/domain/repository"
)

type GetFoldersByDepartamentUseCase struct {
	repo repository.FoldersRepository
}

func NewGetFoldersByDepartamentUseCase(repo repository.FoldersRepository)*GetFoldersByDepartamentUseCase{
	return&GetFoldersByDepartamentUseCase{repo: repo}
}

func(uc *GetFoldersByDepartamentUseCase)Execute(departament string)([]entities.Folders, error){
	return uc.repo.GetFoldersByDepartament(departament)
}