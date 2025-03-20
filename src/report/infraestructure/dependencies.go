package infraestructure

import (
	"rabbitConsumer/src/core"
	"rabbitConsumer/src/report/application"
	_ "rabbitConsumer/src/report/domain/repositories"
	"rabbitConsumer/src/report/infraestructure/adapters"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	ProcessReportUseCase *application.ProcessReport
}

func NewDependencies(router *gin.Engine, db *MySQL, rabbitConn *core.Conn_RabbitMQ) error {
	// Inicializar el servicio de RabbitMQ con la conexión correcta
	rabbitService := adapters.NewRabbitMQAdapter(rabbitConn)

	updateReport := application.NewUpdateReport(db)

	// Crear el caso de uso de procesamiento de reportes
	processReportUseCase := application.NewProcessReport(db, rabbitService, updateReport)

	RegisterRoutes(router, processReportUseCase)
	// Iniciar la escucha de reportes pendientes en un goroutine
	go processReportUseCase.StartProcessingReports()

	return nil
}
