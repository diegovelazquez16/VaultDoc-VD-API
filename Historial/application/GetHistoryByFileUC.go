// VaultDoc-VD-API/Historial/application/GetHistoryByFileUC.go
package application

import (
	"VaultDoc-VD/Historial/domain/entities"
	"VaultDoc-VD/Historial/domain/repository"
)

type GetHistoryByFileUseCase struct {
	repo repository.HistoryPostgreSQLRepo
}

func NewGetHistoryByFileUseCase(repo repository.HistoryPostgreSQLRepo) *GetHistoryByFileUseCase {
	return &GetHistoryByFileUseCase{repo: repo}
}

func (uc *GetHistoryByFileUseCase) Execute(fileID int) ([]entities.SendHistory, error) {
	return uc.repo.GetHistoryByFile(fileID)
}