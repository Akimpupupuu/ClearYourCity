package core_domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testStatusInput struct {
	name     string
	input    string
	expValue status
	expErr   error
}

func TestNewStatus(t *testing.T) {
	inputs := []testStatusInput{
		{
			"success: created",
			"created",
			StatusCreated,
			nil,
		},
		{
			"success: in_progress",
			"in_progress",
			StatusInProgress,
			nil,
		},
		{
			"success: done",
			"done",
			StatusDone,
			nil,
		},
		{
			"success: rejected",
			"rejected",
			StatusRejected,
			nil,
		},
		{
			"fail: empty status",
			"",
			"",
			fmt.Errorf("invalid status: invalid argument"),
		},
		{
			"fail: invalid status",
			"invalid status",
			"",
			fmt.Errorf("invalid status: invalid argument"),
		},
	}

	for _, input := range inputs {
		t.Run(input.name, func(t *testing.T) {
			status, err := NewStatus(input.input)
			if input.expErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, input.expErr.Error())
				assert.Equal(t, input.expValue, status)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, input.expValue, status)
		})
	}
}

func TestNewStatusCreated(t *testing.T) {
	status := NewStatusCreated()

	require.Equal(t, StatusCreated, status)
}
