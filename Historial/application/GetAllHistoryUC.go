// VaultDoc-VD-API/Historial/application/GetAllHistoryUC.go
package application

import (
	"VaultDoc-VD/Historial/domain/entities"
	"VaultDoc-VD/Historial/domain/repository"
)

type GetAllHistoryUseCase struct {
	repo repository.HistoryPostgreSQLRepo
}

func NewGetAllHistoryUseCase(repo repository.HistoryPostgreSQLRepo) *GetAllHistoryUseCase {
	return &GetAllHistoryUseCase{repo: repo}
}

func (uc *GetAllHistoryUseCase) Execute() ([]entities.SendHistory, error) {
	return uc.repo.GetAllHistory()
}