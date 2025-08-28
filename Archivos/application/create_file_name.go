// Archivos/application/create_file_name.go
package application

import (
	"fmt"
	"strings"
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

// Nueva función para generar nombre en actualización manteniendo el año original
func GenerateFilenameForUpdate(folio, departament, originalName string) string {
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

	// Extraer el año del nombre original
	originalYear := extractYearFromFileName(originalName)
	if originalYear == "" {
		// Si no se puede extraer el año, usar el año actual como fallback
		now := time.Now()
		originalYear = fmt.Sprintf("%d", now.Year())
	}

	return "SOTCH-" + abreviation + "-" + folio + "-" + originalYear
}

// Función para extraer el año del nombre de archivo original
func extractYearFromFileName(fileName string) string {
	// Remover la extensión si existe
	nameWithoutExt := strings.TrimSuffix(fileName, getFileExtension(fileName))
	
	// Dividir por guiones y tomar la última parte (debería ser el año)
	parts := strings.Split(nameWithoutExt, "-")
	if len(parts) >= 4 {
		return parts[len(parts)-1] // Último elemento debería ser el año
	}
	return ""
}

// Función para obtener la extensión del archivo
func getFileExtension(fileName string) string {
	lastDot := strings.LastIndex(fileName, ".")
	if lastDot == -1 {
		return ""
	}
	return fileName[lastDot:]
}

// Función para generar nombre completo con extensión
func GenerateFullFileName(baseFileName, extension string) string {
	if extension == "" {
		return baseFileName
	}
	
	// Asegurar que la extensión comience con punto
	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}
	
	return baseFileName + extension
}