package main

import (
	"03/config"
	"03/internal/middleware"
	"03/pkg/db/postgres"
	"database/sql"
	"log"

	iris "github.com/kataras/iris/v12"

	dialogRepo "03/internal/dialog/repository"
	groqBusiness "03/internal/groq_client/business"
	groqDelivery "03/internal/groq_client/delivery/http"
	wordRepo "03/internal/word/repository"
	"time"
)

func main() {
	log.Println("Start!")

	//Load config file
	configPath := "config/config"

	cfgFile, err := config.LoadConfig(configPath)

	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	//Initialize database
	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		log.Printf("Postgresql init: %s\n", err)
	} else {
		sqlDB, err := psqlDB.DB()
		if err != nil {
			log.Printf("Postgresql to sql.DB error: %s\n", err)
		}

		if sqlDB != nil {
			log.Println("Postgres connected, Status: %#v", sqlDB.Stats())
			defer func(sqlDB *sql.DB) {
				err := sqlDB.Close()
				if err != nil {
					log.Println("Postgres close error: %s", err)
				}
			}(sqlDB)

		}
	}

	// Initialize Iris application
	app := iris.New()
	midd := middleware.InitMiddleware()
	app.Use(midd.Cors)

	wr := wordRepo.NewWordRepository(psqlDB)
	dr := dialogRepo.NewDialogRepository(psqlDB)

	timeoutContext := time.Duration(5) * time.Second
	groqBusiness := groqBusiness.NewGroqBusiness(dr, wr, timeoutContext)
	groqDelivery.NewGroqHandler(app, groqBusiness)
	// Set the views directory
	app.RegisterView(iris.HTML("./views", ".html"))

	// Register a route to serve the HTML form
	app.Get("/", func(ctx iris.Context) {
		ctx.View("index.html")
	})

	app.Run(iris.Addr(cfg.Server.Port))
}
