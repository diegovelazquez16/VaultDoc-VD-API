// VaultDoc-VD-API/Historial/domain/repository/history_postgresql_repo.go
package repository

import "VaultDoc-VD/Historial/domain/entities"

type HistoryPostgreSQLRepo interface {
	SaveAction(history entities.ReceiveHistory) error
	GetHistory(departament string) ([]entities.SendHistory, error)
	GetHistoryByID(id int) (*entities.ReceiveHistory, error)
	UpdateHistory(id int, history entities.ReceiveHistory) error
	DeleteHistory(id int) error
	GetAllHistory() ([]entities.SendHistory, error)
	GetHistoryByUser(userID int) ([]entities.SendHistory, error)
	GetHistoryByFolder(folderID int) ([]entities.SendHistory, error)
	GetHistoryByFile(fileID int) ([]entities.SendHistory, error)
}