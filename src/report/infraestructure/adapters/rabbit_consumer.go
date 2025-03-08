package adapters

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type RabbitMQService struct{}

func NewRabbitMQService() *RabbitMQService {
	return &RabbitMQService{}
}

func (s *RabbitMQService) FetchReports() ([]map[string]interface{}, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error cargando el archivo .env: %v", err)
	}

	alertsURL := os.Getenv("REPORTS_API_URL")

	resp, err := http.Get(alertsURL)
	if err != nil {
		log.Println("Error obteniendo alertas del consumidor:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error leyendo respuesta del consumidor:", err)
		return nil, err
	}

	var alerts []map[string]interface{}
	if err := json.Unmarshal(body, &alerts); err != nil {
		log.Println("Error decodificando JSON:", err)
		return nil, err
	}

	return alerts, nil
}
