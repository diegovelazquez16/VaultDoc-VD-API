// Usuarios/application/getUsersByDepartment_useCase.go
package application

import (
	"VaultDoc-VD/Usuarios/domain/entities"
	"VaultDoc-VD/Usuarios/domain/repository"
)

type GetUsersByDepartment struct {
	db repository.UserRepository
}

func NewGetUsersByDepartmentUseCase(db repository.UserRepository) *GetUsersByDepartment {
	return &GetUsersByDepartment{db: db}
}

func (gud *GetUsersByDepartment) Execute(departamento string) ([]entities.User, error) {
	users, err := gud.db.FindByDepartment(departamento)
	if err != nil {
		return nil, err
	}
	return users, nil
}