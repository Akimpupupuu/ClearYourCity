package tasks_transport_http_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	core_domain "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/domain"
	core_jwt "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/jwt"
	core_logger "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/logger"
	core_zap_logger "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/logger/zap"
	tasks_service "github.com/Akimpupupuu/ClearYourCity/task-service/internal/feature/tasks/service"
	mock_tasks_service "github.com/Akimpupupuu/ClearYourCity/task-service/internal/feature/tasks/service/mocks"
	tasks_transport_http "github.com/Akimpupupuu/ClearYourCity/task-service/internal/feature/tasks/transport/http"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	createTaskID      = 1
	createVersion     = 1
	createTitle       = "Clean the park"
	createDescription = "Remove plastic bottles near the lake"
	createUserID      = 1
)

func TestCreateTask(t *testing.T) {
	createdAt := time.Date(2026, time.October, 4, 12, 0, 0, 0, time.UTC)

	ctl := gomock.NewController(t)

	ctx := context.Background()
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)
	ctx = core_logger.ToContext(ctx, log)
	claims, err := core_jwt.NewUserClaims(createUserID, time.Minute*15)
	require.NoError(t, err)
	ctx = core_jwt.ToContext(ctx, claims)

	config := core_jwt.NewConfigMust()
	tokenGenerator := core_jwt.NewTokenGenerator(config)
	repositoryOutput := core_domain.NewTask(createTaskID, createVersion, createUserID, createTitle, createDescription, core_domain.StatusCreated, createdAt, nil)
	validData := `{"title": "Clean the park", "description": "Remove plastic bottles near the lake"}`
	requestBody := bytes.NewBuffer([]byte(validData))

	repository := mock_tasks_service.NewMockTasksRepository(ctl)
	repository.EXPECT().CreateTask(ctx, gomock.Any(), gomock.Any()).Return(repositoryOutput, nil)
	redis := mock_tasks_service.NewMockTasksRedis(ctl)

	service := tasks_service.NewTasksService(repository, redis, log)
	handler := tasks_transport_http.NewTasksHandler(service, tokenGenerator)

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		requestBody,
	)
	req = req.WithContext(ctx)
	handler.CreateTask(rec, req)

	result := rec.Result()
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	expected := `{"id":1,"version":1,"user_id":1,"title":"Clean the park","description":"Remove plastic bottles near the lake","status":"created","created_at":"2026-10-04T12:00:00Z","completed_at":null}`
	assert.JSONEq(t, expected, string(data))
}

func TestCreateTaskNoClaims(t *testing.T) {
	ctl := gomock.NewController(t)

	ctx := context.Background()
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)
	ctx = core_logger.ToContext(ctx, log)

	config := core_jwt.NewConfigMust()
	tokenGenerator := core_jwt.NewTokenGenerator(config)
	validData := `{"title": "Clean the park", "description": "Remove plastic bottles near the lake"}`
	requestBody := bytes.NewBuffer([]byte(validData))

	repository := mock_tasks_service.NewMockTasksRepository(ctl)
	redis := mock_tasks_service.NewMockTasksRedis(ctl)

	service := tasks_service.NewTasksService(repository, redis, log)
	handler := tasks_transport_http.NewTasksHandler(service, tokenGenerator)

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		requestBody,
	)
	req = req.WithContext(ctx)
	handler.CreateTask(rec, req)

	result := rec.Result()
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	expected := `{"message":"failed to get task","error":"get user claims from context: unauthorized"}`
	assert.JSONEq(t, expected, string(data))
}

func TestCreateTaskDecodeError(t *testing.T) {
	ctl := gomock.NewController(t)

	ctx := context.Background()
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)
	ctx = core_logger.ToContext(ctx, log)
	claims, err := core_jwt.NewUserClaims(createUserID, time.Minute*15)
	require.NoError(t, err)
	ctx = core_jwt.ToContext(ctx, claims)

	config := core_jwt.NewConfigMust()
	tokenGenerator := core_jwt.NewTokenGenerator(config)

	invalidData := `{"title": "Clean the park", "description": "Remove plastic bottles near the lake}`
	requestBody := bytes.NewBuffer([]byte(invalidData))

	repository := mock_tasks_service.NewMockTasksRepository(ctl)
	redis := mock_tasks_service.NewMockTasksRedis(ctl)

	service := tasks_service.NewTasksService(repository, redis, log)
	handler := tasks_transport_http.NewTasksHandler(service, tokenGenerator)

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		requestBody,
	)
	req = req.WithContext(ctx)
	handler.CreateTask(rec, req)

	result := rec.Result()
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	expected := `{"error":"decode HTTP request: unexpected EOF: invalid argument", "message":"decode and validate HTTP request"}`
	assert.JSONEq(t, expected, string(data))
}

func TestCreateTaskValidationError(t *testing.T) {
	ctl := gomock.NewController(t)

	ctx := context.Background()
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)
	ctx = core_logger.ToContext(ctx, log)
	claims, err := core_jwt.NewUserClaims(createUserID, time.Minute*15)
	require.NoError(t, err)
	ctx = core_jwt.ToContext(ctx, claims)

	config := core_jwt.NewConfigMust()
	tokenGenerator := core_jwt.NewTokenGenerator(config)

	invalidData := `{"title": "Hi", "description": "Remove plastic bottles near the lake"}`
	requestBody := bytes.NewBuffer([]byte(invalidData))

	repository := mock_tasks_service.NewMockTasksRepository(ctl)
	redis := mock_tasks_service.NewMockTasksRedis(ctl)

	service := tasks_service.NewTasksService(repository, redis, log)
	handler := tasks_transport_http.NewTasksHandler(service, tokenGenerator)

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		requestBody,
	)
	req = req.WithContext(ctx)
	handler.CreateTask(rec, req)

	result := rec.Result()
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	expected := `{"error":"request validation: Key: 'CreateTaskRequest.Title' Error:Field validation for 'Title' failed on the 'min' tag: invalid argument", "message":"decode and validate HTTP request"}`
	assert.JSONEq(t, expected, string(data))
}

func TestCreateTaskServiceError(t *testing.T) {
	ctl := gomock.NewController(t)

	ctx := context.Background()
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)
	ctx = core_logger.ToContext(ctx, log)
	claims, err := core_jwt.NewUserClaims(createUserID, time.Minute*15)
	require.NoError(t, err)
	ctx = core_jwt.ToContext(ctx, claims)

	config := core_jwt.NewConfigMust()
	tokenGenerator := core_jwt.NewTokenGenerator(config)
	validData := `{"title": "Clean the park", "description": "Remove plastic bottles near the lake"}`
	requestBody := bytes.NewBuffer([]byte(validData))
	expErr := fmt.Errorf("sql: connection refused")

	repository := mock_tasks_service.NewMockTasksRepository(ctl)
	repository.EXPECT().CreateTask(ctx, gomock.Any(), gomock.Any()).Return(nil, expErr)
	redis := mock_tasks_service.NewMockTasksRedis(ctl)

	service := tasks_service.NewTasksService(repository, redis, log)
	handler := tasks_transport_http.NewTasksHandler(service, tokenGenerator)

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		requestBody,
	)
	req = req.WithContext(ctx)
	handler.CreateTask(rec, req)

	result := rec.Result()
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	expected := `{"error":"create task in repository: sql: connection refused", "message":"failed to create task"}`
	assert.JSONEq(t, expected, string(data))
}
