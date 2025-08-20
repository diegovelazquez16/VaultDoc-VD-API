package repository

import "VaultDoc-VD/Archivos/domain/entities"

type ViewFileRepository interface {
	GrantPermission(viewFile ent) error
}