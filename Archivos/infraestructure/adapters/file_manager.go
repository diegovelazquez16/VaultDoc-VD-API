// Archivos/infrastructure/adapters/file_manager.go
package adapters

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"github.com/gin-gonic/gin"
)

type FileManager struct{}

func NewFileManager() *FileManager {
	return &FileManager{}
}

// UploadFile sube un archivo a la ruta especificada
func (f *FileManager) UploadFile(file *multipart.FileHeader, foldersDir string, name string, ctx *gin.Context) (string, error) {
	baseDir := os.Getenv("FILES_DIR")
	if baseDir == "" {
		return "", fmt.Errorf("FILES_DIR no está configurado en las variables de entorno")
	}

	// Crear directorio completo
	newDir := filepath.Join(baseDir, foldersDir)
	
	// Crear directorios si no existen
	if err := f.CreateDir("", foldersDir); err != nil {
		return "", fmt.Errorf("error al crear directorios: %v", err)
	}
	
	// Ruta final del archivo
	dst := filepath.Join(newDir, name)
	
	// Guardar archivo
	if err := ctx.SaveUploadedFile(file, dst); err != nil {
		fmt.Printf("Error al almacenar archivo: %v\n", err)
		return dst, err
	}
	
	fmt.Printf("Archivo almacenado en la ubicación: %s\n", dst)
	return dst, nil
}

// DownloadFile permite descargar un archivo
func (f *FileManager) DownloadFile(ctx *gin.Context, dir, filename string) {
	fullPath := filepath.Join(os.Getenv("FILES_DIR"), dir, filename)
	fmt.Printf("Descargando archivo desde: %s\n", fullPath)
	
	// Verificar que el archivo existe
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		ctx.JSON(404, gin.H{"error": "Archivo no encontrado"})
		return
	}
	
	ctx.FileAttachment(fullPath, filename)
}

// DeleteFile elimina un archivo del sistema
func (f *FileManager) DeleteFile(dir_file string) error {
	fullPath := filepath.Join(os.Getenv("FILES_DIR"), dir_file)
	fmt.Printf("Eliminando archivo: %s\n", fullPath)
	
	err := os.Remove(fullPath)
	if err != nil {
		fmt.Printf("Error al eliminar archivo: %v\n", err)
	}
	return err
}

// CreateDir crea directorios de forma recursiva
func (f *FileManager) CreateDir(dir_folder, folder_name string) error {
	baseDir := os.Getenv("FILES_DIR")
	if baseDir == "" {
		return fmt.Errorf("FILES_DIR no está configurado en las variables de entorno")
	}

	var path string
	if dir_folder == "" {
		// Si dir_folder está vacío, usar folder_name directamente
		path = filepath.Join(baseDir, folder_name)
	} else {
		// Combinar ambos
		path = filepath.Join(baseDir, dir_folder, folder_name)
	}
	
	fmt.Printf("Creando directorio: %s\n", path)
	
	// Verificar si ya existe
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			fmt.Printf("Error al crear directorio: %v\n", err)
			return err
		}
		fmt.Printf("Directorio creado exitosamente: %s\n", path)
	} else {
		fmt.Printf("El directorio ya existe: %s\n", path)
	}
	
	return nil
}

// GetFileInfo obtiene información de un archivo
func (f *FileManager) GetFileInfo(filePath string) (os.FileInfo, error) {
	fullPath := filepath.Join(os.Getenv("FILES_DIR"), filePath)
	return os.Stat(fullPath)
}

// FileExists verifica si un archivo existe
func (f *FileManager) FileExists(filePath string) bool {
	fullPath := filepath.Join(os.Getenv("FILES_DIR"), filePath)
	_, err := os.Stat(fullPath)
	return !os.IsNotExist(err)
}