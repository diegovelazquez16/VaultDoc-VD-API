// Archivos/domain/services/user_service.go
package services

type UserService interface {
	GetUsersByRole(roleId int) ([]int, error) // Retorna lista de IDs de usuarios
	GetUsersByRoleAndDepartment(roleId int, department string) ([]int, error) // Retorna usuarios por rol y departamento
}