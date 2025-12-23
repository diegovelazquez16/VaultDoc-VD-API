// VaultDoc-VD-API/Historial/infrastructure/repository/history_postgresql_repository.go
package repository

import (
	"VaultDoc-VD/Historial/domain/entities"
	"VaultDoc-VD/core"
	"database/sql"
	"fmt"
)

type HistoryPostgreSQLRepository struct {
	db *core.Conn_PostgreSQL
}

func NewHistoryPostgreSQLRepository(db *core.Conn_PostgreSQL) *HistoryPostgreSQLRepository {
	return &HistoryPostgreSQLRepository{db: db}
}

func (r *HistoryPostgreSQLRepository) SaveAction(history entities.ReceiveHistory) error {
	query := `INSERT INTO history 
		(movimiento, departamento, id_folder, id_file, id_user) 
		VALUES ($1, $2, $3, $4, $5)`
	
	var folderID, fileID, userID interface{}
	if history.Id_folder == 0 {
		folderID = nil
	} else {
		folderID = history.Id_folder
	}
	if history.Id_file == 0 {
		fileID = nil
	} else {
		fileID = history.Id_file
	}
	if history.Id_user == 0 {
		userID = nil
	} else {
		userID = history.Id_user
	}
	
	_, err := r.db.ExecutePreparedQuery(
		query,
		history.Movimiento,
		history.Departamento,
		folderID,
		fileID,
		userID,
	)
	if err != nil {
		return fmt.Errorf("error al insertar registro en el historial: %v", err)
	}
	return nil
}

func (r *HistoryPostgreSQLRepository) GetHistory(departament string) ([]entities.SendHistory, error) {
	var history []entities.SendHistory
	
	query := `SELECT 
		history.id, 
		history.movimiento, 
		history.departamento, 
		history.fecha_registro, 
		usuarios.nombre, 
		usuarios.apellidos, 
		COALESCE(folders.name, history.folder_name, '') as folder_name,
		CASE 
			WHEN history.id_file IS NULL THEN COALESCE(history.file_name, '')
			ELSE COALESCE(files.nombre, history.file_name, '')
		END as file_nombre
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
		var folderName, fileName sql.NullString
		
		err := rows.Scan(
			&record.Id,
			&record.Movimiento,
			&record.Departamento,
			&record.Fecha_registro,
			&record.Id_user.Nombre,
			&record.Id_user.Apellidos,
			&folderName,
			&fileName,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear registro: %v", err)
		}
		
		record.Id_folder.Name = folderName.String
		record.Id_file.Nombre = fileName.String
		
		history = append(history, record)
	}

	return history, nil
}
func (r *HistoryPostgreSQLRepository) GetAllHistory() ([]entities.SendHistory, error) {
	var history []entities.SendHistory
	query := `SELECT 
		history.id, 
		history.movimiento, 
		history.departamento, 
		history.fecha_registro,
		history.folder_name,
		history.file_name,
		history.user_name,
		COALESCE(usuarios.nombre, '') as user_nombre,
		COALESCE(usuarios.apellidos, '') as user_apellidos,
		COALESCE(folders.name, '') as folder_name_actual,
		COALESCE(files.nombre, '') as file_nombre_actual
		FROM history 
		LEFT JOIN usuarios ON history.id_user = usuarios.id 
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
		var folderNameBackup, fileNameBackup, userNameBackup sql.NullString
		var userNombre, userApellidos, folderNameActual, fileNombreActual sql.NullString
		
		err := rows.Scan(
			&record.Id,
			&record.Movimiento,
			&record.Departamento,
			&record.Fecha_registro,
			&folderNameBackup,
			&fileNameBackup,
			&userNameBackup,
			&userNombre,
			&userApellidos,
			&folderNameActual,
			&fileNombreActual,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear registro: %v", err)
		}
		
		if userNombre.Valid {
			record.Id_user.Nombre = userNombre.String
			record.Id_user.Apellidos = userApellidos.String
		} else if userNameBackup.Valid {
			record.UserName = userNameBackup.String
		}
		
		if folderNameActual.Valid {
			record.Id_folder.Name = folderNameActual.String
		} else if folderNameBackup.Valid {
			record.FolderName = folderNameBackup.String
		}
		
		if fileNombreActual.Valid {
			record.Id_file.Nombre = fileNombreActual.String
		} else if fileNameBackup.Valid {
			record.FileName = fileNameBackup.String
		}
		
		history = append(history, record)
	}

	return history, nil
}

