package main

import (
	"log"
	"rabbitConsumer/src/core"
	"rabbitConsumer/src/report/infraestructure"
	_ "time"

	_ "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	db := infraestructure.NewMySQL()
	rabbit := core.InitRabbitMQ()

	r := gin.Default()

	if err := infraestructure.NewDependencies(r, db, rabbit); err != nil {
		log.Fatal("Error al inicializar dependencias:", err)
	}

	r.Run(":8081")
}
