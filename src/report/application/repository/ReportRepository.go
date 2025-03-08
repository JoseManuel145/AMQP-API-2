package usecases

import (
	"rabbitConsumer/src/report/domain"
	"rabbitConsumer/src/report/domain/entities"
)

type ProcessReport struct {
	repo domain.IReport
}

func NewProcessReport(repo domain.IReport) *ProcessReport {
	return &ProcessReport{repo: repo}
}

func (pr *ProcessReport) Execute(id int, title, content, status string) error {
	// Crear una instancia del Report con el nuevo estado
	report := entities.Report{
		ID:      id,
		Title:   title,
		Content: content,
		Status:  status,
	}

	// Guardar el reporte en la base de datos
	err := pr.repo.Update(report.ID, report.Title, report.Content, report.Status)
	if err != nil {
		return err
	}

	return nil
}
