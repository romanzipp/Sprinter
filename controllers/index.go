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
	Points24h []ChartDataPoint `json:"24h"`
	Points7d  []ChartDataPoint `json:"7d"`
}

func IndexController(c *gin.Context, db *gorm.DB) {
	var checks []models.Check
	db.Order("created_at DESC").Limit(100).Find(&checks)

	stats := Stats{}

	db.Table("checks").Where("created_at >= ?", time.Now().Add(-time.Hour*24)).Count(&stats.ChecksLast24Hours)
	db.Table("checks").Where("created_at >= ?", time.Now().Add(-time.Hour*24)).Where("success = 1").Count(&stats.SuccessLast24Hours)
	db.Table("checks").Where("created_at >= ?", time.Now().Add(-time.Hour*24)).Where("success = 0").Count(&stats.ErrorLast24Hours)

	db.Table("checks").Where("created_at >= ?", time.Now().Add(-time.Hour*24*7)).Count(&stats.ChecksLast7Days)
	db.Table("checks").Where("created_at >= ?", time.Now().Add(-time.Hour*24*7)).Where("success = 1").Count(&stats.SuccessLast7Days)
	db.Table("checks").Where("created_at >= ?", time.Now().Add(-time.Hour*24*7)).Where("success = 0").Count(&stats.ErrorLast7Days)

	var checks24h []models.Check
	var checks7d []models.Check

	db.Order("created_at DESC").Group("check_id").Where("created_at >= ?", time.Now().Add(-time.Hour*24)).Find(&checks24h)

	q := db.Select(`
		id, 
		created_at, 
		updated_at, 
		deleted_at, 
		addr, 
		ip, 
		packet_loss,
		CAST(AVG(latency) AS INT) latency,
		success, 
		check_id,
  		CASE success
		  	WHEN 0
				THEN check_id
		  	WHEN 1
				THEN strftime('%Y-%m-%dT%H:00:00.000', created_at) || '_' || success 
			END group_col
    `)
	q.Order("created_at DESC")
	q.Group("group_col")
	q.Where("created_at >= ?", time.Now().Add(-time.Hour*24*7))
	q.Find(&checks7d)

	chartData := ChartData{}

	for _, v := range checks24h {
		chartData.Points24h = append(chartData.Points24h, ChartDataPoint{v.CreatedAt, int64(v.Latency / time.Millisecond)})
	}

	for _, v := range checks7d {
		chartData.Points7d = append(chartData.Points7d, ChartDataPoint{v.CreatedAt, int64(v.Latency / time.Millisecond)})
	}

	chartDataJson, _ := json.Marshal(chartData)

	c.HTML(http.StatusOK, "index", gin.H{
		"title":  "Home",
		"checks": checks,
		"stats":  stats,
		"chart":  string(chartDataJson),
	})
}
