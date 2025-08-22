package repository

import "VaultDoc-VD/Historial/domain/entities"

type HistoryPostgreSQLRepo interface {
	SaveAction(history entities.ReceiveHistory) (error)
	GetHistory(departament string) ([]entities.SendHistory, error)
	GetHistoryByID(id int) (*entities.ReceiveHistory, error)
}