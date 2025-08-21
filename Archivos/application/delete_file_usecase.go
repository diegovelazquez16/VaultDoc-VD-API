// Archivos/application/delete_file_usecase.go
package application

import (
	"fmt"
	"os"
	"path/filepath"
	
	"VaultDoc-VD/Archivos/domain/repository"
)

type DeleteFileUseCase struct {
	repo repository.FilesRepository
}

func NewDeleteFileUseCase(repo repository.FilesRepository) *DeleteFileUseCase {
	return &DeleteFileUseCase{repo: repo}
}

func (uc *DeleteFileUseCase) Execute(id int) error {
	// 1. Obtener información del archivo antes de eliminarlo
	file, err := uc.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("archivo no encontrado en base de datos: %v", err)
	}

	// 2. Construir la ruta completa del archivo físico
	baseDir := os.Getenv("FILES_DIR")
	if baseDir == "" {
		return fmt.Errorf("FILES_DIR no está configurado")
	}

	fullPath := filepath.Join(baseDir, file.Directorio)
	
	// 3. Eliminar el archivo físico si existe
	if _, err := os.Stat(fullPath); err == nil {
		// El archivo existe, eliminarlo
		if err := os.Remove(fullPath); err != nil {
			fmt.Printf("Advertencia: No se pudo eliminar archivo físico %s: %v\n", fullPath, err)
			// Continuar con la eliminación de la BD aunque falle el archivo físico
		} else {
			fmt.Printf("Archivo físico eliminado: %s\n", fullPath)
		}
	} else {
		fmt.Printf("Archivo físico no encontrado (puede que ya haya sido eliminado): %s\n", fullPath)
	}

	// 4. Eliminar el registro de la base de datos
	if err := uc.repo.Delete(id); err != nil {
		return fmt.Errorf("error al eliminar registro de base de datos: %v", err)
	}

	fmt.Printf("Archivo eliminado exitosamente: ID %d, Nombre: %s\n", id, file.Nombre)
	return nil
}