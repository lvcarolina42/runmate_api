package main

import (
	"log"
	"net/http"

	"runmate_api/config"
	"runmate_api/http/handler"
	"runmate_api/internal/chat"
	"runmate_api/internal/entity"
	"runmate_api/internal/repository"
	"runmate_api/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	err = db.AutoMigrate(
		&entity.User{},
		&entity.Activity{},
		&entity.Coordinate{},
		&entity.Challenge{},
		&entity.ChallengeEvent{},
		&entity.Message{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database %v", err)
	}

	activityRepo := repository.NewActivity(db)
	challengeRepo := repository.NewChallenge(db)
	messageRepo := repository.NewMessage(db)
	userRepo := repository.NewUser(db)

	activityService := service.NewActivity(activityRepo, challengeRepo, userRepo)
	challengeService := service.NewChallenge(challengeRepo, userRepo)
	messageService := service.NewMessage(challengeRepo, messageRepo, userRepo)
	userService := service.NewUser(userRepo)

	chatHub := chat.NewHub()
	chatConsumer := chat.NewConsumer(chatHub, messageService)

	api := handler.NewAPI(activityService, challengeService, userService)
	chat := handler.NewChat(activityService, challengeService, messageService, userService, chatHub, chatConsumer)

	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.RealIP, middleware.Recoverer, middleware.RequestID)
	api.Routes(r)
	chat.Routes(r)

	port := ":" + config.APIPort()
	log.Printf("Listening on %s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
