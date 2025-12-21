// VaultDoc-VD-API/Historial/helpers/history_helper.go
package helpers

import (
	"VaultDoc-VD/Historial/domain/entities"
	"VaultDoc-VD/Historial/domain/repository"
)

// HistoryHelper proporciona funciones auxiliares para registrar acciones en el historial
type HistoryHelper struct {
	repo repository.HistoryPostgreSQLRepo
}

func NewHistoryHelper(repo repository.HistoryPostgreSQLRepo) *HistoryHelper {
	return &HistoryHelper{repo: repo}
}

// Constantes para los tipos de movimientos
const (
	ActionUploadFile      = "Subió archivo"
	ActionDownloadFile    = "Bajó archivo"
	ActionModifyFile      = "Modificó información de un archivo"
	ActionDeleteFile      = "Eliminó archivo"
	ActionGrantAccess     = "Concedió acceso al archivo"
	ActionReceiveAccess   = "Se le otorgó acceso al archivo"
)

// LogFileUpload registra cuando un usuario sube un archivo
func (h *HistoryHelper) LogFileUpload(departamento string, folderID, fileID, userID int) error {
	record := entities.ReceiveHistory{
		Movimiento:   ActionUploadFile,
		Departamento: departamento,
		Id_folder:    folderID,
		Id_file:      fileID,
		Id_user:      userID,
	}
	return h.repo.SaveAction(record)
}

// LogFileDownload registra cuando un usuario descarga un archivo
func (h *HistoryHelper) LogFileDownload(departamento string, folderID, fileID, userID int) error {
	record := entities.ReceiveHistory{
		Movimiento:   ActionDownloadFile,
		Departamento: departamento,
		Id_folder:    folderID,
		Id_file:      fileID,
		Id_user:      userID,
	}
	return h.repo.SaveAction(record)
}

// LogFileModify registra cuando un usuario modifica la información de un archivo
func (h *HistoryHelper) LogFileModify(departamento string, folderID, fileID, userID int) error {
	record := entities.ReceiveHistory{
		Movimiento:   ActionModifyFile,
		Departamento: departamento,
		Id_folder:    folderID,
		Id_file:      fileID,
		Id_user:      userID,
	}
	return h.repo.SaveAction(record)
}

// LogFileDelete registra cuando un usuario elimina un archivo
func (h *HistoryHelper) LogFileDelete(departamento string, folderID, fileID, userID int) error {
	record := entities.ReceiveHistory{
		Movimiento:   ActionDeleteFile,
		Departamento: departamento,
		Id_folder:    folderID,
		Id_file:      fileID,
		Id_user:      userID,
	}
	return h.repo.SaveAction(record)
}

// LogGrantAccess registra cuando un usuario concede acceso a un archivo
func (h *HistoryHelper) LogGrantAccess(departamento string, folderID, fileID, userID int) error {
	record := entities.ReceiveHistory{
		Movimiento:   ActionGrantAccess,
		Departamento: departamento,
		Id_folder:    folderID,
		Id_file:      fileID,
		Id_user:      userID,
	}
	return h.repo.SaveAction(record)
}

// LogReceiveAccess registra cuando a un usuario se le otorga acceso a un archivo
func (h *HistoryHelper) LogReceiveAccess(departamento string, folderID, fileID, userID int) error {
	record := entities.ReceiveHistory{
		Movimiento:   ActionReceiveAccess,
		Departamento: departamento,
		Id_folder:    folderID,
		Id_file:      fileID,
		Id_user:      userID,
	}
	return h.repo.SaveAction(record)
}

// LogAction es un método genérico para registrar cualquier acción
func (h *HistoryHelper) LogAction(movimiento, departamento string, folderID, fileID, userID int) error {
	record := entities.ReceiveHistory{
		Movimiento:   movimiento,
		Departamento: departamento,
		Id_folder:    folderID,
		Id_file:      fileID,
		Id_user:      userID,
	}
	return h.repo.SaveAction(record)
}