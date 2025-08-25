//core/nextcloud.go
package core

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type NextcloudClient struct {
	BaseURL  string
	Username string
	Password string
	Client   *http.Client
}

func NewNextcloudClient() *NextcloudClient {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error al cargar el archivo .env: %v\n", err)
	}

	return &NextcloudClient{
		BaseURL:  os.Getenv("NEXTCLOUD_BASE_URL"),
		Username: os.Getenv("NEXTCLOUD_USERNAME"),
		Password: os.Getenv("NEXTCLOUD_PASSWORD"),
		Client:   &http.Client{},
	}
}

func (nc *NextcloudClient) CreateFolder(folderPath string) error {
	// Crear recursivamente las carpetas necesarias
	if err := nc.ensureFolders(folderPath); err != nil {
		return fmt.Errorf("failed to ensure folder %s: %w", folderPath, err)
	}
	return nil
}

// ensureFolders crea las carpetas necesarias recursivamente
func (nc *NextcloudClient) ensureFolders(folderPath string) error {
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
func (nc *NextcloudClient) createSingleFolder(folderPath string) error {
	
	folderURL := fmt.Sprintf("%s/remote.php/dav/files/%s/%s/",
		nc.BaseURL, nc.Username, folderPath)

	req, err := http.NewRequest("MKCOL", folderURL, nil)
	if err != nil {
		return err
	}

	
	req.SetBasicAuth(nc.Username, nc.Password)
	
	req.Header.Set("User-Agent", "Go-http-client/1.1")

	resp, err := nc.Client.Do(req)
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

func (nc *NextcloudClient) FolderExists(folderPath string) (bool, error) {
	
	cleanPath := strings.Trim(folderPath, "/")
	url := fmt.Sprintf("%s/remote.php/dav/files/%s/%s/", nc.BaseURL, nc.Username, cleanPath)
	
	req, err := http.NewRequest("PROPFIND", url, bytes.NewReader([]byte(`<?xml version="1.0"?>
		<d:propfind xmlns:d="DAV:">
			<d:prop>
				<d:resourcetype/>
			</d:prop>
		</d:propfind>`)))
	if err != nil {
		return false, fmt.Errorf("error al crear request: %w", err)
	}
	
	req.SetBasicAuth(nc.Username, nc.Password)
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Depth", "0")
	
	resp, err := nc.Client.Do(req)
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

func (nc *NextcloudClient) ListFolders(basePath string) ([]string, error) {
	url := fmt.Sprintf("%s/remote.php/dav/files/%s/%s", nc.BaseURL, nc.Username, basePath)
	
	propfindBody := `<?xml version="1.0"?>
		<d:propfind xmlns:d="DAV:">
			<d:prop>
				<d:resourcetype/>
				<d:displayname/>
			</d:prop>
		</d:propfind>`
	
	req, err := http.NewRequest("PROPFIND", url, bytes.NewReader([]byte(propfindBody)))
	if err != nil {
		return nil, fmt.Errorf("error al crear request: %w", err)
	}
	
	req.SetBasicAuth(nc.Username, nc.Password)
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Depth", "1")
	
	resp, err := nc.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusMultiStatus {
		return nil, fmt.Errorf("error al listar carpetas en Nextcloud (status: %d)", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta: %w", err)
	}
	
	
	folders := nc.parseFolderNames(string(body))
	return folders, nil
}

func (nc *NextcloudClient) parseFolderNames(xmlContent string) []string {
	var folders []string
	
	// Parse básico para extraer displayname de carpetas
	lines := strings.Split(xmlContent, "\n")
	isFolder := false
	
	for i, line := range lines {
		if strings.Contains(line, "<d:collection/>") {
			isFolder = true
		} else if strings.Contains(line, "</d:resourcetype>") && isFolder {
			// Buscar displayname en las líneas siguientes
			for j := i + 1; j < len(lines) && j < i+5; j++ {
				if strings.Contains(lines[j], "<d:displayname>") {
					start := strings.Index(lines[j], "<d:displayname>") + len("<d:displayname>")
					end := strings.Index(lines[j], "</d:displayname>")
					if start < end {
						folderName := lines[j][start:end]
						if folderName != "" && folderName != "." {
							folders = append(folders, folderName)
						}
					}
					break
				}
			}
			isFolder = false
		}
	}
	
	return folders
}