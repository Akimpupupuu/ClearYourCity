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
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	patchTaskTaskID                 = 1
	patchTaskVersion                = 1
	patchTaskUserID                 = 1
	patchTaskTitle                  = "Clean the park"
	patchTaskUpdatedTitle           = "Clean at the park"
	patchTaskDescription            = "Remove plastic bottles near the lake"
	patchTaskUpdatedDescription     = "Remove plastic bottles"
	patchTaskInvalidPathValueString = "invalid"
	patchTaskInvalidPathValue       = -1
)

func TestPatchTask(t *testing.T) {
	createdAt := time.Date(2026, time.October, 4, 12, 0, 0, 0, time.UTC)

	ctl := gomock.NewController(t)

	ctx := context.Background()
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)
	ctx = core_logger.ToContext(ctx, log)
	claims, err := core_jwt.NewUserClaims(patchTaskUserID, time.Minute*15)
	require.NoError(t, err)
	ctx = core_jwt.ToContext(ctx, claims)

	config := core_jwt.NewConfigMust()
	tokenGenerator := core_jwt.NewTokenGenerator(config)
	repositoryOutput := core_domain.NewTask(patchTaskTaskID, patchTaskVersion, patchTaskUserID, patchTaskTitle, patchTaskDescription, core_domain.StatusCreated, createdAt, nil)
	repositoryResult := *repositoryOutput
	repositoryResult.Title = patchTaskUpdatedTitle
	repositoryResult.Description = patchTaskUpdatedDescription

	validData := `{"title": "Clean at the park", "description": "Remove plastic bottles"}`
	requestBody := bytes.NewBuffer([]byte(validData))

	repository := mock_tasks_service.NewMockTasksRepository(ctl)
	redis := mock_tasks_service.NewMockTasksRedis(ctl)

	service := tasks_service.NewTasksService(repository, redis, log)
	handler := tasks_transport_http.NewTasksHandler(service, tokenGenerator)

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/%d", patchTaskTaskID),
		requestBody,
	)
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", fmt.Sprintf("%d", patchTaskTaskID))
	ctx = context.WithValue(ctx, chi.RouteCtxKey, routeContext)
	req = req.WithContext(ctx)

	repository.EXPECT().GetTask(ctx, gomock.Any()).Return(repositoryOutput, nil)
	repository.EXPECT().PatchTask(ctx, &repositoryResult).Return(&repositoryResult, nil)

	handler.PatchTask(rec, req)

	result := rec.Result()
	defer func() { _ = result.Body.Close() }()

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	expected := `{"id":1,"version":1,"user_id":1,"title":"Clean at the park","description":"Remove plastic bottles","status":"created","created_at":"2026-10-04T12:00:00Z","completed_at":null}`
	assert.JSONEq(t, expected, string(data))
}

func TestPatchTaskNoClaims(t *testing.T) {
	ctl := gomock.NewController(t)

	ctx := context.Background()
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)
	ctx = core_logger.ToContext(ctx, log)

	config := core_jwt.NewConfigMust()
	tokenGenerator := core_jwt.NewTokenGenerator(config)

	validData := `{"title": "Clean at the park", "description": "Remove plastic bottles"}`
	requestBody := bytes.NewBuffer([]byte(validData))

	repository := mock_tasks_service.NewMockTasksRepository(ctl)
	redis := mock_tasks_service.NewMockTasksRedis(ctl)

	service := tasks_service.NewTasksService(repository, redis, log)
	handler := tasks_transport_http.NewTasksHandler(service, tokenGenerator)

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/%d", patchTaskTaskID),
		requestBody,
	)
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", fmt.Sprintf("%d", patchTaskTaskID))
	ctx = context.WithValue(ctx, chi.RouteCtxKey, routeContext)
	req = req.WithContext(ctx)

	handler.PatchTask(rec, req)

	result := rec.Result()
	defer func() { _ = result.Body.Close() }()

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	expected := `{"error":"get user claims from context: unauthorized", "message":"failed to patch task"}`
	assert.JSONEq(t, expected, string(data))
}

