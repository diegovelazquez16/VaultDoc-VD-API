package application

import (
	"VaultDoc-VD/Historial/domain/entities"
	"VaultDoc-VD/Historial/domain/repository"
)

type GetHistoryUseCase struct {
	repo repository.HistoryPostgreSQLRepo
}

func NewGetHistoryUseCase(repo repository.HistoryPostgreSQLRepo)*GetHistoryUseCase{
	return&GetHistoryUseCase{repo: repo}
}

func(uc *GetHistoryUseCase)Execute(departament string)([]entities.SendHistory, error){
	return uc.repo.GetHistory(departament)
}