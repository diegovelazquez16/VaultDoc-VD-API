// Archivos/infrastructure/repository/files_postgresql_repository.go
package repository

import (
	"database/sql"
	"fmt"
	
	entities "VaultDoc-VD/Archivos/domain/entities"
	"VaultDoc-VD/core"
)

type FilesPostgreSQLRepository struct {
	db *core.Conn_PostgreSQL
}

func NewFilesPostgreSQLRepository(db *core.Conn_PostgreSQL) *FilesPostgreSQLRepository {
	return &FilesPostgreSQLRepository{db: db}
}

func (r *FilesPostgreSQLRepository) Create(file entities.Files) error {
	query := `INSERT INTO files (departamento, nombre, tamano, fecha, folio, extension, id_folder, id_uploader, directorio) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	
	_, err := r.db.DB.Exec(query, 
		file.Departamento, 
		file.Nombre, 
		file.Tamano, 
		file.Fecha, 
		file.Folio, 
		file.Extension, 
		file.Id_Folder, 
		file.Id_Uploader,
		file.Directorio, // Agregado el campo directorio
	)
	
	if err != nil {
		return fmt.Errorf("error al insertar archivo: %v", err)
	}
	
	return nil
}

func (r *FilesPostgreSQLRepository) GetByID(id int) (entities.Files, error) {
	var file entities.Files
	query := `SELECT id, departamento, nombre, tamano, fecha, folio, extension, id_folder, id_uploader, directorio 
			  FROM files WHERE id = $1`
	
	row := r.db.DB.QueryRow(query, id)
	err := row.Scan(
		&file.Id,
		&file.Departamento,
		&file.Nombre,
		&file.Tamano,
		&file.Fecha,
		&file.Folio,
		&file.Extension,
		&file.Id_Folder,
		&file.Id_Uploader,
		&file.Directorio, // Agregado el campo directorio
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return file, fmt.Errorf("archivo con ID %d no encontrado", id)
		}
		return file, fmt.Errorf("error al obtener archivo: %v", err)
	}
	
	return file, nil
}

func (r *FilesPostgreSQLRepository) Update(file entities.Files) error {
	query := `UPDATE files SET 
			  departamento = $2, 
			  nombre = $3, 
			  tamano = $4, 
			  fecha = $5, 
			  folio = $6, 
			  extension = $7, 
			  id_folder = $8, 
			  id_uploader = $9,
			  directorio = $10 
			  WHERE id = $1`
	
	result, err := r.db.DB.Exec(query,
		file.Id,
		file.Departamento,
		file.Nombre,
		file.Tamano,
		file.Fecha,
		file.Folio,
		file.Extension,
		file.Id_Folder,
		file.Id_Uploader,
		file.Directorio, // Agregado el campo directorio
	)
	
	if err != nil {
		return fmt.Errorf("error al actualizar archivo: %v", err)
	}
	
	// Verificar que se actualizó al menos una fila
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas afectadas: %v", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("archivo con ID %d no encontrado para actualizar", file.Id)
	}
	
	return nil
}

func (r *FilesPostgreSQLRepository) Delete(id int) error {
	query := `DELETE FROM files WHERE id = $1`
	result, err := r.db.DB.Exec(query, id)
	
	if err != nil {
		return fmt.Errorf("error al eliminar archivo: %v", err)
	}
	
	// Verificar que se eliminó al menos una fila
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas afectadas: %v", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("archivo con ID %d no encontrado para eliminar", id)
	}
	
	return nil
}

func (r *FilesPostgreSQLRepository) GetAll() ([]entities.Files, error) {
	var files []entities.Files
	query := `SELECT id, departamento, nombre, tamano, fecha, folio, extension, id_folder, id_uploader, directorio 
			  FROM files ORDER BY id ASC`
	
	rows, err := r.db.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener archivos: %v", err)
	}
	defer rows.Close()
	
	for rows.Next() {
		var file entities.Files
		err := rows.Scan(
			&file.Id,
			&file.Departamento,
			&file.Nombre,
			&file.Tamano,
			&file.Fecha,
			&file.Folio,
			&file.Extension,
			&file.Id_Folder,
			&file.Id_Uploader,
			&file.Directorio, // Agregado el campo directorio
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear archivo: %v", err)
		}
		files = append(files, file)
	}
	
	// Verificar errores durante la iteración
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error durante la iteración: %v", err)
	}
	
	return files, nil
}

// GetByDirectorio obtiene archivos por su directorio
func (r *FilesPostgreSQLRepository) GetByDirectorio(directorio string) ([]entities.Files, error) {
	var files []entities.Files
	query := `SELECT id, departamento, nombre, tamano, fecha, folio, extension, id_folder, id_uploader, directorio 
			  FROM files WHERE directorio LIKE $1 ORDER BY id ASC`
	
	rows, err := r.db.DB.Query(query, "%"+directorio+"%")
	if err != nil {
		return nil, fmt.Errorf("error al obtener archivos por directorio: %v", err)
	}
	defer rows.Close()
	
	for rows.Next() {
		var file entities.Files
		err := rows.Scan(
			&file.Id,
			&file.Departamento,
			&file.Nombre,
			&file.Tamano,
			&file.Fecha,
			&file.Folio,
			&file.Extension,
			&file.Id_Folder,
			&file.Id_Uploader,
			&file.Directorio,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear archivo: %v", err)
		}
		files = append(files, file)
	}
	
	return files, nil
}