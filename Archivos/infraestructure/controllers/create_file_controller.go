// Archivos/infrastructure/controllers/create_file_controller.go
package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"time"
	
	"VaultDoc-VD/Archivos/application"
	history "VaultDoc-VD/Historial/application"
	entities "VaultDoc-VD/Archivos/domain/entities"
	historyEntities "VaultDoc-VD/Historial/domain/entities"
	"github.com/gin-gonic/gin"
)

type CreateFileController struct {
	useCase *application.CreateFileUseCase
	historyUseCase *history.SaveActionUseCase
	ucGetbyName *application.GetFileByNameUseCase
}

func NewCreateFileController(
	useCase *application.CreateFileUseCase,
	uc *history.SaveActionUseCase,
	ucGetByName *application.GetFileByNameUseCase,
	) *CreateFileController {
	return &CreateFileController{useCase: useCase, historyUseCase: uc, ucGetbyName: ucGetByName}
}

func (c *CreateFileController) Execute(ctx *gin.Context) {
	// 1. Obtener el archivo del form-data
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error al obtener archivo de la solicitud",
			"error":   err.Error(),
		})
		return
	}

	// 2. Obtener los datos JSON del form-data
	jsonData := ctx.PostForm("json")
	if jsonData == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Campo 'json' es requerido",
			"error":   "No se encontró el campo json en form-data",
		})
		return
	}

	// 3. Parsear el JSON
	var input struct {
		Nombre       string `json:"nombre"`       
		Tamano       int    `json:"tamano"`                         
		Fecha        string `json:"fecha"`                         
		Folio        string `json:"folio" binding:"required"`
		Extension    string `json:"extension"`                      
		Id_Folder    int    `json:"id_folder" binding:"required"`
		Id_Uploader  int    `json:"id_uploader" binding:"required"`
	}

	if err := json.Unmarshal([]byte(jsonData), &input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error al parsear datos JSON",
			"error":   err.Error(),
		})
		return
	}

	// 4. Validar campos requeridos
	if input.Folio == "" || input.Id_Folder == 0 || input.Id_Uploader == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Campos requeridos faltantes",
			"error":   "folio, id_folder, id_uploader son requeridos",
		})
		return
	}

	// 5. Obtener el departamento del usuario desde el middleware JWT (ADAPTACIÓN)
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
		// Los JWT claims pueden venir como interface{}, intentar conversión
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

	// 6. Validar que el departamento esté en los valores permitidos del ENUM
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

	// 7. Obtener información automática del archivo
	if input.Fecha == "" {
		input.Fecha = time.Now().Format("2006-01-02 15:04:05")
	}
	
	// Obtener extensión del archivo
	fileExtension := filepath.Ext(file.Filename)
	if input.Extension == "" {
		input.Extension = fileExtension
	}

	// Obtener tamaño del archivo
	input.Tamano = int(file.Size)

	// 8. Generar nombre final del archivo usando la función del dominio
	
	input.Nombre = application.GenerateFilename(input.Folio, userDepartment)
	finalFileName := input.Nombre
	if filepath.Ext(finalFileName) == "" && fileExtension != "" {
		finalFileName += fileExtension
	}

	// 9. Crear entidad para la base de datos
	
	fileEntity := entities.Files{
		Departamento: userDepartment, // Usar departamento del JWT
		Nombre:       finalFileName,
		Tamano:       input.Tamano,
		Fecha:        input.Fecha,
		Folio:        input.Folio,
		Extension:    input.Extension,
		Id_Folder:    input.Id_Folder,
		Id_Uploader:  input.Id_Uploader,
	}

	// 10. ADAPTACIÓN: Pasar userDepartment al caso de uso para determinar el directorio
	if err := c.useCase.Execute(fileEntity, file, userDepartment); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al crear archivo",
			"error":   err.Error(),
		})
		return
	}

	newFile, err := c.ucGetbyName.Execute(input.Nombre + input.Extension)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Archivo no encontrado",
			"error":   err.Error(),
		})
		return
	}

	var record historyEntities.ReceiveHistory
	record.Departamento = userDepartment
	record.Id_user = input.Id_Uploader
	record.Id_folder = input.Id_Folder
	record.Id_file = newFile.Id
	record.Movimiento = "Subió archivo"

	err = c.historyUseCase.Execute(record)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error interno al crear registro en el historial",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":     "Archivo creado exitosamente en Nextcloud",
		"filename":    finalFileName,
		"size":        input.Tamano,
		"department":  userDepartment, 
		"folio":       input.Folio,
	})
}
