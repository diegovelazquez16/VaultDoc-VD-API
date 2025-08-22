package application

import (
	"VaultDoc-VD/Historial/domain/entities"
	"VaultDoc-VD/Historial/domain/repository"
)

type SaveActionUseCase struct {
	repo repository.HistoryPostgreSQLRepo
}

func NewSaveActionsUseCase(repo repository.HistoryPostgreSQLRepo)*SaveActionUseCase{
	return&SaveActionUseCase{repo: repo}
}

func(uc *SaveActionUseCase)Execute(record entities.ReceiveHistory)(*entities.ReceiveHistory, error){
	id, err := uc.repo.SaveAction(record)
	if err != nil {
		return nil, err
	}
	return uc.repo.GetHistoryByID(id)
}