// Archivos/application/get_files_by_folder_usecase.go
package application

import (
	entities "VaultDoc-VD/Archivos/domain/entities"
	"VaultDoc-VD/Archivos/domain/repository"
)

type GetFilesByFolderUseCase struct {
	repo repository.FilesRepository
}

func NewGetFilesByFolderUseCase(repo repository.FilesRepository) *GetFilesByFolderUseCase {
	return &GetFilesByFolderUseCase{repo: repo}
}

func (uc *GetFilesByFolderUseCase) Execute(folderId int) ([]entities.Files, error) {
	return uc.repo.GetByFolder(folderId)
}