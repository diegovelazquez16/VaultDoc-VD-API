// Archivos/infrastructure/adapters/user_postgresql_adapter.go
package adapters

import (
	"fmt"
	"VaultDoc-VD/core"
)

type UserPostgreSQLAdapter struct {
	db *core.Conn_PostgreSQL
}

func NewUserPostgreSQLAdapter(db *core.Conn_PostgreSQL) *UserPostgreSQLAdapter {
	return &UserPostgreSQLAdapter{db: db}
}

func (u *UserPostgreSQLAdapter) GetUsersByRole(roleId int) ([]int, error) {
	var userIds []int
	query := `SELECT id FROM usuarios WHERE id_rol = $1`
	
	rows, err := u.db.DB.Query(query, roleId)
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuarios por rol: %v", err)
	}
	defer rows.Close()
	
	for rows.Next() {
		var userId int
		err := rows.Scan(&userId)
		if err != nil {
			return nil, fmt.Errorf("error al escanear usuario: %v", err)
		}
		userIds = append(userIds, userId)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error durante la iteración: %v", err)
	}
	
	return userIds, nil
}

func (u *UserPostgreSQLAdapter) GetUsersByRoleAndDepartment(roleId int, department string) ([]int, error) {
	var userIds []int
	query := `SELECT id FROM usuarios WHERE id_rol = $1 AND departamento = $2`
	
	rows, err := u.db.DB.Query(query, roleId, department)
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuarios por rol y departamento: %v", err)
	}
	defer rows.Close()
	
	for rows.Next() {
		var userId int
		err := rows.Scan(&userId)
		if err != nil {
			return nil, fmt.Errorf("error al escanear usuario: %v", err)
		}
		userIds = append(userIds, userId)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error durante la iteración: %v", err)
	}
	
	return userIds, nil
}