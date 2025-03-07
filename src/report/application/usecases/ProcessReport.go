package usecases

import (
	"log"
	"rabbitConsumer/src/report/domain"
	"rabbitConsumer/src/report/domain/entities"
	"rabbitConsumer/src/report/infraestructure/adapters"
	"time"
)

type ProcessReport struct {
	Repo           domain.IReport                   // Interfaz del repositorio de reportes
	Rabbit         *adapters.RabbitMQService        // Servicio para interactuar con RabbitMQ
	PublishService *adapters.RabbitMQPublishService // Servicio para publicar mensajes en RabbitMQ
}

func NewProcessReport(repo domain.IReport, rabbit *adapters.RabbitMQService, publishService *adapters.RabbitMQPublishService) *ProcessReport {
	return &ProcessReport{
		Repo:           repo,
		Rabbit:         rabbit,
		PublishService: publishService,
	}
}

// StartProcessingReports: Método que procesa los reportes en intervalos.
func (pr *ProcessReport) StartProcessingReports() {
	for {
		time.Sleep(5 * time.Second)

		// Fetch reports pendientes desde RabbitMQ
		reports, err := pr.Rabbit.FetchReports()
		if err != nil {
			log.Println("Error obteniendo reportes pendientes:", err)
			continue
		}

		// Procesar cada reporte
		for _, data := range reports {
			report := entities.Report{
				ID:      int(data["id"].(float64)), // JSON lo devuelve como float64
				Title:   data["title"].(string),
				Content: data["content"].(string),
				Status:  "in process", // Cambiamos el estado a "processed"
			}

			log.Printf("Procesando reporte ID %d", report.ID)

			// Actualizar el reporte en la base de datos
			err := pr.Repo.Update(report.ID, report.Title, report.Content, report.Status)
			if err != nil {
				log.Println("Error actualizando el reporte:", err)
				continue
			}

			// Publicar un mensaje a RabbitMQ sobre el estado actualizado
			err = pr.PublishService.PublishReport(&report)
			if err != nil {
				log.Println("Error enviando mensaje a RabbitMQ:", err)
			} else {
				log.Printf("Reporte ID %d procesado y enviado a la cola", report.ID)
			}
		}
	}
}
