package main

import (
	"log"
	"rabbitConsumer/src/report/infraestructure"

	"github.com/gin-gonic/gin"
)

func main() {
	db := infraestructure.NewMySQL()

	deps, err := infraestructure.NewDependencies(db)
	if err != nil {
		log.Fatal("Error al inicializar dependencias:", err)
	}

	r := gin.Default()

	infraestructure.RegisterRoutes(r, deps.ProcessReportUseCase)

	r.Run(":8081")
}
