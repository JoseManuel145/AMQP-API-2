package infraestructure

import (
	"rabbitConsumer/src/report/application"
	"rabbitConsumer/src/report/infraestructure/controllers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, processReport *application.ProcessReport) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	reportController := controllers.NewReportController(processReport.Repo)

	router.GET("/reports", reportController.GetAllReports)
}
