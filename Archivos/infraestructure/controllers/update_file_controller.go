// Archivos/infrastructure/controllers/update_file_controller.go
package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"VaultDoc-VD/Archivos/application"
	history "VaultDoc-VD/Historial/application"
	historyEntities "VaultDoc-VD/Historial/domain/entities"
	"github.com/gin-gonic/gin"
)

type UpdateFileController struct {
	useCase        *application.UpdateFileUseCase
	historyUseCase *history.SaveActionUseCase
	getFileUseCase *application.GetFileByIdUseCase
}

func NewUpdateFileController(
	useCase *application.UpdateFileUseCase, 
	historyUseCase *history.SaveActionUseCase,
	getFileUseCase *application.GetFileByIdUseCase,
) *UpdateFileController {
	return &UpdateFileController{
		useCase:        useCase,
		historyUseCase: historyUseCase,
		getFileUseCase: getFileUseCase,
	}
}

func (c *UpdateFileController) Execute(ctx *gin.Context) {
	// 1. Obtener ID del archivo
	idParam := ctx.Param("id")
	if idParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID del archivo requerido",
		})
		return
	}

	idUserParam := ctx.Param("id_user")
	if idUserParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID del usuario requerido",
		})
		return
	}

	fileId, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID del archivo inválido",
			"error":   "El ID debe ser un número entero válido",
		})
		return
	}

	userId, err := strconv.Atoi(idUserParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID del usuario inválido",
			"error":   "El ID debe ser un número entero válido",
		})
		return
	}

	// 2. Obtener el departamento del usuario desde el middleware JWT
	userDepartmentInterface, exists := ctx.Get("department")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error del sistema",
			"error":   "No se pudo obtener el departamento del usuario del token",
		})
		return
	}
	
	// Convertir correctamente el valor del JWT claim a string
	var userDepartment string
	switch v := userDepartmentInterface.(type) {
	case string:
		userDepartment = v
	case interface{}:
		if str, ok := v.(string); ok {
			userDepartment = str
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error del sistema",
				"error":   "Formato de departamento inválido en el token",
			})
			return
		}
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error del sistema",
			"error":   "Tipo de departamento no reconocido en el token",
		})
		return
	}

	// 3. Validar que el archivo existe y obtener su información actual
	currentFile, err := c.getFileUseCase.Execute(fileId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Archivo no encontrado",
			"error":   err.Error(),
		})
		return
	}

	// 4. Verificar si viene un archivo nuevo (opcional para actualización)
	file, err := ctx.FormFile("file")
	hasNewFile := err == nil && file != nil

	// 5. Obtener los datos JSON del form-data
	jsonData := ctx.PostForm("json")
	if jsonData == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Campo 'json' es requerido",
		})
		return
	}

	// 6. Parsear el JSON
	var input struct {
		Folio       string `json:"folio" binding:"required"`
		Id_Folder   int    `json:"id_folder" binding:"required"`
		Id_Uploader int    `json:"id_uploader" binding:"required"`
	}

	if err := json.Unmarshal([]byte(jsonData), &input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error al parsear datos JSON",
			"error":   err.Error(),
		})
		return
	}

	// 7. Validar campos requeridos
	if input.Folio == "" || input.Id_Folder == 0 || input.Id_Uploader == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Campos requeridos faltantes",
			"error":   "folio, id_folder, id_uploader son requeridos",
		})
		return
	}

	// 8. Validar que el departamento esté en los valores permitidos del ENUM
	validDepartments := []string{
		"Dirección General", 
		"Área Técnica", 
		"Comisaria", 
		"Coordinación Juridica", 
		"Gerencia Administrativa", 
		"Gerencia Operativa", 
		"Departamento de Finanzas", 
		"Departamento de Planeación", 
		"Departamento de Sistema Eléctrico", 
		"Departamento de Sistema Hidrosánitario y Aire Acondicionado", 
		"Departamento de Mantenimiento General", 
		"Departamento de Voz y Datos", 
		"Departamento de Seguridad e Higiene",
	}
	
	isValidDepartment := false
	for _, dept := range validDepartments {
		if userDepartment == dept { 
			isValidDepartment = true
			break
		}
	}
	
	if !isValidDepartment {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Departamento no válido",
			"error":   fmt.Sprintf("El departamento del usuario '%s' no está permitido", userDepartment),
		})
		return
	}

	// 9. Ejecutar la actualización usando el nuevo use case
	if err := c.useCase.Execute(fileId, userId, input.Folio, file, userDepartment); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al actualizar archivo",
			"error":   err.Error(),
		})
		return
	}

	// 10. Obtener el archivo actualizado para mostrar la información
	updatedFile, err := c.getFileUseCase.Execute(fileId)
	if err != nil {
		// No es crítico si no podemos obtener el archivo actualizado
		fmt.Printf("Warning: No se pudo obtener archivo actualizado: %v\n", err)
	}

	// 11. Crear registro en el historial
	var record historyEntities.ReceiveHistory
	record.Departamento = userDepartment
	record.Id_user = userId
	record.Id_folder = input.Id_Folder
	record.Id_file = fileId
	record.Movimiento = "Modificó información de un archivo"

	err = c.historyUseCase.Execute(record)
	if err != nil {
		// Log el error pero no falla la operación
		fmt.Printf("Warning: Error al crear registro en el historial: %v\n", err)
	}

	// 12. Preparar respuesta
	response := gin.H{
		"message":    "Archivo actualizado exitosamente",
		"id":         fileId,
		"department": userDepartment,
		"old_folio":  currentFile.Folio,
		"new_folio":  input.Folio,
	}

	if hasNewFile {
		response["new_file"] = "Archivo físico actualizado en Nextcloud"
		response["size"] = file.Size
		response["original_filename"] = file.Filename
	}

	if updatedFile != (struct {
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
	}{}) {
		response["updated_file"] = map[string]interface{}{
			"id":         updatedFile.Id,
			"name":       updatedFile.Nombre,
			"folio":      updatedFile.Folio,
			"extension":  updatedFile.Extension,
			"size":       updatedFile.Tamano,
			"directory":  updatedFile.Directorio,
		}
	}

	ctx.JSON(http.StatusOK, response)
}