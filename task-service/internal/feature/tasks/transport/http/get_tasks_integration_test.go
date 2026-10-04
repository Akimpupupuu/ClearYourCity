package tasks_transport_http_test

import (
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
	getTaskIDFirst         = 1
	getTaskIDSecond        = 2
	getVersion             = 1
	getTitle               = "Clean the park"
	getDescription         = "Remove plastic bottles near the lake"
	getUserID              = 1
	getLimit               = 20
	getOffset              = 0
	getInvalidLimitString  = "limit"
	getInvalidOffsetString = "offset"
	getInvalidLimit        = -1
	getInvalidOffset       = -1
)

func TestGetTasks(t *testing.T) {
	createdAt := time.Date(2026, time.October, 4, 12, 0, 0, 0, time.UTC)

	ctl := gomock.NewController(t)

	ctx := context.Background()
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)
	ctx = core_logger.ToContext(ctx, log)
	claims, err := core_jwt.NewUserClaims(getUserID, time.Minute*15)
	require.NoError(t, err)
	ctx = core_jwt.ToContext(ctx, claims)

	config := core_jwt.NewConfigMust()
	tokenGenerator := core_jwt.NewTokenGenerator(config)
	repositoryOutput := []*core_domain.Task{
		core_domain.NewTask(getTaskIDFirst, getVersion, getUserID, getTitle, getDescription, core_domain.StatusCreated, createdAt, nil),
		core_domain.NewTask(getTaskIDSecond, getVersion, getUserID, getTitle, getDescription, core_domain.StatusCreated, createdAt, nil),
	}

	repository := mock_tasks_service.NewMockTasksRepository(ctl)
	repository.EXPECT().GetTasks(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(repositoryOutput, nil)
	redis := mock_tasks_service.NewMockTasksRedis(ctl)

	service := tasks_service.NewTasksService(repository, redis, log)
	handler := tasks_transport_http.NewTasksHandler(service, tokenGenerator)

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/?limit=%d&offset=%d", getLimit, getOffset),
		nil,
	)
	req = req.WithContext(ctx)
	handler.GetTasks(rec, req)

	result := rec.Result()
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	expected :=
		`[
		{"id":1,"version":1,"user_id":1,"title":"Clean the park","description":"Remove plastic bottles near the lake","status":"created","created_at":"2026-10-04T12:00:00Z","completed_at":null},
		{"id":2,"version":1,"user_id":1,"title":"Clean the park","description":"Remove plastic bottles near the lake","status":"created","created_at":"2026-10-04T12:00:00Z","completed_at":null}
		]`
	assert.JSONEq(t, expected, string(data))
}

func TestGetTasksNilQueryParams(t *testing.T) {
	createdAt := time.Date(2026, time.October, 4, 12, 0, 0, 0, time.UTC)

	ctl := gomock.NewController(t)

	ctx := context.Background()
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)
	ctx = core_logger.ToContext(ctx, log)
	claims, err := core_jwt.NewUserClaims(getUserID, time.Minute*15)
	require.NoError(t, err)
	ctx = core_jwt.ToContext(ctx, claims)

	config := core_jwt.NewConfigMust()
	tokenGenerator := core_jwt.NewTokenGenerator(config)
	repositoryOutput := []*core_domain.Task{
		core_domain.NewTask(getTaskIDFirst, getVersion, getUserID, getTitle, getDescription, core_domain.StatusCreated, createdAt, nil),
		core_domain.NewTask(getTaskIDSecond, getVersion, getUserID, getTitle, getDescription, core_domain.StatusCreated, createdAt, nil),
	}

	repository := mock_tasks_service.NewMockTasksRepository(ctl)
	repository.EXPECT().GetTasks(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(repositoryOutput, nil)
	redis := mock_tasks_service.NewMockTasksRedis(ctl)

	service := tasks_service.NewTasksService(repository, redis, log)
	handler := tasks_transport_http.NewTasksHandler(service, tokenGenerator)

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)
	req = req.WithContext(ctx)
	handler.GetTasks(rec, req)

	result := rec.Result()
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	expected :=
		`[
		{"id":1,"version":1,"user_id":1,"title":"Clean the park","description":"Remove plastic bottles near the lake","status":"created","created_at":"2026-10-04T12:00:00Z","completed_at":null},
		{"id":2,"version":1,"user_id":1,"title":"Clean the park","description":"Remove plastic bottles near the lake","status":"created","created_at":"2026-10-04T12:00:00Z","completed_at":null}
		]`
	assert.JSONEq(t, expected, string(data))
}

