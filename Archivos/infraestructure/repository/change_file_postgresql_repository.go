// Archivos/infrastructure/repository/change_file_postgresql_repository.go
package repository

import (
	"fmt"
	entities "VaultDoc-VD/Archivos/domain/entities"
	"VaultDoc-VD/core"
)

type ChangeFilePostgreSQLRepository struct {
	db *core.Conn_PostgreSQL
}

func NewChangeFilePostgreSQLRepository(db *core.Conn_PostgreSQL) *ChangeFilePostgreSQLRepository {
	return &ChangeFilePostgreSQLRepository{db: db}
}

func (r *ChangeFilePostgreSQLRepository) GrantPermission(changeFile entities.ChangeFile) error {
	// Verificar que el archivo existe primero
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM files WHERE id = $1)`
	err := r.db.DB.QueryRow(checkQuery, changeFile.Id_File).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error al verificar existencia del archivo: %v", err)
	}
	if !exists {
		return fmt.Errorf("el archivo con ID %d no existe", changeFile.Id_File)
	}

	// Verificar que el usuario existe
	checkUserQuery := `SELECT EXISTS(SELECT 1 FROM usuarios WHERE id = $1)`
	err = r.db.DB.QueryRow(checkUserQuery, changeFile.Id_User).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error al verificar existencia del usuario: %v", err)
	}
	if !exists {
		return fmt.Errorf("el usuario con ID %d no existe", changeFile.Id_User)
	}

	// Insertar el permiso
	query := `INSERT INTO change_files (id_file, id_user) VALUES ($1, $2)
			  ON CONFLICT (id_file, id_user) DO NOTHING`
	
	result, err := r.db.DB.Exec(query, changeFile.Id_File, changeFile.Id_User)
	if err != nil {
		return fmt.Errorf("error al otorgar permiso de edición: %v", err)
	}

	// Verificar si se insertó algo
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("el permiso de edición ya existe para este usuario y archivo")
	}

	return nil
}

func (r *ChangeFilePostgreSQLRepository) RemovePermission(changeFile entities.ChangeFile) error {
	query := `DELETE FROM change_files WHERE id_file = $1 AND id_user = $2`
	
	result, err := r.db.DB.Exec(query, changeFile.Id_File, changeFile.Id_User)
	if err != nil {
		return fmt.Errorf("error al revocar permiso de edición: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas afectadas: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("permiso de edición no encontrado para el usuario y archivo especificados")
	}

	return nil
}

func (r *ChangeFilePostgreSQLRepository) HasPermission(fileId, userId int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM change_files WHERE id_file = $1 AND id_user = $2)`
	err := r.db.DB.QueryRow(query, fileId, userId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error al verificar permiso de edición: %v", err)
	}
	return exists, nil
}