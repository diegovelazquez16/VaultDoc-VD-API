package application

import (
	"VaultDoc-VD/Archivos/domain/repository"
	"VaultDoc-VD/Usuarios/domain/entities"
)

type GetUsersViewPermissionsUseCase struct {
	repo repository.ViewFileRepository
}

func NewGetUsersViewPermissionsUseCase(repo repository.ViewFileRepository)*GetUsersViewPermissionsUseCase{
	return&GetUsersViewPermissionsUseCase{repo: repo}
}

func(uc *GetUsersViewPermissionsUseCase)Execute(fileId int)([]entities.User, []entities.User, error, error){
	users_withp, err1 := uc.repo.GetUsersWithViewPermission(fileId)
	users_withoutp, err2 := uc.repo.GetUsersWithoutViewPermission(fileId)
	return users_withp, users_withoutp, err1, err2
}