// Archivos/infrastructure/repository/view_file_postgresql_repository.go
package repository

import (
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
	query := `INSERT INTO view_files (id_file, id_user) VALUES ($1, $2)
			  ON CONFLICT (id_file, id_user) DO NOTHING`
	
	_, err := r.db.DB.Exec(query, viewFile.Id_File, viewFile.Id_User)
	return err
}

func (r *ViewFilePostgreSQLRepository) RemovePermission(viewFile entities.ViewFile) error {
	query := `DELETE FROM view_files WHERE id_file = $1 AND id_user = $2`
	
	_, err := r.db.DB.Exec(query, viewFile.Id_File, viewFile.Id_User)
	return err
}