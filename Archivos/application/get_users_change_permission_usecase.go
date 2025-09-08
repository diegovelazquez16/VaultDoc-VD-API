package application

import (
	"VaultDoc-VD/Archivos/domain/repository"
	"VaultDoc-VD/Usuarios/domain/entities"
)

type GetUsersChangePermissionsUseCase struct {
	repo repository.ChangeFileRepository
}

func NewGetUsersChangePermissionsUseCase(repo repository.ChangeFileRepository)*GetUsersChangePermissionsUseCase{
	return&GetUsersChangePermissionsUseCase{repo: repo}
}

func(uc *GetUsersChangePermissionsUseCase)Execute(fileId int)([]entities.User, []entities.User, error, error) {
	users_whitp, err1 := uc.repo.GetUsersWithChangePermission(fileId)
	users_whitoutp, err2 := uc.repo.GetUsersWithoutChangePermission(fileId)
	return users_whitp, users_whitoutp, err1, err2
}