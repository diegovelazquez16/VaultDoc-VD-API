// Archivos/application/create_file_usecase.go (Corrección)
package application

import (
	"fmt"
	"mime/multipart"
	entities "VaultDoc-VD/Archivos/domain/entities"
	"VaultDoc-VD/Archivos/domain/repository"
	"VaultDoc-VD/Archivos/domain/services"
)

type CreateFileUseCase struct {
	repo               repository.FilesRepository
	fileStorageService services.FileStorageService
}

func NewCreateFileUseCase(repo repository.FilesRepository, fileStorageService services.FileStorageService) *CreateFileUseCase {
	return &CreateFileUseCase{
		repo:               repo,
		fileStorageService: fileStorageService,
	}
}

func (uc *CreateFileUseCase) Execute(file entities.Files, fileHeader *multipart.FileHeader) error {
	
	if err := uc.validateFile(file); err != nil {
		return fmt.Errorf("validación fallida: %v", err)
	}

	// Construir la ruta de la carpeta (departamento/asunto)
	folderPath := fmt.Sprintf("%s/%s", file.Departamento, file.Asunto) // Usar Asunto

	// Verificar que el archivo no existe ya en Nextcloud
	exists, err := uc.fileStorageService.FileExists(folderPath, file.Nombre)
	if err != nil {
		fmt.Printf("Warning: No se pudo verificar archivo en Nextcloud: %v\n", err)
	} else if exists {
		return fmt.Errorf("el archivo %s ya existe en Nextcloud", file.Nombre)
	}

	// Subir archivo a Nextcloud
	relativePath, err := uc.fileStorageService.UploadFile(folderPath, file.Nombre, fileHeader)
	if err != nil {
		return fmt.Errorf("error al subir archivo a Nextcloud: %v", err)
	}

	// Actualizar el directorio en la entidad con la ruta de Nextcloud
	file.Directorio = relativePath

	// Crear el registro en la base de datos
	if err := uc.repo.Create(file); err != nil {
		
		deleteErr := uc.fileStorageService.DeleteFile(folderPath, file.Nombre)
		if deleteErr != nil {
			fmt.Printf("Error al revertir subida de archivo: %v\n", deleteErr)
		}
		return fmt.Errorf("error al crear archivo en base de datos: %v", err)
	}

	fmt.Printf("Archivo creado exitosamente: %s en %s\n", file.Nombre, relativePath)
	return nil
}

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
	if file.Asunto == "" { // Validar Asunto en lugar de FolderName
		return fmt.Errorf("asunto es requerido")
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