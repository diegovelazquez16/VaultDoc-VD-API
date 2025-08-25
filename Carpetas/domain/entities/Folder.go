package entities
import (
	"time"
)

type Folders struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Departamento string    `json:"departamento"`
	Id_uploader int       `json:"id_uploader"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}