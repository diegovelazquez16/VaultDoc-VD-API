//core/nextcloud.go
package core

import (
	
	"fmt"
	
	"net/http"
	"os"
	

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
