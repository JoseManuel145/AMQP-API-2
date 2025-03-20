package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"rabbitConsumer/src/core"
	"rabbitConsumer/src/report/domain/entities"
	"rabbitConsumer/src/report/domain/repositories"

	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQAdapter struct {
	conn    *core.Conn_RabbitMQ
	channel *amqp091.Channel
}

func NewRabbitMQAdapter(conn *core.Conn_RabbitMQ) repositories.IRabbitMQService {
	if conn.Err != "" {
		log.Fatalf("Error al configurar RabbitMQ: %v", conn.Err)
	}
	return &RabbitMQAdapter{conn: conn, channel: conn.Channel}
}

func (r *RabbitMQAdapter) FetchReports() ([]map[string]interface{}, error) {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Printf("Error cargando el archivo .env: %v", err)
	}

	reportsURL := os.Getenv("REPORTS_API_URL")
	if reportsURL == "" {
		return nil, errors.New("REPORTS_API_URL no está configurado en .env")
	}

	// Reintento en caso de error
	var reports []map[string]interface{}
	var err error
	for i := 0; i < 3; i++ {
		reports, err = r.fetchFromAPI(reportsURL)
		if err == nil {
			break
		}
		log.Printf("Error obteniendo reportes, reintento %d/3: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, err
	}

	return reports, nil
}

func (r *RabbitMQAdapter) PublishReport(report entities.Report) error {
	body, _ := json.Marshal(report)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.channel.PublishWithContext(
		ctx,
		"reports", // Exchange
		"process", // Routing Key
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Println("Error enviando mensaje a RabbitMQ:", err)
	} else {
		log.Println("Reporte enviado a RabbitMQ:", string(body))
	}

	return err
}

func (r *RabbitMQAdapter) fetchFromAPI(url string) ([]map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error obteniendo reportes del API:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("API externa no respondió con estado OK")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error leyendo respuesta del API:", err)
		return nil, err
	}

	var reports []map[string]interface{}
	if err := json.Unmarshal(body, &reports); err != nil {
		log.Println("Error decodificando JSON:", err)
		return nil, err
	}

	return reports, nil
}
