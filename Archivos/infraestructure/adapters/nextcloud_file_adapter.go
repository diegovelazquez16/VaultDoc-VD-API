// Archivos/infrastructure/adapters/nextcloud_file_adapter.go
package adapters

import (
	"VaultDoc-VD/Archivos/domain/services"
	"VaultDoc-VD/core"
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

type NextcloudFileAdapter struct {
	client     *core.NextcloudClient
	workerPool *WorkerPool
}

// Estructura para parsear respuesta XML de PROPFIND
type PropfindResponse struct {
	XMLName   xml.Name `xml:"multistatus"`
	Responses []struct {
		Href     string `xml:"href"`
		Propstat struct {
			Prop struct {
				DisplayName   string `xml:"displayname"`
				ContentLength string `xml:"getcontentlength"`
				LastModified  string `xml:"getlastmodified"`
				ContentType   string `xml:"getcontenttype"`
			} `xml:"prop"`
			Status string `xml:"status"`
		} `xml:"propstat"`
	} `xml:"response"`
}

// UploadJob representa un trabajo de upload
type UploadJob struct {
	FolderPath string
	FileName   string
	Content    []byte
	ResultChan chan UploadResult
}

// UploadResult contiene el resultado de un upload
type UploadResult struct {
	Path  string
	Error error
}

// WorkerPool maneja un pool de workers para uploads concurrentes
type WorkerPool struct {
	jobs       chan UploadJob
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	client     *core.NextcloudClient
	maxRetries int
}

// NewWorkerPool crea un nuevo pool de workers
func NewWorkerPool(client *core.NextcloudClient, maxWorkers int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	
	pool := &WorkerPool{
		jobs:       make(chan UploadJob, maxWorkers*2), // Buffer para evitar bloqueos
		ctx:        ctx,
		cancel:     cancel,
		client:     client,
		maxRetries: client.MaxRetries,
	}

	// Iniciar workers
	for i := 0; i < maxWorkers; i++ {
		pool.wg.Add(1)
		go pool.worker(i)
	}

	return pool
}

// worker procesa los trabajos de upload
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for {
		select {
		case <-wp.ctx.Done():
			return
		case job, ok := <-wp.jobs:
			if !ok {
				return
			}
			
			// Procesar el upload con reintentos
			result := wp.processUpload(job)
			job.ResultChan <- result
		}
	}
}

// processUpload ejecuta el upload con reintentos
func (wp *WorkerPool) processUpload(job UploadJob) UploadResult {
	var lastErr error
	
	for attempt := 0; attempt <= wp.maxRetries; attempt++ {
		if attempt > 0 {
			// Backoff exponencial entre reintentos
			backoff := time.Duration(attempt) * time.Second
			time.Sleep(backoff)
		}

		path, err := wp.executeUpload(job.FolderPath, job.FileName, job.Content)
		if err == nil {
			return UploadResult{Path: path, Error: nil}
		}

		lastErr = err
		
		// Si es un error de cliente (4xx), no reintentar
		if strings.Contains(err.Error(), "status: 4") {
			break
		}
	}

	return UploadResult{
		Path:  "",
		Error: fmt.Errorf("upload falló después de %d intentos: %w", wp.maxRetries+1, lastErr),
	}
}

// executeUpload realiza el upload a Nextcloud
func (wp *WorkerPool) executeUpload(folderPath, fileName string, content []byte) (string, error) {
	// Construir URL
	cleanFolderPath := strings.Trim(folderPath, "/")
	var fileURL string
	if cleanFolderPath == "" {
		fileURL = fmt.Sprintf("%s/remote.php/dav/files/%s/%s",
			wp.client.BaseURL, wp.client.Username, fileName)
	} else {
		fileURL = fmt.Sprintf("%s/remote.php/dav/files/%s/%s/%s",
			wp.client.BaseURL, wp.client.Username, cleanFolderPath, fileName)
	}

	// Crear contexto con timeout
	ctx, cancel := context.WithTimeout(wp.ctx, wp.client.RequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "PUT", fileURL, bytes.NewReader(content))
	if err != nil {
		return "", fmt.Errorf("error al crear request: %w", err)
	}

	req.SetBasicAuth(wp.client.Username, wp.client.Password)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Length", strconv.Itoa(len(content)))

	resp, err := wp.client.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error al subir archivo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("error al subir archivo a Nextcloud (status: %d): %s", resp.StatusCode, string(body))
	}

	return path.Join(cleanFolderPath, fileName), nil
}

