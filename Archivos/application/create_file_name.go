// Archivos/application/create_file_name.go
package application

import (
	"fmt"
	"time"
)

func GenerateFilename(folio, departament string) string {
	var abreviation string
	switch departament {
	case "Dirección General":
		abreviation = "DG"
	case "Área Técnica":
		abreviation = "AT"
	case "Comisaria":
		abreviation = "C"
	case "Coordinación Juridica":
		abreviation = "CJ"
	case "Gerencia Administrativa":
		abreviation = "GA"
	case "Gerencia Operativa":
		abreviation = "GO"
	case "Departamento de Finanzas":
		abreviation = "DF"
	case "Departamento de Planeación":
		abreviation = "DP"
	case "Departamento de Sistema Eléctrico":
		abreviation = "DSE"
	case "Departamento de Sistema Hidrosánitario y Aire Acondicionado":
		abreviation = "DSHAA"
	case "Departamento de Mantenimiento General":
		abreviation = "DMG"
	case "Departamento de Voz y Datos":
		abreviation = "DVD"
	case "Departamento de Seguridad e Higiene":
		abreviation = "DSH"
	}
	now := time.Now()
    year := now.Year()
    yearString := fmt.Sprintf("%d", year)
	return "SOTCH-" + abreviation + "-" + folio + "-" + yearString
}