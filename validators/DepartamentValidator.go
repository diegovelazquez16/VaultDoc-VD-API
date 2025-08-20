package validators

import "strings"

func IsValidDepartamento(departamento string) bool {
	validDepartamentos := []string{"Finanzaz", "Operativo", "General"}
	for _, validDept := range validDepartamentos {
		if strings.EqualFold(departamento, validDept) {
			return true
		}
	}
	return false
}