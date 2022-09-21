package main

import (
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/romanzipp/sprinter/config"
	"github.com/romanzipp/sprinter/controllers"
	"github.com/romanzipp/sprinter/models"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)

	conf := makeConfig()
	db := makeDb()
	router := makeRouter(db)

	go makeCheckScheduler(db, conf)

	log.Fatal().Err(router.Run())
}

func makeConfig() config.Config {
	interval, _ := strconv.ParseInt(os.Getenv("INTERVAL"), 10, 64)
	pingHosts := strings.Split(os.Getenv("PING_HOSTS"), ",")

	return config.Config{
		Interval:  interval,
		PingHosts: pingHosts,
	}
}

func makeDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Check{})

	return db
}

func makeRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.HTMLRender = ginview.Default()

	router.Static("/dist", "./dist")
	router.Static("/assets", "./static")

	router.Use(logger.SetLogger())

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowWildcard: true,
	}))

	routes := router.Group("/")
	routes.GET("/", func(c *gin.Context) {
		controllers.IndexController(c, db)
	})

	return router
}

func makeCheckScheduler(db *gorm.DB, conf config.Config) {
	for {
		for _, host := range conf.PingHosts {
			log.Info().Msgf("ping check %s", host)
			go models.ExecPingCheck(host, db)
		}

		time.Sleep(time.Duration(conf.Interval) * time.Second)
	}
}
