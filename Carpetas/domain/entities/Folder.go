package entities

type Folders struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Departamento string   `json:"departamento"`
	Id_uploader int       `json:"id_uploader"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}