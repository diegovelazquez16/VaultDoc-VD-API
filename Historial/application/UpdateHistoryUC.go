// VaultDoc-VD-API/Historial/application/UpdateHistoryUC.go
package application

import (
	"VaultDoc-VD/Historial/domain/entities"
	"VaultDoc-VD/Historial/domain/repository"
)

type UpdateHistoryUseCase struct {
	repo repository.HistoryPostgreSQLRepo
}

func NewUpdateHistoryUseCase(repo repository.HistoryPostgreSQLRepo) *UpdateHistoryUseCase {
	return &UpdateHistoryUseCase{repo: repo}
}

func (uc *UpdateHistoryUseCase) Execute(id int, record entities.ReceiveHistory) error {
	return uc.repo.UpdateHistory(id, record)
}