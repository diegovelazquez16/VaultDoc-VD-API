// Archivos/domain/entities/files_entities.go
package domain

type Files struct {
	Id           int    `json:"id"`
	Departamento string `json:"departamento"`
	Nombre       string `json:"nombre"`
	Tamano       int    `json:"tamano"`
	Fecha        string `json:"fecha"`
	Folio        string `json:"folio"`
	Extension    string `json:"extension"`
	Id_Folder    int    `json:"id_folder"`
	Id_Uploader  int    `json:"id_uploader"`
	Directorio   string `json:"directorio"`
}
