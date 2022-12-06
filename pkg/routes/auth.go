package routes

import (
	"github.com/gorilla/mux"
	"github.com/ong-gtp/go-chat/pkg/controllers"
	"github.com/ong-gtp/go-chat/pkg/domain/middlewares"
	"github.com/ong-gtp/go-chat/pkg/services"
)

var RegisterAuthRoutes = func(router *mux.Router) {

	sb := router.PathPrefix("/v1/api/auth").Subrouter()
	sb.Use(middlewares.HeaderMiddleware)

	var auth controllers.AuthController
	auth.RegisterService(services.NewAuthService())

	sb.HandleFunc("/login", auth.Login).Methods("POST")
	sb.HandleFunc("/signup", auth.SignUp).Methods("POST")
}
