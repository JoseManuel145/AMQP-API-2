package adapters

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type RabbitMQService struct{}

func NewRabbitMQService() *RabbitMQService {
	return &RabbitMQService{}
}

func (s *RabbitMQService) FetchReports() ([]map[string]interface{}, error) {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error cargando el archivo .env: %v", err)
	}

	alertsURL := os.Getenv("REPORTS_API_URL")
	if alertsURL == "" {
		return nil, errors.New("REPORTS_API_URL no está configurado en .env")
	}

	// Reintento en caso de error
	var reports []map[string]interface{}
	var err error
	for i := 0; i < 3; i++ {
		reports, err = s.fetchFromAPI(alertsURL)
		if err == nil {
			break
		}
		log.Printf("Error obteniendo alertas, reintento %d/3: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, err
	}

	return reports, nil
}

func (s *RabbitMQService) fetchFromAPI(url string) ([]map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error obteniendo alertas del consumidor:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("API externa no respondió con estado OK")
	}

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
