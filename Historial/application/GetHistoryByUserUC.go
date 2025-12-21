// VaultDoc-VD-API/Historial/application/GetHistoryByUserUC.go
package application

import (
	"VaultDoc-VD/Historial/domain/entities"
	"VaultDoc-VD/Historial/domain/repository"
)

type GetHistoryByUserUseCase struct {
	repo repository.HistoryPostgreSQLRepo
}

func NewGetHistoryByUserUseCase(repo repository.HistoryPostgreSQLRepo) *GetHistoryByUserUseCase {
	return &GetHistoryByUserUseCase{repo: repo}
}

func (uc *GetHistoryByUserUseCase) Execute(userID int) ([]entities.SendHistory, error) {
	return uc.repo.GetHistoryByUser(userID)
}