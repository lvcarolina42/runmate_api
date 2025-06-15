package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"runmate_api/config"
	"runmate_api/http/handler"
	"runmate_api/internal/entity"
	"runmate_api/internal/repository"
	"runmate_api/internal/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("failed to load env %v", err)
	}

	db, err := gorm.Open(postgres.Open(config.DatabaseURL()), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database %v", err)
	}

	err = db.AutoMigrate(&entity.Activity{}, &entity.Coordinate{})
	if err != nil {
		log.Fatalf("failed to migrate database %v", err)
	}

	activityRepo := repository.NewActivity(db)
	activityService := service.NewActivity(activityRepo)
	api := handler.NewAPI(activityService)

	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.RealIP, middleware.Recoverer, middleware.RequestID)
	api.Routes(r)

	port := ":" + config.APIPort()
	log.Printf("Listening on %s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
