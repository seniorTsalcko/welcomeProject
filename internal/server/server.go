package server

import (
	"database/sql"
	"net/http"
	"welcomeProject/internal/handlers"
	"welcomeProject/internal/repository"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Server struct {
	router *mux.Router
	db     *sql.DB
}

func NewServer() *Server {
	s := &Server{
		router: mux.NewRouter(),
	}

	s.configureDB()
	s.configureRouter()

	return s
}

func (s *Server) configureRouter() {
	repo := repository.NewRepository(s.db)
	h := handlers.NewHandlers(repo)

	s.router.HandleFunc("/hello", h.HelloHandler).Methods("GET")
	s.router.HandleFunc("/tasks", h.CreateTaskHandler).Methods("POST")
	s.router.HandleFunc("/tasks", h.GetTasksHandler).Methods("GET")
	s.router.HandleFunc("/tasks/{id}", h.GetTaskHandler).Methods("GET")
	s.router.HandleFunc("/tasks/{id}", h.UpdateTaskHandler).Methods("PUT")
	s.router.HandleFunc("/tasks/{id}", h.DeleteTaskHandler).Methods("DELETE")
	s.router.HandleFunc("/tasks/{id}/status", h.UpdateTaskStatusHandler).Methods("PATCH")
}

func (s *Server) configureDB() {
	var err error
	connStr := "user=postgres password=postgres dbname=welcome_project sslmode=disable"
	s.db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	if err = s.db.Ping(); err != nil {
		panic(err)
	}

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			description TEXT NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'new',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CHECK (status IN ('new', 'in progress', 'done'))
		)`)
	if err != nil {
		panic(err)
	}
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
