// Usuarios/infraestructure/repository/users_repo_postgresql.go
package repository

import (
	"VaultDoc-VD/Usuarios/domain/entities"
	"VaultDoc-VD/Usuarios/domain/repository"
	"VaultDoc-VD/core"
	"fmt"
)

type UserPostgreSQLRepository struct {
	db *core.Conn_PostgreSQL
}

func NewUserPostgreSQLRepository(db *core.Conn_PostgreSQL) repository.UserRepository {
	return &UserPostgreSQLRepository{db: db}
}

func (r *UserPostgreSQLRepository) Save(user entities.User) error {
	query := `
		INSERT INTO usuarios (nombre, apellidos, email, password, id_rol, departamento) 
		VALUES ($1, $2, $3, $4, $5, $6)`
	
	_, err := r.db.ExecutePreparedQuery(
		query,
		user.Nombre,
		user.Apellidos,
		user.Email,
		user.Password,
		user.Id_Rol,
		user.Departamento,
	)
	
	if err != nil {
		return fmt.Errorf("error al guardar usuario: %w", err)
	}
	
	return nil
}

func (r *UserPostgreSQLRepository) FindById(id int) (*entities.User, error) {
	query := `
		SELECT id, nombre, apellidos, email, password, id_rol, departamento 
		FROM usuarios 
		WHERE id = $1`
	
	rows := r.db.FetchRows(query, id)
	if rows == nil {
		return nil, fmt.Errorf("error al ejecutar consulta")
	}
	defer rows.Close()
	
	if !rows.Next() {
		return nil, fmt.Errorf("usuario con ID %d no encontrado", id)
	}
	
	var user entities.User
	err := rows.Scan(
		&user.Id,
		&user.Nombre,
		&user.Apellidos,
		&user.Email,
		&user.Password,
		&user.Id_Rol,
		&user.Departamento,
	)
	
	if err != nil {
		return nil, fmt.Errorf("error al escanear usuario: %w", err)
	}
	
	return &user, nil
}

func (r *UserPostgreSQLRepository) FindAll() ([]entities.User, error) {
	query := `
		SELECT id, nombre, apellidos, email, password, id_rol, departamento 
		FROM usuarios 
		ORDER BY id`
	
	rows := r.db.FetchRows(query)
	if rows == nil {
		return nil, fmt.Errorf("error al ejecutar consulta")
	}
	defer rows.Close()
	
	var users []entities.User
	for rows.Next() {
		var user entities.User
		err := rows.Scan(
			&user.Id,
			&user.Nombre,
			&user.Apellidos,
			&user.Email,
			&user.Password,
			&user.Id_Rol,
			&user.Departamento,
		)
		
		if err != nil {
			return nil, fmt.Errorf("error al escanear usuario: %w", err)
		}
		
		// Limpiar password para seguridad
		user.Password = ""
		users = append(users, user)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error en iteración de filas: %w", err)
	}
	
	return users, nil
}

func (r *UserPostgreSQLRepository) FindByEmail(email string) (*entities.User, error) {
	query := `
		SELECT id, nombre, apellidos, email, password, id_rol, departamento 
		FROM usuarios 
		WHERE email = $1`
	
	rows := r.db.FetchRows(query, email)
	if rows == nil {
		return nil, fmt.Errorf("error al ejecutar consulta")
	}
	defer rows.Close()
	
	if !rows.Next() {
		return nil, fmt.Errorf("usuario con email %s no encontrado", email)
	}
	
	var user entities.User
	err := rows.Scan(
		&user.Id,
		&user.Nombre,
		&user.Apellidos,
		&user.Email,
		&user.Password,
		&user.Id_Rol,
		&user.Departamento,
	)
	
	if err != nil {
		return nil, fmt.Errorf("error al escanear usuario: %w", err)
	}
	
	return &user, nil
}

func (r *UserPostgreSQLRepository) Update(user entities.User) error {
	query := `
		UPDATE usuarios 
		SET nombre = $1, apellidos = $2, email = $3, password = $4, id_rol = $5, departamento = $6 
		WHERE id = $7`
	
	result, err := r.db.ExecutePreparedQuery(
		query,
		user.Nombre,
		user.Apellidos,
		user.Email,
		user.Password,
		user.Id_Rol,
		user.Departamento,
		user.Id,
	)
	
	if err != nil {
		return fmt.Errorf("error al actualizar usuario: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas afectadas: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("usuario con ID %d no encontrado para actualizar", user.Id)
	}
	
	return nil
}

func (r *UserPostgreSQLRepository) Delete(id int) error {
	query := `DELETE FROM usuarios WHERE id = $1`
	
	result, err := r.db.ExecutePreparedQuery(query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar usuario: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas afectadas: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("usuario con ID %d no encontrado para eliminar", id)
	}
	
	return nil
}

func (r *UserPostgreSQLRepository) UpdateProfile(user entities.User) error {
	query := `
		UPDATE usuarios 
		SET nombre = $1, apellidos = $2, email = $3, password = $4
		WHERE id = $5`
	
	result, err := r.db.ExecutePreparedQuery(
		query,
		user.Nombre,
		user.Apellidos,
		user.Email,
		user.Password,
		user.Id,
	)
	
	if err != nil {
		return fmt.Errorf("error al actualizar perfil de usuario: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas afectadas: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("usuario con ID %d no encontrado para actualizar perfil", user.Id)
	}
	
	return nil
}

func (r *UserPostgreSQLRepository) GetProfile(userID int) (*entities.User, error) {
	query := `
		SELECT id, nombre, apellidos, email, id_rol, departamento 
		FROM usuarios 
		WHERE id = $1`
	
	rows := r.db.FetchRows(query, userID)
	if rows == nil {
		return nil, fmt.Errorf("error al ejecutar consulta")
	}
	defer rows.Close()
	
	if !rows.Next() {
		return nil, fmt.Errorf("perfil de usuario con ID %d no encontrado", userID)
	}
	
	var user entities.User
	err := rows.Scan(
		&user.Id,
		&user.Nombre,
		&user.Apellidos,
		&user.Email,
		&user.Id_Rol,
		&user.Departamento,
	)
	
	if err != nil {
		return nil, fmt.Errorf("error al escanear datos del perfil: %w", err)
	}
	
	// No incluir password en la consulta del perfil por seguridad
	user.Password = ""
	
	return &user, nil
}

func (r *UserPostgreSQLRepository) FindByDepartment(departamento string) ([]entities.User, error) {
	query := `
		SELECT id, nombre, apellidos, email, password, id_rol, departamento 
		FROM usuarios 
		WHERE departamento = $1
		ORDER BY id`
	
	rows := r.db.FetchRows(query, departamento)
	if rows == nil {
		return nil, fmt.Errorf("error al ejecutar consulta")
	}
	defer rows.Close()
	
	var users []entities.User
	for rows.Next() {
		var user entities.User
		err := rows.Scan(
			&user.Id,
			&user.Nombre,
			&user.Apellidos,
			&user.Email,
			&user.Password,
			&user.Id_Rol,
			&user.Departamento,
		)
		
		if err != nil {
			return nil, fmt.Errorf("error al escanear usuario: %w", err)
		}
		
		// Limpiar password para seguridad
		user.Password = ""
		users = append(users, user)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error en iteración de filas: %w", err)
	}
	
	return users, nil
}