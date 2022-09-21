package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/romanzipp/sprinter/models"
	"gorm.io/gorm"
	"net/http"
)

func IndexController(c *gin.Context, db *gorm.DB) {
	var checks []models.Check
	db.Order("created_at DESC").Limit(200).Find(&checks)

	c.HTML(http.StatusOK, "index", gin.H{
		"title":  "Home",
		"checks": checks,
	})
}
