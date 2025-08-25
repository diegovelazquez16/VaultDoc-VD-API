// Carpetas/infrastructure/services/adapters/nextcloud_adapter.go
package adapters

import (
	"VaultDoc-VD/Carpetas/domain/services"
	"VaultDoc-VD/core"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type NextcloudAdapter struct {
	client *core.NextcloudClient
}


var _ services.CloudStorageService = (*NextcloudAdapter)(nil)

func NewNextcloudAdapter() *NextcloudAdapter {
	return &NextcloudAdapter{
		client: core.NewNextcloudClient(),
	}
}

func (nc *NextcloudAdapter) CreateFolder(folderPath string) error {
	// Crear recursivamente las carpetas necesarias
	if err := nc.ensureFolders(folderPath); err != nil {
		return fmt.Errorf("failed to ensure folder %s: %w", folderPath, err)
	}
	return nil
}

// ensureFolders crea las carpetas necesarias recursivamente
func (nc *NextcloudAdapter) ensureFolders(folderPath string) error {
	// Limpiar y normalizar la ruta
	cleanPath := strings.Trim(folderPath, "/")
	if cleanPath == "" || cleanPath == "." {
		return nil
	}

	// Crear carpetas padre recursivamente
	parts := strings.Split(cleanPath, "/")
	currentPath := ""

	for _, part := range parts {
		if part == "" {
			continue
		}

		if currentPath == "" {
			currentPath = part
		} else {
			currentPath = currentPath + "/" + part
		}

		if err := nc.createSingleFolder(currentPath); err != nil {
			return err
		}
	}

	return nil
}

// createSingleFolder crea una sola carpeta
func (nc *NextcloudAdapter) createSingleFolder(folderPath string) error {
	folderURL := fmt.Sprintf("%s/remote.php/dav/files/%s/%s/",
		nc.client.BaseURL, nc.client.Username, folderPath)

	req, err := http.NewRequest("MKCOL", folderURL, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(nc.client.Username, nc.client.Password)
	req.Header.Set("User-Agent", "Go-http-client/1.1")

	resp, err := nc.client.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusMethodNotAllowed {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("MKCOL failed for %s: %d, %s", folderPath, resp.StatusCode, string(body))
	}

	return nil
}

func (nc *NextcloudAdapter) FolderExists(folderPath string) (bool, error) {
	// Normalizar la ruta
	cleanPath := strings.Trim(folderPath, "/")
	url := fmt.Sprintf("%s/remote.php/dav/files/%s/%s/", nc.client.BaseURL, nc.client.Username, cleanPath)

	req, err := http.NewRequest("PROPFIND", url, bytes.NewReader([]byte(`<?xml version="1.0"?>
<d:propfind xmlns:d="DAV:">
	<d:prop>
		<d:resourcetype/>
	</d:prop>
</d:propfind>`)))
	if err != nil {
		return false, fmt.Errorf("error al crear request: %w", err)
	}

	req.SetBasicAuth(nc.client.Username, nc.client.Password)
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Depth", "0")

	resp, err := nc.client.Client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error al ejecutar request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode != http.StatusMultiStatus {
		return false, fmt.Errorf("error al verificar carpeta en Nextcloud (status: %d)", resp.StatusCode)
	}

	return true, nil
}

func (nc *NextcloudAdapter) ListFolders(folderPath string) ([]string, error) {
	// Normalizar la ruta
	cleanPath := strings.Trim(folderPath, "/")
	var url string
	if cleanPath == "" {
		url = fmt.Sprintf("%s/remote.php/dav/files/%s/", nc.client.BaseURL, nc.client.Username)
	} else {
		url = fmt.Sprintf("%s/remote.php/dav/files/%s/%s/", nc.client.BaseURL, nc.client.Username, cleanPath)
	}

	req, err := http.NewRequest("PROPFIND", url, bytes.NewReader([]byte(`<?xml version="1.0"?>
<d:propfind xmlns:d="DAV:">
	<d:prop>
		<d:resourcetype/>
		<d:displayname/>
	</d:prop>
</d:propfind>`)))
	if err != nil {
		return nil, fmt.Errorf("error al crear request: %w", err)
	}

	req.SetBasicAuth(nc.client.Username, nc.client.Password)
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Depth", "1") 

	resp, err := nc.client.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return []string{}, nil
	}

	if resp.StatusCode != http.StatusMultiStatus {
		return nil, fmt.Errorf("error al listar carpetas en Nextcloud (status: %d)", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta: %w", err)
	}


	folders := nc.parseWebDAVResponse(string(body), cleanPath)
	return folders, nil
}

// parseWebDAVResponse extrae los nombres de carpetas del XML de respuesta WebDAV

func (nc *NextcloudAdapter) parseWebDAVResponse(xmlResponse, basePath string) []string {
	var folders []string
	
	// Buscar elementos que contengan <d:collection/> (indica que es una carpeta)
	lines := strings.Split(xmlResponse, "\n")
	var currentHref string
	isCollection := false
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Extraer href
		if strings.Contains(line, "<d:href>") {
			start := strings.Index(line, "<d:href>") + 8
			end := strings.Index(line, "</d:href>")
			if start < end {
				currentHref = line[start:end]
			}
		}
		
		// Verificar si es una colección (carpeta)
		if strings.Contains(line, "<d:collection/>") {
			isCollection = true
		}
		
		// Al final del response de un recurso
		if strings.Contains(line, "</d:response>") && isCollection && currentHref != "" {
			// Extraer solo el nombre de la carpeta del href
			folderName := nc.extractFolderName(currentHref, basePath)
			if folderName != "" && folderName != basePath {
				folders = append(folders, folderName)
			}
			// Reset para el siguiente item
			currentHref = ""
			isCollection = false
		}
	}
	
	return folders
}


func (nc *NextcloudAdapter) extractFolderName(href, basePath string) string {
	
	href = strings.TrimSuffix(href, "/")
	
	
	userPath := fmt.Sprintf("/remote.php/dav/files/%s/", nc.client.Username)
	if strings.HasPrefix(href, userPath) {
		href = strings.TrimPrefix(href, userPath)
	}
	
	
	if basePath != "" {
		basePrefix := basePath + "/"
		if strings.HasPrefix(href, basePrefix) {
			href = strings.TrimPrefix(href, basePrefix)
		}
	}
	
	
	if !strings.Contains(href, "/") {
		return href
	}
	
	return ""
}