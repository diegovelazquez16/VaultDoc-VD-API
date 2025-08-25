// Archivos/application/update_file_usecase.go (Corrección)
package application

import (
	"fmt"
	"mime/multipart"
	entities "VaultDoc-VD/Archivos/domain/entities"
	"VaultDoc-VD/Archivos/domain/repository"
	"VaultDoc-VD/Archivos/domain/services"
)

type UpdateFileUseCase struct {
	repo               repository.FilesRepository
	fileStorageService services.FileStorageService
}

func NewUpdateFileUseCase(repo repository.FilesRepository, fileStorageService services.FileStorageService) *UpdateFileUseCase {
	return &UpdateFileUseCase{
		repo:               repo,
		fileStorageService: fileStorageService,
	}
}

func (uc *UpdateFileUseCase) Execute(file entities.Files, fileHeader *multipart.FileHeader) error {
	// Si se proporciona un nuevo archivo, actualizarlo en Nextcloud
	if fileHeader != nil {
		// Construir la ruta de la carpeta
		folderPath := fmt.Sprintf("%s/%s", file.Departamento, file.Asunto) 

		// Subir el nuevo archivo a Nextcloud (esto sobrescribirá el existente)
		relativePath, err := uc.fileStorageService.UploadFile(folderPath, file.Nombre, fileHeader)
		if err != nil {
			return fmt.Errorf("error al actualizar archivo en Nextcloud: %v", err)
		}

		
		file.Directorio = relativePath
		fmt.Printf("Archivo actualizado en Nextcloud: %s\n", relativePath)
	}

	
	if err := uc.repo.Update(file); err != nil {
		return fmt.Errorf("error al actualizar archivo en base de datos: %v", err)
	}

	return nil
}