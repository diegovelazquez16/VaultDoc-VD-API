// VaultDoc-VD-API/Historial/infrastructure/repository/history_postgresql_repository.go
package repository

import (
	"VaultDoc-VD/Historial/domain/entities"
	"VaultDoc-VD/core"
	"fmt"
)

type HistoryPostgreSQLRepository struct {
	db *core.Conn_PostgreSQL
}

func NewHistoryPostgreSQLRepository(db *core.Conn_PostgreSQL) *HistoryPostgreSQLRepository {
	return &HistoryPostgreSQLRepository{db: db}
}

func (r *HistoryPostgreSQLRepository) SaveAction(history entities.ReceiveHistory) error {
	_, err := r.db.ExecutePreparedQuery(
		"INSERT INTO history (movimiento, departamento, id_folder, id_file, id_user) VALUES ($1, $2, $3, $4, $5)",
		history.Movimiento,
		history.Departamento,
		history.Id_folder,
		history.Id_file,
		history.Id_user,
	)
	if err != nil {
		return fmt.Errorf("error al insertar registro en el historial: %v", err)
	}
	return nil
}

func (r *HistoryPostgreSQLRepository) GetHistory(departament string) ([]entities.SendHistory, error) {
	var history []entities.SendHistory
	query := `SELECT history.id, history.movimiento, history.departamento, history.fecha_registro, 
		usuarios.nombre, usuarios.apellidos, 
		COALESCE(folders.name, '') as folder_name, 
		COALESCE(files.nombre, '') as file_name 
		FROM history 
		INNER JOIN usuarios ON history.id_user = usuarios.id 
		LEFT JOIN folders ON history.id_folder = folders.id 
		LEFT JOIN files ON history.id_file = files.id 
		WHERE history.departamento = $1
		ORDER BY history.fecha_registro DESC`

	rows, err := r.db.DB.Query(query, departament)
	if err != nil {
		return nil, fmt.Errorf("error al obtener historial: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var record entities.SendHistory
		err := rows.Scan(
			&record.Id,
			&record.Movimiento,
			&record.Departamento,
			&record.Fecha_registro,
			&record.Id_user.Nombre,
			&record.Id_user.Apellidos,
			&record.Id_folder.Name,
			&record.Id_file.Nombre,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear registro: %v", err)
		}
		history = append(history, record)
	}

	return history, nil
}

func (r *HistoryPostgreSQLRepository) GetAllHistory() ([]entities.SendHistory, error) {
	var history []entities.SendHistory
	query := `SELECT history.id, history.movimiento, history.departamento, history.fecha_registro, 
		usuarios.nombre, usuarios.apellidos, 
		COALESCE(folders.name, '') as folder_name, 
		COALESCE(files.nombre, '') as file_name 
		FROM history 
		INNER JOIN usuarios ON history.id_user = usuarios.id 
		LEFT JOIN folders ON history.id_folder = folders.id 
		LEFT JOIN files ON history.id_file = files.id 
		ORDER BY history.fecha_registro DESC`

	rows, err := r.db.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener historial completo: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var record entities.SendHistory
		err := rows.Scan(
			&record.Id,
			&record.Movimiento,
			&record.Departamento,
			&record.Fecha_registro,
			&record.Id_user.Nombre,
			&record.Id_user.Apellidos,
			&record.Id_folder.Name,
			&record.Id_file.Nombre,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear registro: %v", err)
		}
		history = append(history, record)
	}

	return history, nil
}

func (r *HistoryPostgreSQLRepository) GetHistoryByUser(userID int) ([]entities.SendHistory, error) {
	var history []entities.SendHistory
	query := `SELECT history.id, history.movimiento, history.departamento, history.fecha_registro, 
		usuarios.nombre, usuarios.apellidos, 
		COALESCE(folders.name, '') as folder_name, 
		COALESCE(files.nombre, '') as file_name 
		FROM history 
		INNER JOIN usuarios ON history.id_user = usuarios.id 
		LEFT JOIN folders ON history.id_folder = folders.id 
		LEFT JOIN files ON history.id_file = files.id 
		WHERE history.id_user = $1
		ORDER BY history.fecha_registro DESC`

	rows, err := r.db.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener historial por usuario: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var record entities.SendHistory
		err := rows.Scan(
			&record.Id,
			&record.Movimiento,
			&record.Departamento,
			&record.Fecha_registro,
			&record.Id_user.Nombre,
			&record.Id_user.Apellidos,
			&record.Id_folder.Name,
			&record.Id_file.Nombre,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear registro: %v", err)
		}
		history = append(history, record)
	}

	return history, nil
}

