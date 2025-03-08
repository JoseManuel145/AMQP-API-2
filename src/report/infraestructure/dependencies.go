package infraestructure

import (
	"log"
	"rabbitConsumer/src/core"
	"rabbitConsumer/src/report/application/usecases"
	_ "rabbitConsumer/src/report/domain"
	"rabbitConsumer/src/report/infraestructure/adapters"
)

type Dependencies struct {
	ProcessReportUseCase *usecases.ProcessReport
}

func NewDependencies(db *MySQL) (*Dependencies, error) {
	err := core.InitRabbitMQ()
	if err != nil {
		log.Fatal("Error al inicializar RabbitMQ:", err)
	}

	rabbitService := adapters.NewRabbitMQService()
	rabbitPublishService := adapters.NewRabbitMQPublishService()
	reportRepo := NewMySQL()

	// Crear el caso de uso de procesamiento de reportes
	processReportUseCase := usecases.NewProcessReport(reportRepo, rabbitService, rabbitPublishService)

	// Iniciar la escucha de reportes pendientes en un goroutine
	go processReportUseCase.StartProcessingReports()

	return &Dependencies{
		ProcessReportUseCase: processReportUseCase,
	}, nil
}
