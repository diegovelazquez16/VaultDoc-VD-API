package application

import (
	"VaultDoc-VD/Usuarios/domain/entities"
	"VaultDoc-VD/Usuarios/domain/repository"
)

type GetUsers struct {
	db repository.UserRepository
}

func NewGetUsersUseCase(db repository.UserRepository) *GetUsers {
	return &GetUsers{db: db}
}

func (gu *GetUsers) Execute() ([]entities.User, error) {
	users, err := gu.db.FindAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}
