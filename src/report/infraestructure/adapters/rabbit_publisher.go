package adapters

import (
	"context"
	"encoding/json"
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
	if core.RabbitChannel == nil {
		log.Println("no se conecto a rabbit")
		return nil
	}

	body, _ := json.Marshal(report)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := core.RabbitChannel.PublishWithContext(
		ctx,
		"reports",
		"process",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Println("error al enviar el reporte", err)
	} else {
		log.Println("se envio el reporte", string(body))
	}

	return err
}
