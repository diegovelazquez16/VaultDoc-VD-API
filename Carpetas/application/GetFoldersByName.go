//Carpetas/application/GetFoldersByName.go
package application

import (
	"VaultDoc-VD/Carpetas/domain/entities"
	"VaultDoc-VD/Carpetas/domain/repository"
	"VaultDoc-VD/Carpetas/domain/services"
	"VaultDoc-VD/core"
)

type GetFolderByNameUseCase struct {
	repo           repository.FoldersRepository
	nextcloudClient *core.NextcloudClient
}

func NewGetFolderByName(repo repository.FoldersRepository, cloudService services.CloudStorageService) *GetFolderByNameUseCase {
	return &GetFolderByNameUseCase{
		repo:           repo,
		nextcloudClient: core.NewNextcloudClient(),
	}
}

func (uc *GetFolderByNameUseCase) Execute(name string) ([]entities.Folders, error) {
	// Obtener carpetas de la base de datos por nombre
	dbFolders, err := uc.repo.GetFolderByName(name)
	if err != nil {
		return nil, err
	}



	return dbFolders, nil
}