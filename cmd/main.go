package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"runmate_api/config"
	"runmate_api/entity"
	"runmate_api/handler"
	"runmate_api/repository"
	"runmate_api/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load env %v", err)
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
	r.Post("/activities", api.CreateActivity)
	r.Get("/activities", api.GetActivities)
	r.Delete("/activities/{id}", api.DeleteActivity)

	r.Get("/users/{id}/activities", api.GetActivitiesByUser)

	port := ":" + config.APIPort()
	log.Printf("Listening on %s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
