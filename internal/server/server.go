package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"welcomeProject/internal/auth"
	"welcomeProject/internal/config"
	"welcomeProject/internal/handlers"
	"welcomeProject/internal/repository"
	"welcomeProject/middleware"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Server struct {
	router *mux.Router
	db     *sql.DB
}

func NewServer(dbConfig config.DBConfig, jwtConfig config.JWTConfig) *Server {
	s := &Server{
		router: mux.NewRouter(),
	}

	s.configureDB(dbConfig)
	s.configureRouter(jwtConfig.Secret)

	return s
}

func (s *Server) configureRouter(jwtSecret string) {
	repo := repository.NewRepository(s.db)
	authRepo := auth.NewAuthRepository(s.db)

	authService := auth.NewAuthService(authRepo, jwtSecret)

	taskHandlers := handlers.NewHandlers(repo)
	authHandlers := auth.NewAuthHandlers(authService)

	s.router.HandleFunc("/hello", taskHandlers.HelloHandler).Methods("GET")
	s.router.HandleFunc("/signup", authHandlers.SignUp).Methods("POST")
	s.router.HandleFunc("/login", authHandlers.Login).Methods("POST")

	protected := s.router.PathPrefix("/api/").Subrouter()
	protected.Use(middleware.JWTAuth(jwtSecret))

	protected.HandleFunc("/tasks", taskHandlers.CreateTaskHandler).Methods("POST")
	protected.HandleFunc("/tasks", taskHandlers.GetTasksHandler).Methods("GET")
	protected.HandleFunc("/tasks/{id}", taskHandlers.GetTaskHandler).Methods("GET")
	protected.HandleFunc("/tasks/{id}", taskHandlers.UpdateTaskHandler).Methods("PUT")
	protected.HandleFunc("/tasks/{id}", taskHandlers.DeleteTaskHandler).Methods("DELETE")
	protected.HandleFunc("/tasks/{id}/status", taskHandlers.UpdateTaskStatusHandler).Methods("PATCH")
}

func (s *Server) configureDB(dbConfig config.DBConfig) {
	var err error
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Database,
		dbConfig.Host,
		dbConfig.Port,
	)

	s.db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	if err = s.db.Ping(); err != nil {
		panic(err)
	}

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
		    id SERIAL PRIMARY KEY,
		    name TEXT NOT NULL,
		    email TEXT NOT NULL UNIQUE,
		    login TEXT NOT NULL UNIQUE,
		    password TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			description TEXT NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'new',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CHECK (status IN ('new', 'in progress', 'done'))
		)

`)
	if err != nil {
		panic(err)
	}
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
