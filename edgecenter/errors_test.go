package edgecenter

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorResponse_Error(t *testing.T) {
	errorsFilled := json.RawMessage(`{"code":"validation_error"}`)
	tests := []struct {
		name     string
		resp     *ErrorResponse
		expected string
	}{
		{
			name:     "message set returns message",
			resp:     &ErrorResponse{Message: "something went wrong"},
			expected: "something went wrong",
		},
		{
			name:     "message empty errors set returns json",
			resp:     &ErrorResponse{Errors: &errorsFilled},
			expected: "{\n\t\"code\": \"validation_error\"\n}",
		},
		{
			name:     "message empty errors nil returns null",
			resp:     &ErrorResponse{},
			expected: "null",
		},
		{
			name:     "both set message has priority",
			resp:     &ErrorResponse{Message: "priority message", Errors: &errorsFilled},
			expected: "priority message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.resp.Error())
		})
	}
}