func (r *HistoryPostgreSQLRepository) GetHistoryByFolder(folderID int) ([]entities.SendHistory, error) {
	var history []entities.SendHistory
	query := `SELECT history.id, history.movimiento, history.departamento, history.fecha_registro, 
		usuarios.nombre, usuarios.apellidos, 
		folders.name, 
		COALESCE(files.nombre, '') as file_name 
		FROM history 
		INNER JOIN usuarios ON history.id_user = usuarios.id 
		INNER JOIN folders ON history.id_folder = folders.id 
		LEFT JOIN files ON history.id_file = files.id 
		WHERE history.id_folder = $1
		ORDER BY history.fecha_registro DESC`

	rows, err := r.db.DB.Query(query, folderID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener historial por carpeta: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var record entities.SendHistory
		err := rows.Scan(
			&record.Id,
			&record.Movimiento,
			&record.Departamento,
			&record.Fecha_registro,
			&record.Id_user.Nombre,
			&record.Id_user.Apellidos,
			&record.Id_folder.Name,
			&record.Id_file.Nombre,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear registro: %v", err)
		}
		history = append(history, record)
	}

	return history, nil
}

func (r *HistoryPostgreSQLRepository) GetHistoryByFile(fileID int) ([]entities.SendHistory, error) {
	var history []entities.SendHistory
	query := `SELECT history.id, history.movimiento, history.departamento, history.fecha_registro, 
		usuarios.nombre, usuarios.apellidos, 
		COALESCE(folders.name, '') as folder_name, 
		files.nombre 
		FROM history 
		INNER JOIN usuarios ON history.id_user = usuarios.id 
		LEFT JOIN folders ON history.id_folder = folders.id 
		INNER JOIN files ON history.id_file = files.id 
		WHERE history.id_file = $1
		ORDER BY history.fecha_registro DESC`

	rows, err := r.db.DB.Query(query, fileID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener historial por archivo: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var record entities.SendHistory
		err := rows.Scan(
			&record.Id,
			&record.Movimiento,
			&record.Departamento,
			&record.Fecha_registro,
			&record.Id_user.Nombre,
			&record.Id_user.Apellidos,
			&record.Id_folder.Name,
			&record.Id_file.Nombre,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear registro: %v", err)
		}
		history = append(history, record)
	}

	return history, nil
}

func (r *HistoryPostgreSQLRepository) GetHistoryByID(id int) (*entities.ReceiveHistory, error) {
	query := `SELECT * FROM history WHERE id = $1`
	rows := r.db.FetchRows(query, id)
	if rows == nil {
		return nil, fmt.Errorf("error al ejecutar consulta")
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("registro con ID %d no encontrado", id)
	}

	var record entities.ReceiveHistory
	err := rows.Scan(
		&record.Id,
		&record.Movimiento,
		&record.Departamento,
		&record.Id_folder,
		&record.Id_file,
		&record.Id_user,
		&record.Fecha_registro,
	)
	if err != nil {
		return nil, fmt.Errorf("error al escanear registro: %w", err)
	}

	return &record, nil
}

func (r *HistoryPostgreSQLRepository) UpdateHistory(id int, history entities.ReceiveHistory) error {
	_, err := r.db.ExecutePreparedQuery(
		"UPDATE history SET movimiento = $1, departamento = $2, id_folder = $3, id_file = $4, id_user = $5 WHERE id = $6",
		history.Movimiento,
		history.Departamento,
		history.Id_folder,
		history.Id_file,
		history.Id_user,
		id,
	)
	if err != nil {
		return fmt.Errorf("error al actualizar registro en el historial: %v", err)
	}
	return nil
}

func (r *HistoryPostgreSQLRepository) DeleteHistory(id int) error {
	_, err := r.db.ExecutePreparedQuery("DELETE FROM history WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("error al eliminar registro del historial: %v", err)
	}
	return nil
}