// Archivos/application/remove_view_permission_usecase.go
package application

import (
	entities "VaultDoc-VD/Archivos/domain/entities"
	"VaultDoc-VD/Archivos/domain/repository"
)

type RemoveViewPermissionUseCase struct {
	repo repository.ViewFileRepository
}

func NewRemoveViewPermissionUseCase(repo repository.ViewFileRepository) *RemoveViewPermissionUseCase {
	return &RemoveViewPermissionUseCase{repo: repo}
}

func (uc *RemoveViewPermissionUseCase) Execute(viewFile entities.ViewFile) error {
	return uc.repo.RemovePermission(viewFile)
}