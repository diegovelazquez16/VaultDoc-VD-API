// Archivos/infraestructure/adapters/nextcloud_file_adapter.go
package adapters

import (
	"VaultDoc-VD/Archivos/domain/services"
	"VaultDoc-VD/core"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

type NextcloudFileAdapter struct {
	client *core.NextcloudClient
}

// Estructura para parsear respuesta XML de PROPFIND
type PropfindResponse struct {
	XMLName   xml.Name `xml:"multistatus"`
	Responses []struct {
		Href     string `xml:"href"`
		Propstat struct {
			Prop struct {
				DisplayName       string `xml:"displayname"`
				ContentLength     string `xml:"getcontentlength"`
				LastModified      string `xml:"getlastmodified"`
				ContentType       string `xml:"getcontenttype"`
			} `xml:"prop"`
			Status string `xml:"status"`
		} `xml:"propstat"`
	} `xml:"response"`
}

// Verificar que implementa la interfaz
var _ services.FileStorageService = (*NextcloudFileAdapter)(nil)

func NewNextcloudFileAdapter() *NextcloudFileAdapter {
	return &NextcloudFileAdapter{
		client: core.NewNextcloudClient(),
	}
}

// Helper para construir URLs de archivo
func (nf *NextcloudFileAdapter) buildFileURL(folderPath, fileName string) string {
	cleanFolderPath := strings.Trim(folderPath, "/")
	if cleanFolderPath == "" {
		return fmt.Sprintf("%s/remote.php/dav/files/%s/%s",
			nf.client.BaseURL, nf.client.Username, fileName)
	}
	return fmt.Sprintf("%s/remote.php/dav/files/%s/%s/%s",
		nf.client.BaseURL, nf.client.Username, cleanFolderPath, fileName)
}

func (nf *NextcloudFileAdapter) UploadFile(folderPath string, fileName string, fileHeader *multipart.FileHeader) (string, error) {
	// Validar parámetros de entrada
	if fileName == "" {
		return "", fmt.Errorf("nombre de archivo no puede estar vacío")
	}
	if fileHeader == nil {
		return "", fmt.Errorf("fileHeader no puede ser nil")
	}

	// Abrir el archivo desde el multipart
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("error al abrir archivo: %w", err)
	}
	defer src.Close()

	// Leer el contenido del archivo
	fileContent, err := io.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("error al leer archivo: %w", err)
	}

	// Construir la URL de destino
	fileURL := nf.buildFileURL(folderPath, fileName)

	// Crear request PUT para subir el archivo
	req, err := http.NewRequest("PUT", fileURL, bytes.NewReader(fileContent))
	if err != nil {
		return "", fmt.Errorf("error al crear request: %w", err)
	}

	req.SetBasicAuth(nf.client.Username, nf.client.Password)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Length", strconv.FormatInt(fileHeader.Size, 10))

	resp, err := nf.client.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error al subir archivo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("error al subir archivo a Nextcloud (status: %d): %s", resp.StatusCode, string(body))
	}

	cleanFolderPath := strings.Trim(folderPath, "/")
	return path.Join(cleanFolderPath, fileName), nil
}

func (nf *NextcloudFileAdapter) DownloadFile(folderPath string, fileName string) ([]byte, error) {
	// Validar parámetros de entrada
	if fileName == "" {
		return nil, fmt.Errorf("nombre de archivo no puede estar vacío")
	}

	// Construir la URL del archivo
	fileURL := nf.buildFileURL(folderPath, fileName)

	req, err := http.NewRequest("GET", fileURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error al crear request: %w", err)
	}

	req.SetBasicAuth(nf.client.Username, nf.client.Password)

	resp, err := nf.client.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al descargar archivo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("archivo no encontrado: %s/%s", folderPath, fileName)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error al descargar archivo (status: %d): %s", resp.StatusCode, string(body))
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer contenido del archivo: %w", err)
	}

	return content, nil
}

