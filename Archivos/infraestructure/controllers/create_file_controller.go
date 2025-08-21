// Archivos/infrastructure/controllers/create_file_controller.go
package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	_"strconv"
	"time"
	
	"VaultDoc-VD/Archivos/application"
	entities "VaultDoc-VD/Archivos/domain/entities"
	"github.com/gin-gonic/gin"
)

type CreateFileController struct {
	useCase *application.CreateFileUseCase
}

func NewCreateFileController(useCase *application.CreateFileUseCase) *CreateFileController {
	return &CreateFileController{useCase: useCase}
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
		Departamento string `json:"departamento" binding:"required"`
		Asunto       string `json:"asunto" binding:"required"`       // Nueva carpeta dentro del departamento
		Nombre       string `json:"nombre" binding:"required"`       // SOTCH-DVD-004-2025
		Tamano       int    `json:"tamano"`                         // Se calculará automáticamente
		Fecha        string `json:"fecha"`                          // Se asignará automáticamente si no viene
		Folio        string `json:"folio" binding:"required"`
		Extension    string `json:"extension"`                      // Se obtendrá del archivo
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

	// 4. Validar campos requeridos y ENUM departamento
	if input.Departamento == "" || input.Asunto == "" || input.Nombre == "" || 
	   input.Folio == "" || input.Id_Folder == 0 || input.Id_Uploader == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Campos requeridos faltantes",
			"error":   "departamento, asunto, nombre, folio, id_folder, id_uploader son requeridos",
		})
		return
	}

	// Validar que el departamento esté en los valores permitidos del ENUM
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
		if input.Departamento == dept {
			isValidDepartment = true
			break
		}
	}
	
	if !isValidDepartment {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Departamento no válido",
			"error":   fmt.Sprintf("El departamento '%s' no está permitido. Valores válidos: %v", input.Departamento, validDepartments),
		})
		return
	}

	// 5. Obtener información automática del archivo
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

	// 6. Crear estructura de carpetas y guardar archivo
	baseDir := os.Getenv("FILES_DIR")
	if baseDir == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error de configuración",
			"error":   "FILES_DIR no está configurado",
		})
		return
	}

	// Crear ruta: FILES_DIR/departamento/asunto/
	folderPath := filepath.Join(baseDir, input.Departamento, input.Asunto)
	if err := c.createDirectories(folderPath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al crear directorios",
			"error":   err.Error(),
		})
		return
	}

	input.Nombre = application.GenerateFilename(input.Folio, input.Departamento)
	// 7. Generar nombre final del archivo (puede incluir extensión si no la tiene)
	finalFileName := input.Nombre
	if filepath.Ext(finalFileName) == "" && fileExtension != "" {
		finalFileName += fileExtension
	}

	// Ruta completa del archivo
	filePath := filepath.Join(folderPath, finalFileName)
	
	// 8. Guardar archivo físico
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al almacenar archivo",
			"error":   err.Error(),
		})
		return
	}

	// 9. Crear directorio relativo para la base de datos
	// Formato: departamento/asunto/archivo.ext
	relativePath := filepath.Join(input.Departamento, input.Asunto, finalFileName)

	// 10. Crear entidad para la base de datos
	fileEntity := entities.Files{
		Departamento: input.Departamento,
		Nombre:       finalFileName,
		Tamano:       input.Tamano,
		Fecha:        input.Fecha,
		Folio:        input.Folio,
		Extension:    input.Extension,
		Id_Folder:    input.Id_Folder,
		Id_Uploader:  input.Id_Uploader,
		Directorio:   relativePath, // Ruta relativa desde FILES_DIR
	}

	// 11. Guardar en base de datos
	if err := c.useCase.Execute(fileEntity); err != nil {
		// Si falla la BD, eliminar archivo físico
		os.Remove(filePath)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al crear registro en base de datos",
			"error":   err.Error(),
		})
		return
	}

	// 12. Respuesta exitosa
	ctx.JSON(http.StatusCreated, gin.H{
		"message":     "Archivo creado exitosamente",
		"filename":    finalFileName,
		"size":        input.Tamano,
		"path":        relativePath,
		"full_path":   filePath,
		"department":  input.Departamento,
		"subject":     input.Asunto,
	})
}

// createDirectories crea los directorios necesarios si no existen
func (c *CreateFileController) createDirectories(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("error al crear directorios: %v", err)
		}
		fmt.Printf("Directorios creados: %s\n", path)
	}
	return nil
}