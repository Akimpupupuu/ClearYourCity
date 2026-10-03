package core_domain

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testTaskInput struct {
	name   string
	input  *Task
	expErr error
}

func TestNewTask(t *testing.T) {
	id := 1
	version := 1
	userID := 1
	title := "Clean the park"
	description := "Remove plastic bottles near the lake"
	taskStatus := StatusCreated
	createdAt := time.Now()
	completedAt := createdAt.Add(time.Hour)

	task := NewTask(id, version, userID, title, description, taskStatus, createdAt, &completedAt)

	require.NotNil(t, task)
	require.Equal(t, id, task.ID)
	require.Equal(t, version, task.Version)
	require.Equal(t, userID, task.UserID)
	require.Equal(t, title, task.Title)
	require.Equal(t, description, task.Description)
	require.Equal(t, taskStatus, task.Status)
	require.Equal(t, createdAt, task.CreatedAt)
	require.Equal(t, &completedAt, task.CompletedAt)
}

func TestNewTaskUninitialized(t *testing.T) {
	userID := 1
	title := "Clean the park"
	description := "Remove plastic bottles near the lake"
	startTime := time.Now()

	task := NewTaskUninitialized(title, description, userID)

	require.NotNil(t, task)
	require.Equal(t, UninitializedID, task.ID)
	require.Equal(t, UninitializedVersion, task.Version)
	require.Equal(t, userID, task.UserID)
	require.Equal(t, title, task.Title)
	require.Equal(t, description, task.Description)
	require.Equal(t, StatusCreated, task.Status)
	require.WithinDuration(t, startTime, task.CreatedAt, time.Second)
	require.Nil(t, task.CompletedAt)
}

