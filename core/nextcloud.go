// core/nextcloud.go
package core

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type NextcloudClient struct {
	BaseURL  string
	Username string
	Password string
	Client   *http.Client
	// Configuración de concurrencia
	MaxWorkers      int
	RequestTimeout  time.Duration
	MaxRetries      int
}

// NextcloudConfig contiene la configuración del cliente
type NextcloudConfig struct {
	BaseURL        string
	Username       string
	Password       string
	MaxWorkers     int         
	RequestTimeout time.Duration
	MaxRetries     int          
}

func NewNextcloudClient() *NextcloudClient {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error al cargar el archivo .env: %v\n", err)
	}

	// Valores por defecto
	maxWorkers := 5
	requestTimeout := 120 * time.Second
	maxRetries := 3

	// Leer valores desde .env si existen
	if workers := os.Getenv("NEXTCLOUD_MAX_WORKERS"); workers != "" {
		if val, err := strconv.Atoi(workers); err == nil && val > 0 {
			maxWorkers = val
		}
	}

	if timeout := os.Getenv("NEXTCLOUD_REQUEST_TIMEOUT_SECONDS"); timeout != "" {
		if val, err := strconv.Atoi(timeout); err == nil && val > 0 {
			requestTimeout = time.Duration(val) * time.Second
		}
	}

	if retries := os.Getenv("NEXTCLOUD_MAX_RETRIES"); retries != "" {
		if val, err := strconv.Atoi(retries); err == nil && val >= 0 {
			maxRetries = val
		}
	}

	return &NextcloudClient{
		BaseURL:        os.Getenv("NEXTCLOUD_BASE_URL"),
		Username:       os.Getenv("NEXTCLOUD_USERNAME"),
		Password:       os.Getenv("NEXTCLOUD_PASSWORD"),
		MaxWorkers:     maxWorkers,
		RequestTimeout: requestTimeout,
		MaxRetries:     maxRetries,
		Client: &http.Client{
			Timeout: requestTimeout,
		},
	}
}

func NewNextcloudClientWithConfig(config NextcloudConfig) *NextcloudClient {
	return &NextcloudClient{
		BaseURL:        config.BaseURL,
		Username:       config.Username,
		Password:       config.Password,
		MaxWorkers:     config.MaxWorkers,
		RequestTimeout: config.RequestTimeout,
		MaxRetries:     config.MaxRetries,
		Client: &http.Client{
			Timeout: config.RequestTimeout,
		},
	}
}