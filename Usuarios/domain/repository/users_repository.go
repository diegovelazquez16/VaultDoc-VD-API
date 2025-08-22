// Usuarios/domain/repository/users_repository.go
package repository

import "VaultDoc-VD/Usuarios/domain/entities"

type UserRepository interface {
	Save(user entities.User) error
	FindById(id int) (*entities.User, error)
	FindAll() ([]entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	Update(user entities.User) error
	UpdateProfile(user entities.User) error
	Delete(id int) error
}
