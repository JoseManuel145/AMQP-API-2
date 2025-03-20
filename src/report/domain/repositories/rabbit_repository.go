package repositories

import "rabbitConsumer/src/report/domain/entities"

type IRabbitMQService interface {
	FetchReports() ([]map[string]interface{}, error)
	PublishReport(report entities.Report) error
}
