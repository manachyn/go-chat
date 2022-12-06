package main

import (
	"fmt"
	stdlog "log"
	"net/http"
	"os"

	log "github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/ong-gtp/go-chat/pkg/config"
	"github.com/ong-gtp/go-chat/pkg/domain/middlewares"
	"github.com/ong-gtp/go-chat/pkg/intetrnal/rabbitmq"
	"github.com/ong-gtp/go-chat/pkg/models"
	"github.com/ong-gtp/go-chat/pkg/routes"
	"github.com/rs/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		stdlog.Fatal("Error loading .env file")
	}

	config.ConnectDB()
	db := config.GetDB()
	db.AutoMigrate(&models.User{}, &models.ChatRoom{}, &models.Chat{})

	conn, ch := rabbitmq.InitilizeBroker()
	defer conn.Close()
	defer ch.Close()
	port := os.Getenv("PORT")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		stdlog.Fatal("JWT Secret not set")
	}

	r := mux.NewRouter()
	routes.RegisterAuthRoutes(r)
	routes.RegisterChatRoutes(r)
	routes.RegisterWebsocketRoute(r)

	// Logging setup
	var logger log.Logger
	// Logfmt is a structured, key=val logging format that is easy to read and parse
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	// Direct any attempts to use Go's log package to our structured logger
	stdlog.SetOutput(log.NewStdlibAdapter(logger))
	// Log the timestamp (in UTC) and the callsite (file + line number) of the logging
	// call for debugging in the future.
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "loc", log.DefaultCaller)
	loggingMiddleware := middlewares.LoggingMiddleware(logger)
	loggedRoutes := loggingMiddleware(r)

	logger.Log("Server", "starting", "port", port)
	handler := corsSetup(loggedRoutes)
	stdlog.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), handler))

}

func corsSetup(loggedRoutes http.Handler) http.Handler {
	allowedHeaders := []string{"Authorization", "Content-Type"}
	corsDebug := os.Getenv("CORS_DEBUG")
	if corsDebug == "true" {
		return cors.New(cors.Options{Debug: true, AllowedHeaders: allowedHeaders}).Handler(loggedRoutes)
	}
	return cors.New(cors.Options{AllowedHeaders: allowedHeaders}).Handler(loggedRoutes)

}