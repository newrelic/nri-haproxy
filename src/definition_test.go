package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMillisToSecondsAmendment(t *testing.T) {
	t.Parallel()

	md := metricDefinition{
		Amend: millisToSeconds,
	}

	testCases := []struct {
		name     string
		value    any
		expected any
	}{
		{
			name:     "As string",
			value:    "1000",
			expected: 1.0,
		},
		{
			name:     "As int",
			value:    1,
			expected: 0.001,
		},
		{
			name:     "As float",
			value:    1.0,
			expected: 0.001,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			amended, err := md.value(tc.value)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, amended)
		})
	}

	errCases := []struct {
		name  string
		value any
	}{
		{
			name:  "Invalid type",
			value: "invalid",
		},
	}

	for _, tc := range errCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := md.value(tc.value)
			assert.Error(t, err)
		})
	}
}
