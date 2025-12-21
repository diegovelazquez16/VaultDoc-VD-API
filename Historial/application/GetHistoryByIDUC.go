// VaultDoc-VD-API/Historial/application/GetHistoryByIDUC.go
package application

import (
	"VaultDoc-VD/Historial/domain/entities"
	"VaultDoc-VD/Historial/domain/repository"
)

type GetHistoryByIDUseCase struct {
	repo repository.HistoryPostgreSQLRepo
}

func NewGetHistoryByIDUseCase(repo repository.HistoryPostgreSQLRepo) *GetHistoryByIDUseCase {
	return &GetHistoryByIDUseCase{repo: repo}
}

func (uc *GetHistoryByIDUseCase) Execute(id int) (*entities.ReceiveHistory, error) {
	return uc.repo.GetHistoryByID(id)
}