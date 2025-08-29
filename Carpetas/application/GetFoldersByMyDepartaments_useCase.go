package application

import (
	"fmt"
	"VaultDoc-VD/Carpetas/domain/entities"
	"VaultDoc-VD/Carpetas/domain/repository"
	"VaultDoc-VD/Carpetas/domain/services"
	"VaultDoc-VD/core"
)

type GetFolderByMyDepartamentUseCase struct {
	repo repository.FoldersRepository
	nextcloudClient *core.NextcloudClient
}

func NewGetFoldersByMyDepartamentUseCase(repo repository.FoldersRepository, cloudService services.CloudStorageService) *GetFolderByMyDepartamentUseCase{
	return &GetFolderByMyDepartamentUseCase{
		repo: repo,
		nextcloudClient: core.NewNextcloudClient(),
	}
}

func (uc *GetFolderByMyDepartamentUseCase) Execute(myDepartament string) ([]entities.Folders, error){
	if myDepartament == "" {
		return nil, fmt.Errorf("Departamento no pueda estar vacio")
	}
	folders, err := uc.repo.GetFoldersByMyDepartament(myDepartament)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener los folders de su departamento: %w", err)
	}

	return folders, nil
}