func (nf *NextcloudFileAdapter) DeleteFile(folderPath string, fileName string) error {
	// Validar parámetros de entrada
	if fileName == "" {
		return fmt.Errorf("nombre de archivo no puede estar vacío")
	}

	// Construir la URL del archivo
	fileURL := nf.buildFileURL(folderPath, fileName)

	req, err := http.NewRequest("DELETE", fileURL, nil)
	if err != nil {
		return fmt.Errorf("error al crear request: %w", err)
	}

	req.SetBasicAuth(nf.client.Username, nf.client.Password)

	resp, err := nf.client.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error al eliminar archivo: %w", err)
	}
	defer resp.Body.Close()

	// 204 No Content es éxito, 404 Not Found también se puede considerar éxito
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error al eliminar archivo de Nextcloud (status: %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func (nf *NextcloudFileAdapter) FileExists(folderPath string, fileName string) (bool, error) {
	// Validar parámetros de entrada
	if fileName == "" {
		return false, fmt.Errorf("nombre de archivo no puede estar vacío")
	}

	// Construir la URL del archivo
	fileURL := nf.buildFileURL(folderPath, fileName)

	req, err := http.NewRequest("HEAD", fileURL, nil)
	if err != nil {
		return false, fmt.Errorf("error al crear request: %w", err)
	}

	req.SetBasicAuth(nf.client.Username, nf.client.Password)

	resp, err := nf.client.Client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error al verificar archivo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("error al verificar archivo (status: %d)", resp.StatusCode)
	}

	return true, nil
}

func (nf *NextcloudFileAdapter) GetFileInfo(folderPath string, fileName string) (*services.FileInfo, error) {
	// Validar parámetros de entrada
	if fileName == "" {
		return nil, fmt.Errorf("nombre de archivo no puede estar vacío")
	}

	// Construir la URL del archivo
	fileURL := nf.buildFileURL(folderPath, fileName)

	// XML para PROPFIND request
	propfindXML := `<?xml version="1.0"?>
<d:propfind xmlns:d="DAV:">
	<d:prop>
		<d:displayname/>
		<d:getcontentlength/>
		<d:getlastmodified/>
		<d:getcontenttype/>
	</d:prop>
</d:propfind>`

	req, err := http.NewRequest("PROPFIND", fileURL, strings.NewReader(propfindXML))
	if err != nil {
		return nil, fmt.Errorf("error al crear request: %w", err)
	}

	req.SetBasicAuth(nf.client.Username, nf.client.Password)
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Depth", "0")

	resp, err := nf.client.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al obtener información del archivo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("archivo no encontrado: %s/%s", folderPath, fileName)
	}

	if resp.StatusCode != http.StatusMultiStatus {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error al obtener información del archivo (status: %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta: %w", err)
	}

	fileInfo, err := nf.parseFileInfoResponse(body, fileName)
	if err != nil {
		return nil, fmt.Errorf("error al parsear respuesta XML: %w", err)
	}

	return fileInfo, nil
}

func (nf *NextcloudFileAdapter) parseFileInfoResponse(xmlData []byte, fileName string) (*services.FileInfo, error) {
	var propfindResp PropfindResponse
	
	err := xml.Unmarshal(xmlData, &propfindResp)
	if err != nil {
		return nil, fmt.Errorf("error al parsear XML: %w", err)
	}

	if len(propfindResp.Responses) == 0 {
		return nil, fmt.Errorf("no se encontró información del archivo en la respuesta")
	}

	response := propfindResp.Responses[0]
	prop := response.Propstat.Prop

	fileInfo := &services.FileInfo{
		Name:        fileName,
		ContentType: prop.ContentType,
	}

	// Parsear tamaño del archivo
	if prop.ContentLength != "" {
		if size, err := strconv.ParseInt(prop.ContentLength, 10, 64); err == nil {
			fileInfo.Size = size
		}
	}

	// Parsear fecha de modificación
	if prop.LastModified != "" {
		// Nextcloud devuelve fechas en formato RFC1123
		if parsedTime, err := time.Parse(time.RFC1123, prop.LastModified); err == nil {
			fileInfo.LastModified = parsedTime.Format(time.RFC3339)
		} else {
			// Fallback: usar la fecha como string
			fileInfo.LastModified = prop.LastModified
		}
	}

	return fileInfo, nil
}

func (nf *NextcloudFileAdapter) UploadFileFromBytes(folderPath string, fileName string, fileContent []byte) (string, error) {
	// Validar parámetros de entrada
	if fileName == "" {
		return "", fmt.Errorf("nombre de archivo no puede estar vacío")
	}
	if fileContent == nil {
		return "", fmt.Errorf("contenido del archivo no puede ser nil")
	}

	// Construir la URL de destino
	fileURL := nf.buildFileURL(folderPath, fileName)

	// Crear request PUT para subir el archivo
	req, err := http.NewRequest("PUT", fileURL, bytes.NewReader(fileContent))
	if err != nil {
		return "", fmt.Errorf("error al crear request: %w", err)
	}

	req.SetBasicAuth(nf.client.Username, nf.client.Password)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Length", strconv.Itoa(len(fileContent)))

	resp, err := nf.client.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error al subir archivo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("error al subir archivo a Nextcloud (status: %d): %s", resp.StatusCode, string(body))
	}

	cleanFolderPath := strings.Trim(folderPath, "/")
	return path.Join(cleanFolderPath, fileName), nil
}