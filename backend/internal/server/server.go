package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"backend/internal/database"
	"backend/internal/handler"
	"backend/internal/repository"
	"backend/internal/service"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	router	*mux.Router
	db 		*database.DB
}

func NewServer(db *database.DB) *Server {
	return &Server{
		router: mux.NewRouter(),
		db: db,
	}
}

func (s *Server) SetupRoutes() {
	// Инициализация репозиториев и сервисов
	userRepo := repository.NewUserRepository(s.db)
	userServices := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userServices)

	// API routes
	api := s.router.PathPrefix("/api/v1").Subrouter()

	// User routes
	api.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	api.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")

	// Health check
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}).Methods("GET")
}

func (s *Server) Start(port string) error {
	s.SetupRoutes()

	//CORS Setup for Frontend
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:8080"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(s.router)

	srv := &http.Server{
		Handler: handler,
		Addr: ":"+port,
		WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server starting on port %s", port)
    return srv.ListenAndServe()
}