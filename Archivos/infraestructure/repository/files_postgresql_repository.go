// Archivos/infrastructure/repository/files_postgresql_repository.go
package repository

import (
	"database/sql"
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
	query := `INSERT INTO files (departamento, nombre, tamano, fecha, folio, extension, id_folder, id_uploader) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	
	_, err := r.db.DB.Exec(query, 
		file.Departamento, 
		file.Nombre, 
		file.Tamano, 
		file.Fecha, 
		file.Folio, 
		file.Extension, 
		file.Id_Folder, 
		file.Id_Uploader,
	)
	
	return err
}

func (r *FilesPostgreSQLRepository) GetByID(id int) (entities.Files, error) {
	var file entities.Files
	query := `SELECT id, departamento, nombre, tamano, fecha, folio, extension, id_folder, id_uploader 
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
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return file, err // Archivo no encontrado
		}
		return file, err
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
			  id_uploader = $9 
			  WHERE id = $1`
	
	_, err := r.db.DB.Exec(query,
		file.Id,
		file.Departamento,
		file.Nombre,
		file.Tamano,
		file.Fecha,
		file.Folio,
		file.Extension,
		file.Id_Folder,
		file.Id_Uploader,
	)
	
	return err
}

func (r *FilesPostgreSQLRepository) Delete(id int) error {
	query := `DELETE FROM files WHERE id = $1`
	_, err := r.db.DB.Exec(query, id)
	return err
}

func (r *FilesPostgreSQLRepository) GetAll() ([]entities.Files, error) {
	var files []entities.Files
	query := `SELECT id, departamento, nombre, tamano, fecha, folio, extension, id_folder, id_uploader 
			  FROM files ORDER BY id ASC`
	
	rows, err := r.db.DB.Query(query)
	if err != nil {
		return nil, err
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
		)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	
	return files, nil
}
