package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/romanzipp/sprinter/models"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Stats struct {
	ChecksLast24Hours  int64
	SuccessLast24Hours int64
	ErrorLast24Hours   int64
	ChecksLast7Days    int64
	SuccessLast7Days   int64
	ErrorLast7Days     int64
}

type ChartDataPoint struct {
	Date  time.Time `json:"date"`
	Value int64     `json:"value"`
}

type ChartData struct {
	Points []ChartDataPoint `json:"points"`
}

func IndexController(c *gin.Context, db *gorm.DB) {
	var checks []models.Check
	db.Order("created_at DESC").Limit(200).Find(&checks)

	stats := Stats{}

	db.Table("checks").Where("created_at >= ?", time.Now().Add(-time.Hour*24)).Count(&stats.ChecksLast24Hours)
	db.Table("checks").Where("created_at >= ?", time.Now().Add(-time.Hour*24)).Where("success = 1").Count(&stats.SuccessLast24Hours)
	db.Table("checks").Where("created_at >= ?", time.Now().Add(-time.Hour*24)).Where("success = 0").Count(&stats.ErrorLast24Hours)

	db.Table("checks").Where("created_at >= ?", time.Now().Add(-time.Hour*24*7)).Count(&stats.ChecksLast7Days)
	db.Table("checks").Where("created_at >= ?", time.Now().Add(-time.Hour*24*7)).Where("success = 1").Count(&stats.SuccessLast7Days)
	db.Table("checks").Where("created_at >= ?", time.Now().Add(-time.Hour*24*7)).Where("success = 0").Count(&stats.ErrorLast7Days)

	chartData := ChartData{}

	for _, v := range checks {
		chartData.Points = append(chartData.Points, ChartDataPoint{v.CreatedAt, int64(v.PacketLoss)})
	}

	ch1, _ := json.Marshal(chartData)

	c.HTML(http.StatusOK, "index", gin.H{
		"title":  "Home",
		"checks": checks,
		"stats":  stats,
		"chart":  string(ch1),
	})
}
