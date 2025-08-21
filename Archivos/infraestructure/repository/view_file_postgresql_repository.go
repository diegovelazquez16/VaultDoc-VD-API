// Archivos/infrastructure/repository/view_file_postgresql_repository.go
package repository

import (
	"fmt"
	entities "VaultDoc-VD/Archivos/domain/entities"
	"VaultDoc-VD/core"
)

type ViewFilePostgreSQLRepository struct {
	db *core.Conn_PostgreSQL
}

func NewViewFilePostgreSQLRepository(db *core.Conn_PostgreSQL) *ViewFilePostgreSQLRepository {
	return &ViewFilePostgreSQLRepository{db: db}
}

func (r *ViewFilePostgreSQLRepository) GrantPermission(viewFile entities.ViewFile) error {
	// Verificar que el archivo existe primero
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM files WHERE id = $1)`
	err := r.db.DB.QueryRow(checkQuery, viewFile.Id_File).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error al verificar existencia del archivo: %v", err)
	}
	if !exists {
		return fmt.Errorf("el archivo con ID %d no existe", viewFile.Id_File)
	}

	// Verificar que el usuario existe
	checkUserQuery := `SELECT EXISTS(SELECT 1 FROM usuarios WHERE id = $1)`
	err = r.db.DB.QueryRow(checkUserQuery, viewFile.Id_User).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error al verificar existencia del usuario: %v", err)
	}
	if !exists {
		return fmt.Errorf("el usuario con ID %d no existe", viewFile.Id_User)
	}

	// Insertar el permiso
	query := `INSERT INTO view_files (id_file, id_user) VALUES ($1, $2)
			  ON CONFLICT (id_file, id_user) DO NOTHING`
	
	result, err := r.db.DB.Exec(query, viewFile.Id_File, viewFile.Id_User)
	if err != nil {
		return fmt.Errorf("error al otorgar permiso de visualización: %v", err)
	}

	// Verificar si se insertó algo
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("el permiso de visualización ya existe para este usuario y archivo")
	}

	return nil
}

func (r *ViewFilePostgreSQLRepository) RemovePermission(viewFile entities.ViewFile) error {
	query := `DELETE FROM view_files WHERE id_file = $1 AND id_user = $2`
	
	result, err := r.db.DB.Exec(query, viewFile.Id_File, viewFile.Id_User)
	if err != nil {
		return fmt.Errorf("error al revocar permiso de visualización: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas afectadas: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("permiso de visualización no encontrado para el usuario y archivo especificados")
	}

	return nil
}

func (r *ViewFilePostgreSQLRepository) HasPermission(fileId, userId int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM view_files WHERE id_file = $1 AND id_user = $2)`
	err := r.db.DB.QueryRow(query, fileId, userId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error al verificar permiso de visualización: %v", err)
	}
	return exists, nil
}