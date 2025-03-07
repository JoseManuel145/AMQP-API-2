package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"rabbitConsumer/src/core"
	"rabbitConsumer/src/report/domain/entities"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublishService struct{}

func NewRabbitMQPublishService() *RabbitMQPublishService {
	return &RabbitMQPublishService{}
}

func (s *RabbitMQPublishService) PublishReport(report *entities.Report) error {
	// Verificar si la conexión RabbitMQ está activa
	if core.RabbitChannel == nil {
		log.Println("Error: RabbitMQ no está conectado")
		return errors.New("RabbitMQ no está conectado")
	}

	// Serializar el reporte a formato JSON
	body, err := json.Marshal(report)
	if err != nil {
		log.Println("Error al serializar el reporte:", err)
		return err
	}

	// Configurar el contexto para el timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Publicar el mensaje en la cola "reports"
	err = core.RabbitChannel.PublishWithContext(
		ctx,
		"reports", // Nombre del exchange
		"process", // Nombre de la cola
		false,     // Mandatory
		false,     // Immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Println("Error al enviar el reporte a RabbitMQ:", err)
		return err
	}

	log.Println("Reporte enviado correctamente a la cola 'process'.")
	return nil
}
