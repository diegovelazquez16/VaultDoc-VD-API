//Carpetas/domain/repository/FoldersRepo.go
package repository

import "VaultDoc-VD/Carpetas/domain/entities"

type FoldersRepository interface {
	CreateFolder(newFolder entities.Folders) error
	GetFoldersByDepartament(department string) ([]entities.Folders, error)
	GetFoldersByDepartamentComplete(department string) ([]entities.Folders, error)
	GetFolderByFullName(name string) ([]entities.Folders, error)
	GetFolderByName(name string) ([]entities.Folders, error)
	GetFoldersByMyDepartament(departament string) ([]entities.Folders, error)
}