package application

import (
	"log"
	"rabbitConsumer/src/report/domain/entities"
	"rabbitConsumer/src/report/domain/repositories"
)

type UpdateReport struct {
	Repo repositories.IReport
}

func NewUpdateReport(repo repositories.IReport) *UpdateReport {
	return &UpdateReport{
		Repo: repo,
	}
}

func (pr *UpdateReport) Execute(report entities.Report) {
	log.Printf("Procesando reporte ID %d", report.ID)

	err := pr.Repo.Update(report.ID, report.Title, report.Content, report.Status)
	if err != nil {
		log.Println("Error actualizando el reporte:", err)
	} else {
		log.Printf("Reporte ID %d actualizado correctamente", report.ID)
	}
}
