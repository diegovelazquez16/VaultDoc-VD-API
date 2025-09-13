package application

import "VaultDoc-VD/Archivos/domain/repository"

type GetChangePermissionsOfAFolderUseCase struct {
	repo repository.ChangeFileRepository
}

func NewGetChangePermissionsOfAFolderUseCase(repo repository.ChangeFileRepository)*GetChangePermissionsOfAFolderUseCase{
	return&GetChangePermissionsOfAFolderUseCase{repo: repo}
}

func(uc *GetChangePermissionsOfAFolderUseCase)Execute(idFolder, idUser int) ([]int, error) {
	return uc.repo.GetChangePermissionsOfAFolder(idFolder, idUser)
}