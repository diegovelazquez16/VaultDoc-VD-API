//Usuarios/domain/entities/users.go
package entities

type User struct {
	Id           int    `json:"id" db:"id"`
	Nombre       string `json:"nombre" db:"nombre" binding:"required"`
	Apellidos    string `json:"apellidos" db:"apellidos" binding:"required"`
	Email        string `json:"email" db:"email" binding:"required,email"`
	Password     string `json:"password,omitempty" db:"password"`
	Id_Rol       int    `json:"id_rol" db:"id_rol"`
	Departamento string `json:"departamento" db:"departamento"`
}