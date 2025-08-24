package application

import (
	"VaultDoc-VD/Carpetas/domain/entities"
	"VaultDoc-VD/Carpetas/domain/repository"
	"fmt"
)

type CreateFolderUseCase struct {
	repo repository.FoldersRepository
}

func NewCreateFolderUseCase(repo repository.FoldersRepository)*CreateFolderUseCase{
	return&CreateFolderUseCase{repo: repo}
}

func(uc *CreateFolderUseCase)Execute(name, departament string, id_uploader int)(*entities.Folders, error){
	equalFolders, err := uc.repo.GetFolderByFullName(name)
	if len(equalFolders) > 0 {
		return nil, fmt.Errorf("el folder %s ya está registrado", name)
	}

	if err != nil {
		return nil, fmt.Errorf("Error al buscar folder: %s", err)
	}

	newFolder := &entities.Folders{Id: 0, Name: name, Departamento: departament, Id_uploader: id_uploader}
	err = uc.repo.CreateFolder(*newFolder)

	if err != nil {
		return nil, fmt.Errorf("Error al registrar archivo: %s", err)
	}

	folder, err := uc.repo.GetFolderByFullName(name)
	if err != nil {
		return nil, err
	}
	return &folder[0], nil
}