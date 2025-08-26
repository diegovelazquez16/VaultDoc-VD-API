// Archivos/application/download_file_usecase.go
package application

import (
	"fmt"
	"strings"
	"VaultDoc-VD/Archivos/domain/repository"
	"VaultDoc-VD/Archivos/domain/services"
)

type DownloadFileUseCase struct {
	repo               repository.FilesRepository
	fileStorageService services.FileStorageService
}

func NewDownloadFileUseCase(repo repository.FilesRepository, fileStorageService services.FileStorageService) *DownloadFileUseCase {
	return &DownloadFileUseCase{
		repo:               repo,
		fileStorageService: fileStorageService,
	}
}

func (uc *DownloadFileUseCase) Execute(id int) ([]byte, string, error) {
	// 1. Obtener información del archivo de la BD
	file, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, "", fmt.Errorf("archivo no encontrado en base de datos: %v", err)
	}

	// 2. Extraer folderPath y fileName del directorio almacenado
	folderPath, fileName := uc.parsePath(file.Directorio)

	// 3. Descargar archivo de Nextcloud
	content, err := uc.fileStorageService.DownloadFile(folderPath, fileName)
	if err != nil {
		return nil, "", fmt.Errorf("error al descargar archivo de Nextcloud: %v", err)
	}

	fmt.Printf("Preparando descarga del archivo: %s\n", file.Directorio)
	return content, fileName, nil
}


func (uc *DownloadFileUseCase) parsePath(fullPath string) (string, string) {
	parts := strings.Split(fullPath, "/")
	if len(parts) <= 1 {
		return "", fullPath
	}
	
	fileName := parts[len(parts)-1]
	folderPath := strings.Join(parts[:len(parts)-1], "/")
	return folderPath, fileName
}

