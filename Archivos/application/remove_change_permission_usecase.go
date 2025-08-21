// Archivos/application/remove_change_permission_usecase.go
package application

import (
	entities "VaultDoc-VD/Archivos/domain/entities"
	"VaultDoc-VD/Archivos/domain/repository"
)

type RemoveChangePermissionUseCase struct {
	repo repository.ChangeFileRepository
}

func NewRemoveChangePermissionUseCase(repo repository.ChangeFileRepository) *RemoveChangePermissionUseCase {
	return &RemoveChangePermissionUseCase{repo: repo}
}

func (uc *RemoveChangePermissionUseCase) Execute(changeFile entities.ChangeFile) error {
	return uc.repo.RemovePermission(changeFile)
}