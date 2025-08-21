// Archivos/application/download_file_usecase.go
package application

import (
	"fmt"
	"os"
	"path/filepath"
)

type DownloadFileUseCase struct {
}

func NewDownloadFileUseCase() *DownloadFileUseCase {
	return &DownloadFileUseCase{}
}

func (uc *DownloadFileUseCase) Execute(dir, filename string) (string, error) {
	// 1. Verificar que FILES_DIR está configurado
	baseDir := os.Getenv("FILES_DIR")
	if baseDir == "" {
		return "", fmt.Errorf("FILES_DIR no está configurado")
	}

	// 2. Construir la ruta completa: FILES_DIR + dir + filename
	fullPath := filepath.Join(baseDir, dir, filename)
	
	// 3. Verificar que el archivo existe físicamente
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("archivo no encontrado en: %s", fullPath)
	}

	fmt.Printf("Preparando descarga del archivo: %s\n", fullPath)
	return fullPath, nil
}