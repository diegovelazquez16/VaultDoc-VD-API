// Archivos/application/grant_view_permission_usecase.go
package application

import (
	entities "VaultDoc-VD/Archivos/domain/entities"
	"VaultDoc-VD/Archivos/domain/repository"
)

type GrantViewPermissionUseCase struct {
	repo repository.ViewFileRepository
}

func NewGrantViewPermissionUseCase(repo repository.ViewFileRepository) *GrantViewPermissionUseCase {
	return &GrantViewPermissionUseCase{repo: repo}
}

func (uc *GrantViewPermissionUseCase) Execute(viewFile entities.ViewFile) error {
	return uc.repo.GrantPermission(viewFile)
}