func TestPatchTaskInvalidPathValueString(t *testing.T) {
	ctl := gomock.NewController(t)

	ctx := context.Background()
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)
	ctx = core_logger.ToContext(ctx, log)
	claims, err := core_jwt.NewUserClaims(patchTaskUserID, time.Minute*15)
	require.NoError(t, err)
	ctx = core_jwt.ToContext(ctx, claims)

	config := core_jwt.NewConfigMust()
	tokenGenerator := core_jwt.NewTokenGenerator(config)

	validData := `{"title": "Clean at the park", "description": "Remove plastic bottles"}`
	requestBody := bytes.NewBuffer([]byte(validData))

	repository := mock_tasks_service.NewMockTasksRepository(ctl)
	redis := mock_tasks_service.NewMockTasksRedis(ctl)

	service := tasks_service.NewTasksService(repository, redis, log)
	handler := tasks_transport_http.NewTasksHandler(service, tokenGenerator)

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/%s", patchTaskInvalidPathValueString),
		requestBody,
	)
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", fmt.Sprintf("%s", patchTaskInvalidPathValueString))
	ctx = context.WithValue(ctx, chi.RouteCtxKey, routeContext)
	req = req.WithContext(ctx)

	handler.PatchTask(rec, req)

	result := rec.Result()
	defer func() { _ = result.Body.Close() }()

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	expected := `{"error":"path value = 'invalid' by key = 'id' is not a valid integer: strconv.Atoi: parsing \"invalid\": invalid syntax: invalid argument", "message":"failed to patch task"}`
	assert.JSONEq(t, expected, string(data))
}

func TestPatchTaskInvalidPathValue(t *testing.T) {
	ctl := gomock.NewController(t)

	ctx := context.Background()
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)
	ctx = core_logger.ToContext(ctx, log)
	claims, err := core_jwt.NewUserClaims(patchTaskUserID, time.Minute*15)
	require.NoError(t, err)
	ctx = core_jwt.ToContext(ctx, claims)

	config := core_jwt.NewConfigMust()
	tokenGenerator := core_jwt.NewTokenGenerator(config)

	validData := `{"title": "Clean at the park", "description": "Remove plastic bottles"}`
	requestBody := bytes.NewBuffer([]byte(validData))

	repository := mock_tasks_service.NewMockTasksRepository(ctl)
	redis := mock_tasks_service.NewMockTasksRedis(ctl)

	service := tasks_service.NewTasksService(repository, redis, log)
	handler := tasks_transport_http.NewTasksHandler(service, tokenGenerator)

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/%d", patchTaskInvalidPathValue),
		requestBody,
	)
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", fmt.Sprintf("%d", patchTaskInvalidPathValue))
	ctx = context.WithValue(ctx, chi.RouteCtxKey, routeContext)
	req = req.WithContext(ctx)

	handler.PatchTask(rec, req)

	result := rec.Result()
	defer func() { _ = result.Body.Close() }()

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	expected := `{"error":"invalid id: invalid argument", "message":"failed to patch task"}`
	assert.JSONEq(t, expected, string(data))
}

