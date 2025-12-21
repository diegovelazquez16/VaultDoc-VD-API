// VaultDoc-VD-API/Historial/application/GetHistoryByFolderUC.go
package application

import (
	"VaultDoc-VD/Historial/domain/entities"
	"VaultDoc-VD/Historial/domain/repository"
)

type GetHistoryByFolderUseCase struct {
	repo repository.HistoryPostgreSQLRepo
}

func NewGetHistoryByFolderUseCase(repo repository.HistoryPostgreSQLRepo) *GetHistoryByFolderUseCase {
	return &GetHistoryByFolderUseCase{repo: repo}
}

func (uc *GetHistoryByFolderUseCase) Execute(folderID int) ([]entities.SendHistory, error) {
	return uc.repo.GetHistoryByFolder(folderID)
}