func (r *HistoryPostgreSQLRepository) GetHistoryByUser(userID int) ([]entities.SendHistory, error) {
	var history []entities.SendHistory
	query := `SELECT 
		history.id, 
		history.movimiento, 
		history.departamento, 
		history.fecha_registro,
		history.folder_name,
		history.file_name,
		history.user_name,
		COALESCE(usuarios.nombre, '') as user_nombre,
		COALESCE(usuarios.apellidos, '') as user_apellidos,
		COALESCE(folders.name, '') as folder_name_actual,
		COALESCE(files.nombre, '') as file_nombre_actual
		FROM history 
		LEFT JOIN usuarios ON history.id_user = usuarios.id 
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
		var folderNameBackup, fileNameBackup, userNameBackup sql.NullString
		var userNombre, userApellidos, folderNameActual, fileNombreActual sql.NullString
		
		err := rows.Scan(
			&record.Id,
			&record.Movimiento,
			&record.Departamento,
			&record.Fecha_registro,
			&folderNameBackup,
			&fileNameBackup,
			&userNameBackup,
			&userNombre,
			&userApellidos,
			&folderNameActual,
			&fileNombreActual,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear registro: %v", err)
		}
		
		if userNombre.Valid {
			record.Id_user.Nombre = userNombre.String
			record.Id_user.Apellidos = userApellidos.String
		} else if userNameBackup.Valid {
			record.UserName = userNameBackup.String
		}
		
		if folderNameActual.Valid {
			record.Id_folder.Name = folderNameActual.String
		} else if folderNameBackup.Valid {
			record.FolderName = folderNameBackup.String
		}
		
		if fileNombreActual.Valid {
			record.Id_file.Nombre = fileNombreActual.String
		} else if fileNameBackup.Valid {
			record.FileName = fileNameBackup.String
		}
		
		history = append(history, record)
	}

	return history, nil
}

func (r *HistoryPostgreSQLRepository) GetHistoryByFolder(folderID int) ([]entities.SendHistory, error) {
	var history []entities.SendHistory
	query := `SELECT 
		history.id, 
		history.movimiento, 
		history.departamento, 
		history.fecha_registro,
		history.folder_name,
		history.file_name,
		history.user_name,
		COALESCE(usuarios.nombre, '') as user_nombre,
		COALESCE(usuarios.apellidos, '') as user_apellidos,
		COALESCE(folders.name, '') as folder_name_actual,
		COALESCE(files.nombre, '') as file_nombre_actual
		FROM history 
		LEFT JOIN usuarios ON history.id_user = usuarios.id 
		LEFT JOIN folders ON history.id_folder = folders.id 
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
		var folderNameBackup, fileNameBackup, userNameBackup sql.NullString
		var userNombre, userApellidos, folderNameActual, fileNombreActual sql.NullString
		
		err := rows.Scan(
			&record.Id,
			&record.Movimiento,
			&record.Departamento,
			&record.Fecha_registro,
			&folderNameBackup,
			&fileNameBackup,
			&userNameBackup,
			&userNombre,
			&userApellidos,
			&folderNameActual,
			&fileNombreActual,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear registro: %v", err)
		}
		
		if userNombre.Valid {
			record.Id_user.Nombre = userNombre.String
			record.Id_user.Apellidos = userApellidos.String
		} else if userNameBackup.Valid {
			record.UserName = userNameBackup.String
		}
		
		if folderNameActual.Valid {
			record.Id_folder.Name = folderNameActual.String
		} else if folderNameBackup.Valid {
			record.FolderName = folderNameBackup.String
		}
		
		if fileNombreActual.Valid {
			record.Id_file.Nombre = fileNombreActual.String
		} else if fileNameBackup.Valid {
			record.FileName = fileNameBackup.String
		}
		
		history = append(history, record)
	}

	return history, nil
}

