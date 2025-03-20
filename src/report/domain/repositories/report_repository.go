package repositories

import "rabbitConsumer/src/report/domain/entities"

type IReport interface {
	Update(id int, title, content, status string) error
	GetAll() ([]entities.Report, error)
}
