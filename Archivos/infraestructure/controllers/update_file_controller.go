// Archivos/infrastructure/controllers/update_file_controller.go (Actualizado)
package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"VaultDoc-VD/Archivos/application"
	entities "VaultDoc-VD/Archivos/domain/entities"
	"github.com/gin-gonic/gin"
	
)

type UpdateFileController struct {
	useCase *application.UpdateFileUseCase
}

func NewUpdateFileController(useCase *application.UpdateFileUseCase) *UpdateFileController {
	return &UpdateFileController{useCase: useCase}
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

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ID inválido",
			"error":   "El ID debe ser un número entero válido",
		})
		return
	}

	// 2. Verificar si viene un archivo nuevo (opcional para actualización)
	file, err := ctx.FormFile("file")
	hasNewFile := err == nil && file != nil

	// 3. Obtener los datos JSON del form-data
	jsonData := ctx.PostForm("json")
	if jsonData == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Campo 'json' es requerido",
		})
		return
	}

	// 4. Parsear el JSON
	var input struct {
		Departamento string `json:"departamento"`
		Asunto       string `json:"asunto"`
		Nombre       string `json:"nombre"`
		Tamano       int    `json:"tamano"`
		Fecha        string `json:"fecha"`
		Folio        string `json:"folio"`
		Extension    string `json:"extension"`
		Id_Folder    int    `json:"id_folder"`
		Id_Uploader  int    `json:"id_uploader"`
	}

	if err := json.Unmarshal([]byte(jsonData), &input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error al parsear datos JSON",
			"error":   err.Error(),
		})
		return
	}

	
	/* // 5. Si hay archivo nuevo, actualizar información del archivo
	if hasNewFile {
		input.Tamano = int(file.Size)
		if input.Extension == "" {
			input.Extension = filepath.Ext(file.Filename)
		}
	}
	*/
	// 6. Crear entidad para actualizar
	fileEntity := entities.Files{
		Id:           id,
		Departamento: input.Departamento,
		Nombre:       input.Nombre,
		Tamano:       input.Tamano,
		Fecha:        input.Fecha,
		Folio:        input.Folio,
		Extension:    input.Extension,
		Id_Folder:    input.Id_Folder,
		Id_Uploader:  input.Id_Uploader,
		Asunto:       input.Asunto, 
	}


	if hasNewFile {
		err = c.useCase.Execute(fileEntity, file)
	} else {
		err = c.useCase.Execute(fileEntity, nil)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al actualizar archivo",
			"error":   err.Error(),
		})
		return
	}


	response := gin.H{
		"message":    "Archivo actualizado exitosamente",
		"id":         id,
		"department": input.Departamento,
		"subject":    input.Asunto,
	}

	if hasNewFile {
		response["new_file"] = "Archivo físico actualizado en Nextcloud"
		response["size"] = input.Tamano
	}

	ctx.JSON(http.StatusOK, response)
}