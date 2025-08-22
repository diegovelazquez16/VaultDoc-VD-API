package repository

import (
	"VaultDoc-VD/Historial/domain/entities"
	"VaultDoc-VD/core"
	"fmt"
)

type HistoryPostgreSQLRepository struct {
	db *core.Conn_PostgreSQL
}

func NewHistoryPostgreSQLRepository(db *core.Conn_PostgreSQL)*HistoryPostgreSQLRepository{
	return&HistoryPostgreSQLRepository{db: db}
}

func (r *HistoryPostgreSQLRepository) SaveAction(history entities.ReceiveHistory) (error) {
	_, err := r.db.ExecutePreparedQuery(
		"INSERT INTO history (movimiento, departamento, id_folder, id_file, id_user) VALUES	($1, $2, $3, $4, $5)",
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
	query := `SELECT history.id, history.movimiento, history.departamento, history.fecha_registro, usuarios.nombre, 
	usuarios.apellidos, folders.name, files.nombre FROM history INNER JOIN usuarios 
	ON history.id_user = usuarios.id INNER JOIN folders ON history.id_folder = folders.id INNER JOIN files 
	ON history.id_file = files.id WHERE history.departamento = $1`
	
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
			return nil, fmt.Errorf("error al escanear archivo: %v", err)
		}
		history = append(history, record)
	}
	
	return history, nil
}

func (r *HistoryPostgreSQLRepository) GetHistoryByID(id int) (*entities.ReceiveHistory, error) {
	query := `
		SELECT * FROM history WHERE id = $1`
	
	rows := r.db.FetchRows(query, id)
	if rows == nil {
		return nil, fmt.Errorf("error al ejecutar consulta")
	}
	defer rows.Close()
	
	if !rows.Next() {
		return nil, fmt.Errorf("usuario con ID %d no encontrado", id)
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
		return nil, fmt.Errorf("error al escanear usuario: %w", err)
	}
	
	return &record, nil
}