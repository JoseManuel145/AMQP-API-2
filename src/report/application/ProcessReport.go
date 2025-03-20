package application

import (
	"log"
	"time"

	"rabbitConsumer/src/report/domain/entities"
	"rabbitConsumer/src/report/domain/repositories"
)

type ProcessReport struct {
	Repo       repositories.IReport
	RabbitMQ   repositories.IRabbitMQService
	UpdateRepo *UpdateReport
}

func NewProcessReport(DB repositories.IReport, Rabbit repositories.IRabbitMQService, updateRepo *UpdateReport) *ProcessReport {
	return &ProcessReport{
		Repo:       DB,
		RabbitMQ:   Rabbit,
		UpdateRepo: updateRepo,
	}
}

func (pr *ProcessReport) StartProcessingReports() {
	for {
		time.Sleep(5 * time.Second)

		reports, err := pr.RabbitMQ.FetchReports()
		if err != nil {
			log.Println("Error obteniendo reportes pendientes:", err)
			continue
		}

		for _, data := range reports {
			report := entities.Report{
				ID:      int(data["id"].(float64)),
				Title:   data["title"].(string),
				Content: data["content"].(string),
				Status:  "in process",
			}

			pr.UpdateRepo.Execute(report)

			err := pr.RabbitMQ.PublishReport(report)
			if err != nil {
				log.Println("Error enviando mensaje a RabbitMQ:", err)
			} else {
				log.Printf("Reporte ID %d procesado y enviado a la cola", report.ID)
			}
		}
	}
}
