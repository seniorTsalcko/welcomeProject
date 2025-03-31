package repository

import (
	"database/sql"
	"welcomeProject/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateTask(task *models.Task) error {
	sqlStatement := `
		INSERT INTO tasks (description, status)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(sqlStatement, task.Description, task.Status).
		Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
}

func (r *Repository) GetAllTasks() ([]models.Task, error) {
	rows, err := r.db.Query("SELECT id, description, status, created_at, updated_at FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *Repository) GetTaskByID(id int) (*models.Task, error) {
	var task models.Task
	err := r.db.QueryRow("SELECT id, description, status, created_at, updated_at FROM tasks WHERE id = $1", id).
		Scan(&task.ID, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *Repository) UpdateTask(task *models.Task) error {
	sqlStatement := `
		UPDATE tasks 
		SET description = $1, status = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING updated_at`

	return r.db.QueryRow(sqlStatement, task.Description, task.Status, task.ID).
		Scan(&task.UpdatedAt)
}

func (r *Repository) DeleteTask(id int) error {
	_, err := r.db.Exec("DELETE FROM tasks WHERE id = $1", id)
	return err
}

func (r *Repository) UpdateTaskStatus(id int, status string) (*models.Task, error) {
	var task models.Task
	sqlStatement := `
		UPDATE tasks 
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING id, description, status, created_at, updated_at`

	err := r.db.QueryRow(sqlStatement, status, id).
		Scan(&task.ID, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &task, nil
}
