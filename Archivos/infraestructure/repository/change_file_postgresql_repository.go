// Archivos/infrastructure/repository/change_file_postgresql_repository.go
package repository

import (
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
	query := `INSERT INTO change_files (id_file, id_user) VALUES ($1, $2)
			  ON CONFLICT (id_file, id_user) DO NOTHING`
	
	_, err := r.db.DB.Exec(query, changeFile.Id_File, changeFile.Id_User)
	return err
}

func (r *ChangeFilePostgreSQLRepository) RemovePermission(changeFile entities.ChangeFile) error {
	query := `DELETE FROM change_files WHERE id_file = $1 AND id_user = $2`
	
	_, err := r.db.DB.Exec(query, changeFile.Id_File, changeFile.Id_User)
	return err
}
