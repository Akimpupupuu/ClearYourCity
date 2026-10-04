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
	patchStatusTaskID      = 1
	patchStatusVersion     = 1
	patchStatusTitle       = "Clean the park"
	patchStatusDescription = "Remove plastic bottles near the lake"
	patchStatusUserID      = 1
	patchStatusToken       = "40d668c1f2fd4f24a85dfed76d5e76b543540bca17bb928dad4af0bc71671d63"
	patchStatusStatus      = "in_progress"
)

func TestPatchStatus(t *testing.T) {
	createdAt := time.Date(2026, time.October, 4, 12, 0, 0, 0, time.UTC)

	ctl := gomock.NewController(t)

	ctx := context.Background()
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)
	ctx = core_logger.ToContext(ctx, log)
	claims, err := core_jwt.NewUserClaims(patchStatusUserID, time.Minute*15)
	require.NoError(t, err)
	ctx = core_jwt.ToContext(ctx, claims)

	config := core_jwt.NewConfigMust()
	tokenGenerator := core_jwt.NewTokenGenerator(config)
	repositoryOutput := core_domain.NewTask(patchStatusTaskID, patchStatusVersion, patchStatusUserID, patchStatusTitle, patchStatusDescription, core_domain.StatusCreated, createdAt, nil)
	repositoryResult := *repositoryOutput
	repositoryResult.Status = core_domain.StatusInProgress

	repository := mock_tasks_service.NewMockTasksRepository(ctl)
	redis := mock_tasks_service.NewMockTasksRedis(ctl)
	redis.EXPECT().GetAction(ctx, patchStatusToken).Return(patchStatusTaskID, nil)
	repository.EXPECT().GetTask(ctx, patchStatusTaskID).Return(repositoryOutput, nil)
	repository.EXPECT().PatchStatus(ctx, repositoryOutput).Return(&repositoryResult, nil)
	redis.EXPECT().ProlongAction(ctx, patchStatusToken).Return(nil)

	service := tasks_service.NewTasksService(repository, redis, log)
	handler := tasks_transport_http.NewTasksHandler(service, tokenGenerator)

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("/status?token=%s&status=%s", patchStatusToken, patchStatusStatus),
		nil,
	)
	req = req.WithContext(ctx)
	handler.PatchStatus(rec, req)

	result := rec.Result()
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	expected := `{"id":1,"version":1,"user_id":1,"title":"Clean the park","description":"Remove plastic bottles near the lake","status":"in_progress","created_at":"2026-10-04T12:00:00Z","completed_at":null}`
	assert.JSONEq(t, expected, string(data))
}

func TestPatchStatusEmptyTokenError(t *testing.T) {
	ctl := gomock.NewController(t)

	ctx := context.Background()
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)
	ctx = core_logger.ToContext(ctx, log)
	claims, err := core_jwt.NewUserClaims(patchStatusUserID, time.Minute*15)
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
		http.MethodPatch,
		fmt.Sprintf("/status?status=%s", patchStatusStatus),
		nil,
	)
	req = req.WithContext(ctx)
	handler.PatchStatus(rec, req)

	result := rec.Result()
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	expected := `{"error":"'token' query param can't be empty: invalid argument", "message":"query params are empty"}`
	assert.JSONEq(t, expected, string(data))
}

func TestPatchStatusEmptyStatusError(t *testing.T) {
	ctl := gomock.NewController(t)

	ctx := context.Background()
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)
	ctx = core_logger.ToContext(ctx, log)
	claims, err := core_jwt.NewUserClaims(patchStatusUserID, time.Minute*15)
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
		http.MethodPatch,
		fmt.Sprintf("/status?token=%s", patchStatusToken),
		nil,
	)
	req = req.WithContext(ctx)
	handler.PatchStatus(rec, req)

	result := rec.Result()
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	expected := `{"error":"'status' query param can't be empty: invalid argument", "message":"query params are empty"}`
	assert.JSONEq(t, expected, string(data))
}

func TestPatchStatusServiceError(t *testing.T) {
	ctl := gomock.NewController(t)

	ctx := context.Background()
	log, err := core_zap_logger.NewLogger()
	require.NoError(t, err)
	ctx = core_logger.ToContext(ctx, log)
	claims, err := core_jwt.NewUserClaims(patchStatusUserID, time.Minute*15)
	require.NoError(t, err)
	ctx = core_jwt.ToContext(ctx, claims)

	config := core_jwt.NewConfigMust()
	tokenGenerator := core_jwt.NewTokenGenerator(config)
	expErr := fmt.Errorf("sql: connection refused")

	repository := mock_tasks_service.NewMockTasksRepository(ctl)
	redis := mock_tasks_service.NewMockTasksRedis(ctl)
	redis.EXPECT().GetAction(ctx, patchStatusToken).Return(patchStatusTaskID, nil)
	repository.EXPECT().GetTask(ctx, patchStatusTaskID).Return(nil, expErr)

	service := tasks_service.NewTasksService(repository, redis, log)
	handler := tasks_transport_http.NewTasksHandler(service, tokenGenerator)

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("/status?token=%s&status=%s", patchStatusToken, patchStatusStatus),
		nil,
	)
	req = req.WithContext(ctx)
	handler.PatchStatus(rec, req)

	result := rec.Result()
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	expected := `{"error":"get task from repository: sql: connection refused", "message":"failed to patch status"}`
	assert.JSONEq(t, expected, string(data))
}
