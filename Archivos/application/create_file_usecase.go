// Archivos/application/create_file_usecase.go
package application

import (
	"fmt"
	"mime/multipart"
	entities "VaultDoc-VD/Archivos/domain/entities"
	"VaultDoc-VD/Archivos/domain/repository"
	"VaultDoc-VD/Archivos/domain/services"
	folderEntities "VaultDoc-VD/Carpetas/domain/entities"
	folderRepository "VaultDoc-VD/Carpetas/domain/repository"
)

type CreateFileUseCase struct {
	repo                    repository.FilesRepository
	fileStorageService      services.FileStorageService
	changeFileRepo          repository.ChangeFileRepository
	viewFileRepo           repository.ViewFileRepository
	userService            services.UserService
	folderRepo             folderRepository.FoldersRepository 
}

func NewCreateFileUseCase(
	repo repository.FilesRepository, 
	fileStorageService services.FileStorageService,
	changeFileRepo repository.ChangeFileRepository,
	viewFileRepo repository.ViewFileRepository,
	userService services.UserService,
	folderRepo folderRepository.FoldersRepository, 
) *CreateFileUseCase {
	return &CreateFileUseCase{
		repo:                    repo,
		fileStorageService:      fileStorageService,
		changeFileRepo:          changeFileRepo,
		viewFileRepo:           viewFileRepo,
		userService:            userService,
		folderRepo:             folderRepo,
	}
}


func (uc *CreateFileUseCase) Execute(file entities.Files, fileHeader *multipart.FileHeader, userDepartment string) error {
	
	if err := uc.validateFile(file); err != nil {
		return fmt.Errorf("validación fallida: %v", err)
	}

	
	// Como no existe GetByID, buscar por departamento y filtrar por ID
	folders, err := uc.folderRepo.GetFoldersByDepartament(userDepartment)
	if err != nil {
		return fmt.Errorf("error al obtener carpetas del departamento: %v", err)
	}

	var targetFolder *folderEntities.Folders
	for _, folder := range folders {
		if folder.Id == file.Id_Folder {
			targetFolder = &folder
			break
		}
	}

	if targetFolder == nil {
		return fmt.Errorf("carpeta con ID %d no encontrada en el departamento %s", file.Id_Folder, userDepartment)
	}

	// Construir la ruta completa: departamento/nombre_carpeta
	folderPath := userDepartment + "/" + targetFolder.Name

	// Verificar que el archivo no existe ya en Nextcloud
	exists, err := uc.fileStorageService.FileExists(folderPath, file.Nombre)
	if err != nil {
		fmt.Printf("Warning: No se pudo verificar archivo en Nextcloud: %v\n", err)
	} else if exists {
		return fmt.Errorf("el archivo %s ya existe en Nextcloud en el directorio %s", file.Nombre, folderPath)
	}

	
	relativePath, err := uc.fileStorageService.UploadFile(folderPath, file.Nombre, fileHeader)
	if err != nil {
		return fmt.Errorf("error al subir archivo a Nextcloud en directorio %s: %v", folderPath, err)
	}

	
	file.Directorio = relativePath

	// Crear el registro en la base de datos
	if err := uc.repo.Create(file); err != nil {
		
		deleteErr := uc.fileStorageService.DeleteFile(folderPath, file.Nombre)
		if deleteErr != nil {
			fmt.Printf("Error al revertir subida de archivo de directorio %s: %v\n", folderPath, deleteErr)
		}
		return fmt.Errorf("error al crear archivo en base de datos: %v", err)
	}

	
	createdFile, err := uc.repo.GetByFolio(file.Folio)
	if err != nil {
		fmt.Printf("Warning: No se pudo obtener ID del archivo creado para otorgar permisos: %v\n", err)
		return nil // El archivo se creó, pero no se pudieron otorgar permisos automáticos
	}

	// Otorgar permisos automáticamente usando el departamento del JWT
	if err := uc.grantAutomaticPermissions(createdFile.Id, file.Id_Uploader, userDepartment); err != nil {
		fmt.Printf("Warning: Error al otorgar permisos automáticos: %v\n", err)
		// No retornamos error porque el archivo se creó correctamente
	}

	fmt.Printf("Archivo creado exitosamente: %s en directorio %s (ruta completa: %s)\n", file.Nombre, folderPath, relativePath)
	return nil
}

func (uc *CreateFileUseCase) grantAutomaticPermissions(fileId, uploaderId int, userDepartment string) error {
	// Lista de usuarios a los que otorgar permisos
	usersToGrant := []int{uploaderId} // El usuario que subió el archivo

	// Obtener jefes del mismo departamento (id_rol = 2 y departamento específico)
	departmentBosses, err := uc.userService.GetUsersByRoleAndDepartment(2, userDepartment)
	if err != nil {
		fmt.Printf("Warning: No se pudieron obtener jefes del departamento %s: %v\n", userDepartment, err)
	} else {
		usersToGrant = append(usersToGrant, departmentBosses...)
	}

	// Obtener administradores principales (id_rol = 3) - estos no tienen departamento específico
	mainAdmins, err := uc.userService.GetUsersByRole(3)
	if err != nil {
		fmt.Printf("Warning: No se pudieron obtener administradores principales: %v\n", err)
	} else {
		usersToGrant = append(usersToGrant, mainAdmins...)
	}

	// Remover duplicados usando un mapa
	uniqueUsers := make(map[int]bool)
	for _, userId := range usersToGrant {
		uniqueUsers[userId] = true
	}

	// Otorgar permisos de visualización y edición a todos los usuarios únicos
	for userId := range uniqueUsers {
		// Otorgar permiso de visualización
		viewFile := entities.ViewFile{
			Id_File: fileId,
			Id_User: userId,
		}
		if err := uc.viewFileRepo.GrantPermission(viewFile); err != nil {
			fmt.Printf("Warning: No se pudo otorgar permiso de visualización al usuario %d para archivo %d: %v\n", userId, fileId, err)
		}

		// Otorgar permiso de edición
		changeFile := entities.ChangeFile{
			Id_File: fileId,
			Id_User: userId,
		}
		if err := uc.changeFileRepo.GrantPermission(changeFile); err != nil {
			fmt.Printf("Warning: No se pudo otorgar permiso de edición al usuario %d para archivo %d: %v\n", userId, fileId, err)
		}
	}

	fmt.Printf("Permisos automáticos otorgados a %d usuarios del departamento '%s' y admins para archivo %d\n", len(uniqueUsers), userDepartment, fileId)
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
