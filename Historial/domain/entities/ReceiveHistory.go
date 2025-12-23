// VaultDoc-VD-API/Historial/domain/entities/ReceiveHistory.go
package entities

type ReceiveHistory struct {
	Id             int    `json:"id"`
	Movimiento     string `json:"movimiento"`
	Departamento   string `json:"departamento"`
	Id_folder      int    `json:"id_folder"`
	Id_file        int    `json:"id_file"`
	Id_user        int    `json:"id_user"`
	
	FolderName     string `json:"folder_name,omitempty"`
	FileName       string `json:"file_name,omitempty"`
	UserName       string `json:"user_name,omitempty"`
	
	Fecha_registro string `json:"fecha_registro"`
}