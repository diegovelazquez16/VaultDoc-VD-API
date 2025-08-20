// Archivos/domain/entities/files_entities.go
package domain

type Files struct {
	Id           int
	Departamento string
	Nombre       string
	Tamano       int
	Fecha        string
	Folio        string
    Extension    string
	Id_Folder    int
	Id_Uploader  int
}