func TestGetTasksNoClaims(t *testing.T) {
	ctl := gomock.NewController(t)

	ctx := context.Background()
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)
	ctx = core_logger.ToContext(ctx, log)

	config := core_jwt.NewConfigMust()
	tokenGenerator := core_jwt.NewTokenGenerator(config)

	repository := mock_tasks_service.NewMockTasksRepository(ctl)
	redis := mock_tasks_service.NewMockTasksRedis(ctl)

	service := tasks_service.NewTasksService(repository, redis, log)
	handler := tasks_transport_http.NewTasksHandler(service, tokenGenerator)

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/?limit=%d&offset=%d", getLimit, getOffset),
		nil,
	)
	req = req.WithContext(ctx)
	handler.GetTasks(rec, req)

	result := rec.Result()
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	expected := `{"error":"get user claims from context: unauthorized", "message":"failed to get task"}`
	assert.JSONEq(t, expected, string(data))
}

func TestGetTasksInvalidLimitQueryParam(t *testing.T) {
	ctl := gomock.NewController(t)

	ctx := context.Background()
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)
	ctx = core_logger.ToContext(ctx, log)
	claims, err := core_jwt.NewUserClaims(getUserID, time.Minute*15)
	require.NoError(t, err)
	ctx = core_jwt.ToContext(ctx, claims)

	config := core_jwt.NewConfigMust()
	tokenGenerator := core_jwt.NewTokenGenerator(config)

	repository := mock_tasks_service.NewMockTasksRepository(ctl)
	redis := mock_tasks_service.NewMockTasksRedis(ctl)

	service := tasks_service.NewTasksService(repository, redis, log)
	handler := tasks_transport_http.NewTasksHandler(service, tokenGenerator)

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/?limit=%s&offset=%d", getInvalidLimitString, getOffset),
		nil,
	)
	req = req.WithContext(ctx)
	handler.GetTasks(rec, req)

	result := rec.Result()
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	expected := `{"error":"get limit query param: param: 'limit', by key: 'limit' is not a valid integer: strconv.Atoi: parsing \"limit\": invalid syntax: invalid argument", "message":"failed to get query params"}`
	assert.JSONEq(t, expected, string(data))
}

func TestGetTasksInvalidOffsetQueryParam(t *testing.T) {
	ctl := gomock.NewController(t)

	ctx := context.Background()
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)
	ctx = core_logger.ToContext(ctx, log)
	claims, err := core_jwt.NewUserClaims(getUserID, time.Minute*15)
	require.NoError(t, err)
	ctx = core_jwt.ToContext(ctx, claims)

	config := core_jwt.NewConfigMust()
	tokenGenerator := core_jwt.NewTokenGenerator(config)

	repository := mock_tasks_service.NewMockTasksRepository(ctl)
	redis := mock_tasks_service.NewMockTasksRedis(ctl)

	service := tasks_service.NewTasksService(repository, redis, log)
	handler := tasks_transport_http.NewTasksHandler(service, tokenGenerator)

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/?limit=%d&offset=%s", getLimit, getInvalidOffsetString),
		nil,
	)
	req = req.WithContext(ctx)
	handler.GetTasks(rec, req)

	result := rec.Result()
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	expected := `{"error":"get offset query param: param: 'offset', by key: 'offset' is not a valid integer: strconv.Atoi: parsing \"offset\": invalid syntax: invalid argument", "message":"failed to get query params"}`
	assert.JSONEq(t, expected, string(data))
}

