// Archivos/application/delete_file_usecase.go (Modificado)
package application

import (
	"fmt"
	"strings"
	"VaultDoc-VD/Archivos/domain/repository"
	"VaultDoc-VD/Archivos/domain/services"
)

type DeleteFileUseCase struct {
	repo               repository.FilesRepository
	fileStorageService services.FileStorageService
}

func NewDeleteFileUseCase(repo repository.FilesRepository, fileStorageService services.FileStorageService) *DeleteFileUseCase {
	return &DeleteFileUseCase{
		repo:               repo,
		fileStorageService: fileStorageService,
	}
}

func (uc *DeleteFileUseCase) Execute(id int) error {
	// 1. Obtener información del archivo antes de eliminarlo
	file, err := uc.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("archivo no encontrado en base de datos: %v", err)
	}

	// 2. Extraer folderPath y fileName del directorio almacenado
	folderPath, fileName := uc.parsePath(file.Directorio)

	// 3. Eliminar el archivo de Nextcloud
	if err := uc.fileStorageService.DeleteFile(folderPath, fileName); err != nil {
		fmt.Printf("Warning: No se pudo eliminar archivo de Nextcloud: %v\n", err)
		// Continuar con la eliminación de la BD aunque falle Nextcloud
	} else {
		fmt.Printf("Archivo eliminado de Nextcloud: %s\n", file.Directorio)
	}

	// 4. Eliminar el registro de la base de datos
	if err := uc.repo.Delete(id); err != nil {
		return fmt.Errorf("error al eliminar registro de base de datos: %v", err)
	}

	fmt.Printf("Archivo eliminado exitosamente: ID %d, Nombre: %s\n", id, file.Nombre)
	return nil
}

// Para rutas completas
func (uc *DeleteFileUseCase) parsePath(fullPath string) (string, string) {
	parts := strings.Split(fullPath, "/")
	if len(parts) <= 1 {
		return "", fullPath
	}
	
	fileName := parts[len(parts)-1]
	folderPath := strings.Join(parts[:len(parts)-1], "/")
	return folderPath, fileName
}