func TestValidate(t *testing.T) {
	const (
		taskID                  = 1
		userID                  = 1
		invalidUserID           = -1
		version                 = 1
		title                   = "Clean the park"
		invalidShortTitle       = "Hi"
		invalidLongTitle        = "1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111"
		description             = "Remove plastic bottles near the lake"
		invalidLongDescription  = "11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111"
		invalidShortDescription = "plastic"
	)

	createdAt := time.Date(2026, time.September, 28, 15, 0, 0, 0, time.UTC)
	completedAt := time.Date(2026, time.September, 28, 16, 0, 0, 0, time.UTC)
	completedAtInvalid := time.Date(2026, time.September, 28, 14, 0, 0, 0, time.UTC)

	inputs := []testTaskInput{
		{
			"success: comletedAt == nil",
			NewTask(taskID, version, userID, title, description, StatusCreated, createdAt, nil),
			nil,
		},
		{
			"success: comletedAt != nil",
			NewTask(taskID, version, userID, title, description, StatusDone, createdAt, &completedAt),
			nil,
		},
		{
			"fail: invalid title length: more",
			NewTask(taskID, version, userID, invalidLongTitle, description, StatusCreated, createdAt, nil),
			fmt.Errorf("invalid 'title' length: 103: invalid argument"),
		},
		{
			"fail: invalid title length: less",
			NewTask(taskID, version, userID, invalidShortTitle, description, StatusCreated, createdAt, nil),
			fmt.Errorf("invalid 'title' length: 2: invalid argument"),
		},
		{
			"fail: invalid description length: more",
			NewTask(taskID, version, userID, title, invalidLongDescription, StatusCreated, createdAt, nil),
			fmt.Errorf("invalid 'description' length: 1133: invalid argument"),
		},
		{
			"fail: invalid description length: less",
			NewTask(taskID, version, userID, title, invalidShortDescription, StatusCreated, createdAt, nil),
			fmt.Errorf("invalid 'description' length: 7: invalid argument"),
		},
		{
			"fail: invalid combination of status and completedAt values: status = created and comletedAt != nil",
			NewTask(taskID, version, userID, title, description, StatusCreated, createdAt, &completedAt),
			fmt.Errorf("'completed_at' must be nil if 'status' = 'created' or 'in_progress': invalid argument"),
		},
		{
			"fail: invalid combination of status and completedAt values: status = in_progress and comletedAt != nil",
			NewTask(taskID, version, userID, title, description, StatusInProgress, createdAt, &completedAt),
			fmt.Errorf("'completed_at' must be nil if 'status' = 'created' or 'in_progress': invalid argument"),
		},
		{
			"fail: invalid combination of status and completedAt values: status = done and comletedAt == nil",
			NewTask(taskID, version, userID, title, description, StatusDone, createdAt, nil),
			fmt.Errorf("'completed_at' can't be nil if 'status' = 'done' or 'rejected': invalid argument"),
		},
		{
			"fail: invalid combination of status and completedAt values: status = rejected and comletedAt == nil",
			NewTask(taskID, version, userID, title, description, StatusRejected, createdAt, nil),
			fmt.Errorf("'completed_at' can't be nil if 'status' = 'done' or 'rejected': invalid argument"),
		},
		{
			"fail: invalid userID",
			NewTask(taskID, version, invalidUserID, title, description, StatusCreated, createdAt, nil),
			fmt.Errorf("'user_id' can't be nil: invalid argument"),
		},
		{
			"fail: invalid combination of createdAt and completedAt: createdAt > completedAt",
			NewTask(taskID, version, userID, title, description, StatusDone, createdAt, &completedAtInvalid),
			fmt.Errorf("invalid complition time: invalid argument"),
		},
	}

	for _, input := range inputs {
		t.Run(input.name, func(t *testing.T) {
			task := input.input
			err := task.Validate()
			if input.expErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, input.expErr.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestApplyPatch(t *testing.T) {
	const (
		taskID      = 1
		userID      = 1
		version     = 1
		title       = "Clean the park"
		description = "Remove plastic bottles near the lake"
	)

	createdAt := time.Date(2026, time.September, 28, 15, 0, 0, 0, time.UTC)

	validTitle := "New Valid Title"
	validDescription := "New Valid Description"
	invalidShortTitle := "Hi"
	invalidShortDescription := "plastic"

	initialTaskCreated := NewTask(taskID, version, userID, title, description, StatusCreated, createdAt, nil)
	initialTaskInProgress := NewTask(taskID, version, userID, title, description, StatusInProgress, createdAt, nil)

	inputs := []testTaskInput{
		{
			"success: patch title",
			initialTaskCreated,
			nil,
		},
		{
			"success: patch description",
			initialTaskCreated,
			nil,
		},
		{
			"success: patch title, description",
			initialTaskCreated,
			nil,
		},
		{
			"fail: status != created",
			initialTaskInProgress,
			fmt.Errorf("you can modify task only in status: created, current status: in_progress: invalid argument"),
		},
		{
			"fail: title fails validation",
			initialTaskCreated,
			fmt.Errorf("validate task: invalid 'description' length: 2: invalid argument"),
		},
		{
			"fail: description fails validation",
			initialTaskCreated,
			fmt.Errorf("validate task: invalid 'description' length: 7: invalid argument"),
		},
	}

	for _, input := range inputs {
		t.Run(input.name, func(t *testing.T) {
			task := input.input

			var titlePtr *string
			var descPtr *string

			switch input.name {
			case "success: patch title":
				titlePtr = &validTitle
			case "success: patch description":
				descPtr = &validDescription
			case "success: patch title, description":
				titlePtr = &validTitle
				descPtr = &validDescription
			case "fail: status != created":
				titlePtr = &validTitle
				descPtr = &validDescription
			case "fail: title fails validation":
				descPtr = &invalidShortTitle
			case "fail: description fails validation":
				descPtr = &invalidShortDescription
			}

			err := task.ApplyPatch(titlePtr, descPtr)
			if input.expErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, input.expErr.Error())
				assert.Equal(t, input.input.Title, task.Title)
				assert.Equal(t, input.input.Description, task.Description)
				return
			}
			require.NoError(t, err)
			if titlePtr != nil {
				assert.Equal(t, *titlePtr, task.Title)
			}
			if descPtr != nil {
				assert.Equal(t, *descPtr, task.Description)
			}
		})
	}
}

func TestApplyStatusPatch(t *testing.T) {
	const (
		taskID                  = 1
		userID                  = 1
		invalidUserID           = -1
		version                 = 1
		title                   = "Clean the park"
		invalidShortTitle       = "Hi"
		invalidLongTitle        = "1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111"
		description             = "Remove plastic bottles near the lake"
		invalidLongDescription  = "11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111"
		invalidShortDescription = "plastic"
	)

	createdAt := time.Date(2026, time.September, 28, 15, 0, 0, 0, time.UTC)
	invalidCreatedAt := time.Date(2026, time.October, 28, 15, 0, 0, 0, time.UTC)

	inputs := []testTaskInput{
		{
			"success: transition from created to in_progress",
			NewTask(taskID, version, userID, title, description, StatusCreated, createdAt, nil),
			nil,
		},
		{
			"success: transition from created to rejected",
			NewTask(taskID, version, userID, title, description, StatusCreated, createdAt, nil),
			nil,
		},
		{
			"success: transition from in_progress to done",
			NewTask(taskID, version, userID, title, description, StatusInProgress, createdAt, nil),
			nil,
		},
		{
			"success: transition from in_progress to rejected",
			NewTask(taskID, version, userID, title, description, StatusInProgress, createdAt, nil),
			nil,
		},
		{
			"fail: change to the same status",
			NewTask(taskID, version, userID, title, description, StatusCreated, createdAt, nil),
			fmt.Errorf("status of the task is already: created: invalid argument"),
		},
		{
			"fail: transition from created to done",
			NewTask(taskID, version, userID, title, description, StatusCreated, createdAt, nil),
			fmt.Errorf("transition from: created to: done is not allowed: invalid argument"),
		},
		{
			"fail: validation error",
			NewTask(taskID, version, userID, title, description, StatusCreated, invalidCreatedAt, nil),
			fmt.Errorf("validate task: invalid complition time: invalid argument"),
		},
	}

	for _, input := range inputs {
		t.Run(input.name, func(t *testing.T) {
			task := input.input

			var targetStatus status
			switch input.name {
			case "success: transition from created to in_progress":
				targetStatus = StatusInProgress
			case "success: transition from created to rejected":
				targetStatus = StatusRejected
			case "success: transition from in_progress to done":
				targetStatus = StatusDone
			case "success: transition from in_progress to rejected":
				targetStatus = StatusRejected
			case "fail: change to the same status":
				targetStatus = StatusCreated
			case "fail: transition from created to done":
				targetStatus = StatusDone
			case "fail: validation error":
				targetStatus = StatusRejected
			}

			startTime := time.Now()
			err := task.ApplyStatusPatch(targetStatus)

			if input.expErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, input.expErr.Error())
				assert.Equal(t, input.input.Status, task.Status)
				assert.Equal(t, input.input.CompletedAt, task.CompletedAt)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, targetStatus, task.Status)

			if targetStatus == StatusDone || targetStatus == StatusRejected {
				require.NotNil(t, task.CompletedAt)
				assert.WithinDuration(t, startTime, *task.CompletedAt, time.Second)
			} else {
				assert.Nil(t, task.CompletedAt)
			}
		})
	}
}
