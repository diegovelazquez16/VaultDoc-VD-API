// VaultDoc-VD-API/Historial/application/DeleteHistoryUC.go
package application

import (
	"VaultDoc-VD/Historial/domain/repository"
)

type DeleteHistoryUseCase struct {
	repo repository.HistoryPostgreSQLRepo
}

func NewDeleteHistoryUseCase(repo repository.HistoryPostgreSQLRepo) *DeleteHistoryUseCase {
	return &DeleteHistoryUseCase{repo: repo}
}

func (uc *DeleteHistoryUseCase) Execute(id int) error {
	return uc.repo.DeleteHistory(id)
}