package validators

import "strings"

func IsValidDepartamento(departamento string) bool {
	validDepartamentos := []string{
		"Dirección General", "Área Técnica", "Comisaria", "Coordinación Juridica", 
		"Gerencia Administrativa", "Gerencia Operativa", "Departamento de Finanzas", 
		"Departamento de Planeación", "Departamento de Sistema Eléctrico", 
		"Departamento de Sistema Hidrosánitario y Aire Acondicionado", 
		"Departamento de Mantenimiento General", "Departamento de Voz y Datos", 
		"Departamento de Seguridad e Higiene",
	}
	for _, validDept := range validDepartamentos {
		if strings.EqualFold(departamento, validDept) {
			return true
		}
	}
	return false
}