//Carpetas/application/CreateFolderUC.go
package application

import (
	"VaultDoc-VD/Carpetas/domain/entities"
	"VaultDoc-VD/Carpetas/domain/repository"
	"VaultDoc-VD/core"
	"fmt"
	"strings"
)

type CreateFolderUseCase struct {
	repo           repository.FoldersRepository
	nextcloudClient *core.NextcloudClient
}

func NewCreateFolderUseCase(repo repository.FoldersRepository) *CreateFolderUseCase {
	return &CreateFolderUseCase{
		repo:           repo,
		nextcloudClient: core.NewNextcloudClient(),
	}
}

func (uc CreateFolderUseCase) Execute(name, departament string, id_uploader int) (*entities.Folders, error) {
	// Verificar si la carpeta ya existe en la base de datos
	equalFolders, err := uc.repo.GetFolderByFullName(name)
	if err != nil {
		return nil, fmt.Errorf("Error al buscar folder: %s", err)
	}
	if len(equalFolders) > 0 {
		return nil, fmt.Errorf("el folder %s ya está registrado", name)
	}

	// Crear la ruta completa para Nextcloud
	// Normalizar nombres para evitar problemas con caracteres especiales
	cleanDepartment := strings.TrimSpace(departament)
	cleanName := strings.TrimSpace(name)
	folderPath := fmt.Sprintf("%s/%s", cleanDepartment, cleanName)
	
	// Verificar si la carpeta completa ya existe en Nextcloud
	exists, err := uc.nextcloudClient.FolderExists(folderPath)
	if err != nil {
		// Log del error pero continuamos (Nextcloud podría estar temporalmente no disponible)
		fmt.Printf("Warning: No se pudo verificar carpeta en Nextcloud: %v\n", err)
	} else if exists {
		return nil, fmt.Errorf("la carpeta %s ya existe en Nextcloud", folderPath)
	}

	// Crear la carpeta en Nextcloud (esto creará recursivamente departamento/nombre)
	err = uc.nextcloudClient.CreateFolder(folderPath)
	if err != nil {
		return nil, fmt.Errorf("Error al crear carpeta en Nextcloud: %s", err)
	}

	// Si la creación en Nextcloud fue exitosa, guardar en base de datos
	newFolder := &entities.Folders{
		Id:          0,
		Name:        cleanName,
		Departamento: cleanDepartment,
		Id_uploader: id_uploader,
	}

	err = uc.repo.CreateFolder(*newFolder)
	if err != nil {
		// Si falla la BD, podríamos considerar eliminar la carpeta de Nextcloud pero por simplicidad, solo devolvemos el error
		
		return nil, fmt.Errorf("Error al registrar folder en base de datos: %s", err)
	}

	
	folder, err := uc.repo.GetFolderByFullName(cleanName)
	if err != nil {
		return nil, err
	}

	return &folder[0], nil
}