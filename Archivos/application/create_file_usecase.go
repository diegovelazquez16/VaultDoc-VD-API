// Archivos/application/create_file_usecase.go
package application

import (
	"fmt"
	"os"
	_"path/filepath"
	
	entities "VaultDoc-VD/Archivos/domain/entities"
	"VaultDoc-VD/Archivos/domain/repository"
)

type CreateFileUseCase struct {
	repo repository.FilesRepository
}

func NewCreateFileUseCase(repo repository.FilesRepository) *CreateFileUseCase {
	return &CreateFileUseCase{repo: repo}
}

func (uc *CreateFileUseCase) Execute(file entities.Files) error {
	// Validaciones de negocio
	if err := uc.validateFile(file); err != nil {
		return fmt.Errorf("validación fallida: %v", err)
	}

	// Verificar que el directorio base existe
	if err := uc.validateBaseDirectory(); err != nil {
		return fmt.Errorf("error de directorio: %v", err)
	}

	// Crear el registro en la base de datos
	if err := uc.repo.Create(file); err != nil {
		return fmt.Errorf("error al crear archivo en base de datos: %v", err)
	}

	fmt.Printf("Archivo creado exitosamente: %s\n", file.Nombre)
	return nil
}

// validateFile valida los datos del archivo
func (uc *CreateFileUseCase) validateFile(file entities.Files) error {
	if file.Departamento == "" {
		return fmt.Errorf("departamento es requerido")
	}
	
	if file.Nombre == "" {
		return fmt.Errorf("nombre es requerido")
	}
	
	if file.Folio == "" {
		return fmt.Errorf("folio es requerido")
	}
	
	if file.Id_Folder <= 0 {
		return fmt.Errorf("id_folder debe ser mayor a 0")
	}
	
	if file.Id_Uploader <= 0 {
		return fmt.Errorf("id_uploader debe ser mayor a 0")
	}
	
	if file.Tamano < 0 {
		return fmt.Errorf("tamaño no puede ser negativo")
	}

	return nil
}

// validateBaseDirectory verifica que el directorio base existe
func (uc *CreateFileUseCase) validateBaseDirectory() error {
	baseDir := os.Getenv("FILES_DIR")
	if baseDir == "" {
		return fmt.Errorf("FILES_DIR no está configurado")
	}
	
	// Verificar que existe y es un directorio
	info, err := os.Stat(baseDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("directorio base no existe: %s", baseDir)
	}
	
	if err != nil {
		return fmt.Errorf("error al acceder al directorio base: %v", err)
	}
	
	if !info.IsDir() {
		return fmt.Errorf("FILES_DIR no es un directorio válido: %s", baseDir)
	}

	return nil
}