func (r *HistoryPostgreSQLRepository) GetHistoryByFile(fileID int) ([]entities.SendHistory, error) {
	var history []entities.SendHistory
	query := `SELECT 
		history.id, 
		history.movimiento, 
		history.departamento, 
		history.fecha_registro,
		history.folder_name,
		history.file_name,
		history.user_name,
		COALESCE(usuarios.nombre, '') as user_nombre,
		COALESCE(usuarios.apellidos, '') as user_apellidos,
		COALESCE(folders.name, '') as folder_name_actual,
		COALESCE(files.nombre, '') as file_nombre_actual
		FROM history 
		LEFT JOIN usuarios ON history.id_user = usuarios.id 
		LEFT JOIN folders ON history.id_folder = folders.id 
		LEFT JOIN files ON history.id_file = files.id 
		WHERE history.id_file = $1
		ORDER BY history.fecha_registro DESC`

	rows, err := r.db.DB.Query(query, fileID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener historial por archivo: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var record entities.SendHistory
		var folderNameBackup, fileNameBackup, userNameBackup sql.NullString
		var userNombre, userApellidos, folderNameActual, fileNombreActual sql.NullString
		
		err := rows.Scan(
			&record.Id,
			&record.Movimiento,
			&record.Departamento,
			&record.Fecha_registro,
			&folderNameBackup,
			&fileNameBackup,
			&userNameBackup,
			&userNombre,
			&userApellidos,
			&folderNameActual,
			&fileNombreActual,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear registro: %v", err)
		}
		
		if userNombre.Valid {
			record.Id_user.Nombre = userNombre.String
			record.Id_user.Apellidos = userApellidos.String
		} else if userNameBackup.Valid {
			record.UserName = userNameBackup.String
		}
		
		if folderNameActual.Valid {
			record.Id_folder.Name = folderNameActual.String
		} else if folderNameBackup.Valid {
			record.FolderName = folderNameBackup.String
		}
		
		if fileNombreActual.Valid {
			record.Id_file.Nombre = fileNombreActual.String
		} else if fileNameBackup.Valid {
			record.FileName = fileNameBackup.String
		}
		
		history = append(history, record)
	}

	return history, nil
}

func (r *HistoryPostgreSQLRepository) GetHistoryByID(id int) (*entities.ReceiveHistory, error) {
	query := `SELECT id, movimiento, departamento, 
		COALESCE(id_folder, 0), COALESCE(id_file, 0), COALESCE(id_user, 0),
		COALESCE(folder_name, ''), COALESCE(file_name, ''), COALESCE(user_name, ''),
		fecha_registro 
		FROM history WHERE id = $1`
	
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
		&record.FolderName,
		&record.FileName,
		&record.UserName,
		&record.Fecha_registro,
	)
	if err != nil {
		return nil, fmt.Errorf("error al escanear registro: %w", err)
	}

	return &record, nil
}

func (r *HistoryPostgreSQLRepository) UpdateHistory(id int, history entities.ReceiveHistory) error {
	// Convertir 0 a NULL
	var folderID, fileID, userID interface{}
	if history.Id_folder == 0 {
		folderID = nil
	} else {
		folderID = history.Id_folder
	}
	if history.Id_file == 0 {
		fileID = nil
	} else {
		fileID = history.Id_file
	}
	if history.Id_user == 0 {
		userID = nil
	} else {
		userID = history.Id_user
	}
	
	_, err := r.db.ExecutePreparedQuery(
		`UPDATE history SET movimiento = $1, departamento = $2, 
		id_folder = $3, id_file = $4, id_user = $5 WHERE id = $6`,
		history.Movimiento,
		history.Departamento,
		folderID,
		fileID,
		userID,
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