// Submit envía un trabajo al pool
func (wp *WorkerPool) Submit(job UploadJob) {
	wp.jobs <- job
}

// Shutdown cierra el pool de workers
func (wp *WorkerPool) Shutdown() {
	close(wp.jobs)
	wp.wg.Wait()
	wp.cancel()
}

// Verificar que implementa la interfaz
var _ services.FileStorageService = (*NextcloudFileAdapter)(nil)

func NewNextcloudFileAdapter() *NextcloudFileAdapter {
	client := core.NewNextcloudClient()
	
	return &NextcloudFileAdapter{
		client:     client,
		workerPool: NewWorkerPool(client, client.MaxWorkers),
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

// UploadFile ahora usa el worker pool
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

	// Usar UploadFileFromBytes que ya implementa el worker pool
	return nf.UploadFileFromBytes(folderPath, fileName, fileContent)
}

// UploadFileFromBytes ahora usa el worker pool
func (nf *NextcloudFileAdapter) UploadFileFromBytes(folderPath string, fileName string, fileContent []byte) (string, error) {
	// Validar parámetros de entrada
	if fileName == "" {
		return "", fmt.Errorf("nombre de archivo no puede estar vacío")
	}
	if fileContent == nil {
		return "", fmt.Errorf("contenido del archivo no puede ser nil")
	}

	// Crear canal para resultado
	resultChan := make(chan UploadResult, 1)

	// Crear job y enviarlo al pool
	job := UploadJob{
		FolderPath: folderPath,
		FileName:   fileName,
		Content:    fileContent,
		ResultChan: resultChan,
	}

	nf.workerPool.Submit(job)

	// Esperar resultado con timeout
	ctx, cancel := context.WithTimeout(context.Background(), nf.client.RequestTimeout+5*time.Second)
	defer cancel()

	select {
	case result := <-resultChan:
		if result.Error != nil {
			return "", result.Error
		}
		return result.Path, nil
	case <-ctx.Done():
		return "", fmt.Errorf("timeout esperando resultado del upload")
	}
}

// DownloadFile - Sin cambios, ya que los downloads individuales son rápidos
func (nf *NextcloudFileAdapter) DownloadFile(folderPath string, fileName string) ([]byte, error) {
	// Validar parámetros de entrada
	if fileName == "" {
		return nil, fmt.Errorf("nombre de archivo no puede estar vacío")
	}

	// Construir la URL del archivo
	fileURL := nf.buildFileURL(folderPath, fileName)

	// Crear contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), nf.client.RequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", fileURL, nil)
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

// DeleteFile - Agregado contexto con timeout
func (nf *NextcloudFileAdapter) DeleteFile(folderPath string, fileName string) error {
	// Validar parámetros de entrada
	if fileName == "" {
		return fmt.Errorf("nombre de archivo no puede estar vacío")
	}

	// Construir la URL del archivo
	fileURL := nf.buildFileURL(folderPath, fileName)

	// Crear contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), nf.client.RequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "DELETE", fileURL, nil)
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

// FileExists - Agregado contexto con timeout
func (nf *NextcloudFileAdapter) FileExists(folderPath string, fileName string) (bool, error) {
	// Validar parámetros de entrada
	if fileName == "" {
		return false, fmt.Errorf("nombre de archivo no puede estar vacío")
	}

	// Construir la URL del archivo
	fileURL := nf.buildFileURL(folderPath, fileName)

	// Crear contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), nf.client.RequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "HEAD", fileURL, nil)
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

// GetFileInfo - Agregado contexto con timeout
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

	// Crear contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), nf.client.RequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "PROPFIND", fileURL, strings.NewReader(propfindXML))
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

// Cleanup cierra el worker pool (llamar al cerrar la aplicación)
func (nf *NextcloudFileAdapter) Cleanup() {
	if nf.workerPool != nil {
		nf.workerPool.Shutdown()
	}
}