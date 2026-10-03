package core_domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testTaskCommandInput struct {
	name        string
	title       *string
	description *string
}

func TestNewTaskCommand(t *testing.T) {
	title := "Clean the park"
	description := "Remove plastic bottles near the lake"

	inputs := []testTaskCommandInput{
		{
			"success: nil, nil",
			nil,
			nil,
		},
		{
			"success: title, nil",
			&title,
			nil,
		},
		{
			"success: nil, description",
			nil,
			&description,
		},
		{
			"success: title, description",
			&title,
			&description,
		},
	}

	for _, input := range inputs {
		t.Run(input.name, func(t *testing.T) {
			command := NewPatchTaskCommand(input.title, input.description)

			require.Equal(t, input.title, command.Title)
			require.Equal(t, input.description, command.Description)
		})
	}
}