func TestPatchTaskDecodeError(t *testing.T) {
	ctl := gomock.NewController(t)

	ctx := context.Background()
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)
	ctx = core_logger.ToContext(ctx, log)
	claims, err := core_jwt.NewUserClaims(patchTaskUserID, time.Minute*15)
	require.NoError(t, err)
	ctx = core_jwt.ToContext(ctx, claims)

	config := core_jwt.NewConfigMust()
	tokenGenerator := core_jwt.NewTokenGenerator(config)

	validData := `{"title": "Clean at the park", "description": "Remove plastic bottles"`
	requestBody := bytes.NewBuffer([]byte(validData))

	repository := mock_tasks_service.NewMockTasksRepository(ctl)
	redis := mock_tasks_service.NewMockTasksRedis(ctl)

	service := tasks_service.NewTasksService(repository, redis, log)
	handler := tasks_transport_http.NewTasksHandler(service, tokenGenerator)

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/%d", patchTaskTaskID),
		requestBody,
	)
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", fmt.Sprintf("%d", patchTaskTaskID))
	ctx = context.WithValue(ctx, chi.RouteCtxKey, routeContext)
	req = req.WithContext(ctx)

	handler.PatchTask(rec, req)

	result := rec.Result()
	defer func() { _ = result.Body.Close() }()

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	expected := `{"error":"decode HTTP request: unexpected EOF: invalid argument", "message":"decode and validate HTTP request"}`
	assert.JSONEq(t, expected, string(data))
}

func TestPatchTaskValidationError(t *testing.T) {
	ctl := gomock.NewController(t)

	ctx := context.Background()
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)
	ctx = core_logger.ToContext(ctx, log)
	claims, err := core_jwt.NewUserClaims(patchTaskUserID, time.Minute*15)
	require.NoError(t, err)
	ctx = core_jwt.ToContext(ctx, claims)

	config := core_jwt.NewConfigMust()
	tokenGenerator := core_jwt.NewTokenGenerator(config)

	invalidData := `{"title": "Clean at the park", "description": "plastic"}`
	requestBody := bytes.NewBuffer([]byte(invalidData))

	repository := mock_tasks_service.NewMockTasksRepository(ctl)
	redis := mock_tasks_service.NewMockTasksRedis(ctl)

	service := tasks_service.NewTasksService(repository, redis, log)
	handler := tasks_transport_http.NewTasksHandler(service, tokenGenerator)

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/%d", patchTaskTaskID),
		requestBody,
	)
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", fmt.Sprintf("%d", patchTaskTaskID))
	ctx = context.WithValue(ctx, chi.RouteCtxKey, routeContext)
	req = req.WithContext(ctx)

	handler.PatchTask(rec, req)

	result := rec.Result()
	defer func() { _ = result.Body.Close() }()

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	expected := `{"error":"request validation: Key: 'PatchTaskRequest.Description' Error:Field validation for 'Description' failed on the 'min' tag: invalid argument", "message":"decode and validate HTTP request"}`
	assert.JSONEq(t, expected, string(data))
}

func TestPatchTaskServiceError(t *testing.T) {
	ctl := gomock.NewController(t)

	ctx := context.Background()
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)
	ctx = core_logger.ToContext(ctx, log)
	claims, err := core_jwt.NewUserClaims(patchTaskUserID, time.Minute*15)
	require.NoError(t, err)
	ctx = core_jwt.ToContext(ctx, claims)

	config := core_jwt.NewConfigMust()
	tokenGenerator := core_jwt.NewTokenGenerator(config)
	expErr := fmt.Errorf("sql: connection refused")

	validData := `{"title": "Clean at the park", "description": "Remove plastic bottles"}`
	requestBody := bytes.NewBuffer([]byte(validData))

	repository := mock_tasks_service.NewMockTasksRepository(ctl)
	redis := mock_tasks_service.NewMockTasksRedis(ctl)

	service := tasks_service.NewTasksService(repository, redis, log)
	handler := tasks_transport_http.NewTasksHandler(service, tokenGenerator)

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/%d", patchTaskTaskID),
		requestBody,
	)
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", fmt.Sprintf("%d", patchTaskTaskID))
	ctx = context.WithValue(ctx, chi.RouteCtxKey, routeContext)
	req = req.WithContext(ctx)

	repository.EXPECT().GetTask(ctx, gomock.Any()).Return(nil, expErr)

	handler.PatchTask(rec, req)

	result := rec.Result()
	defer func() { _ = result.Body.Close() }()

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	expected := `{"error":"get task from repository: sql: connection refused", "message":"failed to patch task"}`
	assert.JSONEq(t, expected, string(data))
}
