// Archivos/domain/services/file_storage_service.go
package services

import (
	"mime/multipart"
)

type FileStorageService interface {
	UploadFile(folderPath string, fileName string, fileHeader *multipart.FileHeader) (string, error)
	UploadFileFromBytes(folderPath string, fileName string, fileContent []byte) (string, error)
	DownloadFile(folderPath string, fileName string) ([]byte, error)
	DeleteFile(folderPath string, fileName string) error
	FileExists(folderPath string, fileName string) (bool, error)
	GetFileInfo(folderPath string, fileName string) (*FileInfo, error)
}

type FileInfo struct {
	Name         string
	Size         int64
	LastModified string
	ContentType  string
}