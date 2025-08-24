// Archivos/application/grant_change_permission_usecase.go
package application

import (
	entities "VaultDoc-VD/Archivos/domain/entities"
	"VaultDoc-VD/Archivos/domain/repository"
)

type GrantChangePermissionUseCase struct {
	repo repository.ChangeFileRepository
}

func NewGrantChangePermissionUseCase(repo repository.ChangeFileRepository) *GrantChangePermissionUseCase {
	return &GrantChangePermissionUseCase{repo: repo}
}

func (uc *GrantChangePermissionUseCase) Execute(changeFile entities.ChangeFile) error {
	return uc.repo.GrantPermission(changeFile)
}