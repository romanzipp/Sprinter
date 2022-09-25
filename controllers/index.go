package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/romanzipp/sprinter/models"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Stats struct {
	ChecksLast24Hours  int64
	SuccessLast24Hours int64
	ErrorLast24Hours   int64
	Uptime24Hours      float64
	ChecksLast7Days    int64
	SuccessLast7Days   int64
	ErrorLast7Days     int64
	Uptime7Days        float64
}

type ChartDataPoint struct {
	Date  time.Time `json:"date"`
	Value int64     `json:"value"`
}

type ChartData struct {
	Points24h []ChartDataPoint `json:"24h"`
	Points7d  []ChartDataPoint `json:"7d"`
}

const FilterAll = "all"
const FilterFailed = "failed"

func IndexController(c *gin.Context, db *gorm.DB) {
	var checks []models.Check

	filter := FilterAll
	if c.Query("filter") == "failed" {
		filter = FilterFailed
	}

	latestQuery := db.Order("created_at DESC").Limit(100)

	if filter == FilterFailed {
		latestQuery.Where("success = 0")
	}

	latestQuery.Find(&checks)

	stats := Stats{}

	db.Table("checks").Where("created_at >= ?", time.Now().Add(-time.Hour*24)).Count(&stats.ChecksLast24Hours)
	db.Table("checks").Where("created_at >= ?", time.Now().Add(-time.Hour*24)).Where("success = 1").Count(&stats.SuccessLast24Hours)
	db.Table("checks").Where("created_at >= ?", time.Now().Add(-time.Hour*24)).Where("success = 0").Count(&stats.ErrorLast24Hours)

	db.Table("checks").Where("created_at >= ?", time.Now().Add(-time.Hour*24*7)).Count(&stats.ChecksLast7Days)
	db.Table("checks").Where("created_at >= ?", time.Now().Add(-time.Hour*24*7)).Where("success = 1").Count(&stats.SuccessLast7Days)
	db.Table("checks").Where("created_at >= ?", time.Now().Add(-time.Hour*24*7)).Where("success = 0").Count(&stats.ErrorLast7Days)

	if stats.ChecksLast24Hours == 0 {
		stats.Uptime24Hours = 0.0
	} else {
		stats.Uptime24Hours = float64(stats.SuccessLast24Hours) / float64(stats.ChecksLast24Hours) * 100.0
	}

	if stats.ChecksLast7Days == 0 {
		stats.Uptime7Days = 0.0
	} else {
		stats.Uptime7Days = float64(stats.SuccessLast7Days) / float64(stats.ChecksLast7Days) * 100.0
	}

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
	up := true

	if len(checks24h) > 0 {
		up = checks24h[len(checks24h)-1].Success
	}

	c.HTML(http.StatusOK, "index", gin.H{
		"title":   "Home",
		"checks":  checks,
		"stats":   stats,
		"chart":   string(chartDataJson),
		"up":      up,
		"filter":  filter,
		"printer": message.NewPrinter(language.English),
	})
}
