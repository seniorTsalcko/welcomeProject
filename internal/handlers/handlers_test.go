package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"welcomeProject/internal/models"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateTask(task *models.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockRepository) GetAllTasks() ([]models.Task, error) {
	args := m.Called()
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *MockRepository) GetTaskByID(id int) (*models.Task, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockRepository) UpdateTask(task *models.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockRepository) DeleteTask(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) UpdateTaskStatus(id int, status string) (*models.Task, error) {
	args := m.Called(id, status)
	return args.Get(0).(*models.Task), args.Error(1)
}

func TestHelloHandler(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/hello", nil)
	responce := httptest.NewRecorder()
	handler := NewHandlers(nil)
	handler.HelloHandler(responce, request)
	assert.Equal(t, http.StatusOK, responce.Code)
	assert.JSONEq(t, `{"hello": "world"}`, responce.Body.String())
}

func TestCreateTaskHandler(t *testing.T) {
	mockRepo := &MockRepository{}
	handler := NewHandlers(mockRepo)

	task := models.Task{
		Description: "Hello World",
		Status:      "new",
	}

	mockRepo.On("CreateTask", mock.AnythingOfType("*models.Task")).Return(nil)

	body, _ := json.Marshal(task)
	request := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	responce := httptest.NewRecorder()
	handler.CreateTaskHandler(responce, request)
	assert.Equal(t, http.StatusCreated, responce.Code)
	mockRepo.AssertExpectations(t)
}

func TestGetTaskHandler(t *testing.T) {
	mockRepo := &MockRepository{}
	handler := NewHandlers(mockRepo)

	testTask := &models.Task{
		ID:          1,
		Description: "Hello World",
		Status:      "new",
	}

	mockRepo.On("GetTaskByID", testTask.ID).Return(testTask, nil)

	request := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	responce := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/tasks/{id}", handler.GetTaskHandler)
	router.ServeHTTP(responce, request)

	assert.Equal(t, http.StatusOK, responce.Code)
	mockRepo.AssertExpectations(t)
}

func TestUpdateTaskStatusHandler(t *testing.T) {
	mockRepo := new(MockRepository)
	handler := NewHandlers(mockRepo)

	statusUpdate := models.StatusUpdate{Status: "done"}
	updatedTask := &models.Task{
		ID:          1,
		Description: "Hello World",
		Status:      "done",
	}

	mockRepo.On("UpdateTaskStatus", 1, "done").Return(updatedTask, nil)

	jsonStatus, _ := json.Marshal(statusUpdate)
	request, _ := http.NewRequest(http.MethodPatch, "/tasks/1/status", bytes.NewBuffer(jsonStatus))
	request.Header.Set("Content-Type", "application/json")

	responce := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/tasks/{id}/status", handler.UpdateTaskStatusHandler).Methods("PATCH")
	router.ServeHTTP(responce, request)

	assert.Equal(t, http.StatusOK, responce.Code)
	mockRepo.AssertExpectations(t)
}
