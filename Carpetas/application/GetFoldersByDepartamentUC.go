//Carpetas/application/GetFolderByDepartamentUseCase.go
package application

import (
	"VaultDoc-VD/Carpetas/domain/entities"
	"VaultDoc-VD/Carpetas/domain/repository"
	"VaultDoc-VD/Carpetas/domain/services"
	"VaultDoc-VD/core"
)

type GetFoldersByDepartamentUseCase struct {
	repo           repository.FoldersRepository
	nextcloudClient *core.NextcloudClient
}

func NewGetFoldersByDepartamentUseCase(repo repository.FoldersRepository, cloudService services.CloudStorageService ) *GetFoldersByDepartamentUseCase {
	return &GetFoldersByDepartamentUseCase{
		repo:           repo,
		nextcloudClient: core.NewNextcloudClient(),
	}
}

func (uc *GetFoldersByDepartamentUseCase) Execute(departament string) ([]entities.Folders, error) {
	
	dbFolders, err := uc.repo.GetFoldersByDepartament(departament)
	if err != nil {
		return nil, err
	}


	return dbFolders, nil
}
