//Carpetas/infraestructure/repository/folders_repo_postgresql.go
package repository

import (
	"VaultDoc-VD/Carpetas/domain/entities"
	"VaultDoc-VD/core"
	"fmt"
	"time"
)

type FoldersPostgreSQLRepository struct {
	db core.Conn_PostgreSQL
}

func NewFoldersPostgreSQLRepository(db *core.Conn_PostgreSQL) *FoldersPostgreSQLRepository {
	return &FoldersPostgreSQLRepository{db: *db}
}

func (r *FoldersPostgreSQLRepository) CreateFolder(newFolder entities.Folders) error {
	_, err := r.db.ExecutePreparedQuery(
		"INSERT INTO folders (name, departamento, id_uploader) VALUES ($1, $2, $3)",
		newFolder.Name,
		newFolder.Departamento,
		newFolder.Id_uploader,
	)
	if err != nil {
		return fmt.Errorf("Error al agregar carpeta: %w", err)
	}
	return nil
}


func (r *FoldersPostgreSQLRepository) GetFoldersByDepartament(department string) ([]entities.Folders, error) {
	
	rows := r.db.FetchRows("SELECT id, name, departamento, id_uploader FROM folders WHERE departamento = $1", department)
	if rows == nil {
		return nil, fmt.Errorf("Error al ejecutar consulta")
	}
	defer rows.Close()

	var folders []entities.Folders
	for rows.Next() {
		var folder entities.Folders
		err := rows.Scan(
			&folder.Id,
			&folder.Name,
			&folder.Departamento,
			&folder.Id_uploader,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear folder: %w", err)
		}
		folders = append(folders, folder)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error en iteración de filas: %w", err)
	}

	return folders, nil
}

func (r *FoldersPostgreSQLRepository) GetFolderByFullName(name string) ([]entities.Folders, error) {
	
	rows := r.db.FetchRows("SELECT id, name, departamento, id_uploader FROM folders WHERE name ILIKE $1", name)
	if rows == nil {
		return nil, fmt.Errorf("Error al ejecutar consulta")
	}
	defer rows.Close()

	var folders []entities.Folders
	for rows.Next() {
		var folder entities.Folders
		err := rows.Scan(
			&folder.Id,
			&folder.Name,
			&folder.Departamento,
			&folder.Id_uploader,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear folder: %w", err)
		}
		folders = append(folders, folder)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error en iteración de filas: %w", err)
	}

	return folders, nil
}

func (r *FoldersPostgreSQLRepository) GetFolderByName(name string) ([]entities.Folders, error) {
	query := "SELECT id, name, departamento, id_uploader FROM folders WHERE name ILIKE $1"
	searchPattern := "%" + name + "%"
	fmt.Println("Ejecutando query:", query, "con patrón:", searchPattern)
	
	rows := r.db.FetchRows(query, searchPattern)
	if rows == nil {
		return nil, fmt.Errorf("Error al ejecutar consulta")
	}
	defer rows.Close()

	var folders []entities.Folders
	for rows.Next() {
		var folder entities.Folders
		err := rows.Scan(
			&folder.Id,
			&folder.Name,
			&folder.Departamento,
			&folder.Id_uploader,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear folder: %w", err)
		}
		folders = append(folders, folder)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error en iteración de filas: %w", err)
	}

	return folders, nil
}


func (r *FoldersPostgreSQLRepository) GetFoldersByDepartamentComplete(department string) ([]entities.Folders, error) {
	
	rows := r.db.FetchRows("SELECT * FROM folders WHERE departamento = $1", department)
	if rows == nil {
		return nil, fmt.Errorf("Error al ejecutar consulta")
	}
	defer rows.Close()

	var folders []entities.Folders
	for rows.Next() {
		var folder entities.Folders
		var createdAt, updatedAt time.Time 
		
		err := rows.Scan(
			&folder.Id,
			&folder.Name,
			&folder.Departamento,
			&folder.Id_uploader,
			&createdAt,    
			&updatedAt,    
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear folder: %w", err)
		}
		folders = append(folders, folder)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error en iteración de filas: %w", err)
	}

	return folders, nil
}