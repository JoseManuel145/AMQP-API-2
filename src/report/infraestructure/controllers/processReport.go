package controllers

import (
	"net/http"
	"rabbitConsumer/src/report/domain/repositories"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	repo repositories.IReport
}

func NewReportController(repo repositories.IReport) *ReportController {
	return &ReportController{repo: repo}
}

func (ctrl *ReportController) GetAllReports(c *gin.Context) {
	alerts, err := ctrl.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo los reportes"})
		return
	}

	c.JSON(http.StatusOK, alerts)
}
