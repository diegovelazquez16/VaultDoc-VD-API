// Carpetas/domain/services/cloud_storage_service.go
package services


type CloudStorageService interface {
	CreateFolder(folderPath string) error
	FolderExists(folderPath string) (bool, error)
	ListFolders(basePath string) ([]string, error)
}