func TestGetTasksValidateLimitQueryParamError(t *testing.T) {
	ctl := gomock.NewController(t)

	ctx := context.Background()
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)
	ctx = core_logger.ToContext(ctx, log)
	claims, err := core_jwt.NewUserClaims(getUserID, time.Minute*15)
	require.NoError(t, err)
	ctx = core_jwt.ToContext(ctx, claims)

	config := core_jwt.NewConfigMust()
	tokenGenerator := core_jwt.NewTokenGenerator(config)

	repository := mock_tasks_service.NewMockTasksRepository(ctl)
	redis := mock_tasks_service.NewMockTasksRedis(ctl)

	service := tasks_service.NewTasksService(repository, redis, log)
	handler := tasks_transport_http.NewTasksHandler(service, tokenGenerator)

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/?limit=%d&offset=%d", getInvalidLimit, getOffset),
		nil,
	)
	req = req.WithContext(ctx)
	handler.GetTasks(rec, req)

	result := rec.Result()
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	expected := `{"error":"validate pagination values: 'limit' must be non-negative value and smaller then 100: -1: invalid argument", "message":"failed to get tasks"}`
	assert.JSONEq(t, expected, string(data))
}

func TestGetTasksValidateOffsetQueryParamError(t *testing.T) {
	ctl := gomock.NewController(t)

	ctx := context.Background()
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)
	ctx = core_logger.ToContext(ctx, log)
	claims, err := core_jwt.NewUserClaims(getUserID, time.Minute*15)
	require.NoError(t, err)
	ctx = core_jwt.ToContext(ctx, claims)

	config := core_jwt.NewConfigMust()
	tokenGenerator := core_jwt.NewTokenGenerator(config)

	repository := mock_tasks_service.NewMockTasksRepository(ctl)
	redis := mock_tasks_service.NewMockTasksRedis(ctl)

	service := tasks_service.NewTasksService(repository, redis, log)
	handler := tasks_transport_http.NewTasksHandler(service, tokenGenerator)

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/?limit=%d&offset=%d", getLimit, getInvalidOffset),
		nil,
	)
	req = req.WithContext(ctx)
	handler.GetTasks(rec, req)

	result := rec.Result()
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	expected := `{"error":"validate pagination values: 'offset' must be non-negative value: -1: invalid argument", "message":"failed to get tasks"}`
	assert.JSONEq(t, expected, string(data))
}

func TestGetTasksServiceError(t *testing.T) {
	ctl := gomock.NewController(t)

	ctx := context.Background()
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)
	ctx = core_logger.ToContext(ctx, log)
	claims, err := core_jwt.NewUserClaims(getUserID, time.Minute*15)
	require.NoError(t, err)
	ctx = core_jwt.ToContext(ctx, claims)

	config := core_jwt.NewConfigMust()
	tokenGenerator := core_jwt.NewTokenGenerator(config)
	expErr := fmt.Errorf("sql: connection refused")

	repository := mock_tasks_service.NewMockTasksRepository(ctl)
	repository.EXPECT().GetTasks(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, expErr)
	redis := mock_tasks_service.NewMockTasksRedis(ctl)

	service := tasks_service.NewTasksService(repository, redis, log)
	handler := tasks_transport_http.NewTasksHandler(service, tokenGenerator)

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/?limit=%d&offset=%d", getLimit, getOffset),
		nil,
	)
	req = req.WithContext(ctx)
	handler.GetTasks(rec, req)

	result := rec.Result()
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	expected := `{"error":"get tasks from repository: sql: connection refused", "message":"failed to get tasks"}`
	assert.JSONEq(t, expected, string(data))
}
