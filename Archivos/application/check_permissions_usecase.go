// Archivos/application/check_permissions_usecase.go
package application

import (
	"VaultDoc-VD/Archivos/domain/repository"
)

type CheckPermissionsUseCase struct {
	changeRepo repository.ChangeFileRepository
	viewRepo   repository.ViewFileRepository
	filesRepo  repository.FilesRepository
}

func NewCheckPermissionsUseCase(changeRepo repository.ChangeFileRepository, viewRepo repository.ViewFileRepository, filesRepo repository.FilesRepository) *CheckPermissionsUseCase {
	return &CheckPermissionsUseCase{
		changeRepo: changeRepo,
		viewRepo:   viewRepo,
		filesRepo:  filesRepo,
	}
}

type UserPermissions struct {
	CanView   bool `json:"can_view"`
	CanChange bool `json:"can_change"`
	CanDelete bool `json:"can_delete"`
}

func (uc *CheckPermissionsUseCase) Execute(fileId, userId int) (UserPermissions, error) {
	permissions := UserPermissions{
		CanView:   false,
		CanChange: false,
		CanDelete: false,
	}

	// Verificar que el archivo existe
	_, err := uc.filesRepo.GetByID(fileId)
	if err != nil {
		return permissions, err
	}

	// Verificar permisos de visualización
	canView, err := uc.viewRepo.HasPermission(fileId, userId)
	if err != nil {
		return permissions, err
	}
	permissions.CanView = canView

	// Verificar permisos de edición
	canChange, err := uc.changeRepo.HasPermission(fileId, userId)
	if err != nil {
		return permissions, err
	}
	permissions.CanChange = canChange

	// Los permisos de eliminación están basados en los permisos de edición
	permissions.CanDelete = canChange

	return permissions, nil
}
