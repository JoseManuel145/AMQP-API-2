package core

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
)

type Conn_RabbitMQ struct {
	Conn    *amqp091.Connection
	Channel *amqp091.Channel
	Err     string
}

var RabbitInstance *Conn_RabbitMQ

func InitRabbitMQ() *Conn_RabbitMQ {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error al cargar el archivo .env: %v", err)
	}

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Fatal("RABBITMQ_URL no definido en el archivo .env")
	}

	conn, err := amqp091.Dial(rabbitURL)
	if err != nil {
		return &Conn_RabbitMQ{Conn: nil, Channel: nil, Err: fmt.Sprintf("Error conectando con RabbitMQ: %v", err)}
	}

	channel, err := conn.Channel()
	if err != nil {
		return &Conn_RabbitMQ{Conn: conn, Channel: nil, Err: fmt.Sprintf("Error creando el canal de RabbitMQ: %v", err)}
	}

	RabbitInstance = &Conn_RabbitMQ{Conn: conn, Channel: channel, Err: ""}
	return RabbitInstance
}
