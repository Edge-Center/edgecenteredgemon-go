package checkrabbitmq

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Edge-Center/edgecenteredgemon-go/checks"
	"github.com/Edge-Center/edgecenteredgemon-go/internal/testutil"
)

func TestService_Create(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(checks.CreateResponse{ID: 1})

	svc := New(mock)
	resp, err := svc.Create(context.Background(), &Request{Name: "rabbitmq-check"})

	require.NoError(t, err)
	assert.Equal(t, 1, resp.ID)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodPost, mock.Calls[0].Method)
	assert.Equal(t, "/rmon/check/rabbitmq", mock.Calls[0].Path)
}

func TestService_Get(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(Response{Name: "rabbitmq-check"})

	svc := New(mock)
	resp, err := svc.Get(context.Background(), 5)

	require.NoError(t, err)
	assert.Equal(t, "rabbitmq-check", resp.Name)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodGet, mock.Calls[0].Method)
	assert.Equal(t, "/rmon/check/rabbitmq/5", mock.Calls[0].Path)
}

func TestService_Update(t *testing.T) {
	mock := &testutil.MockRequester{}
	svc := New(mock)

	err := svc.Update(context.Background(), 5, &Request{Name: "updated"})

	require.NoError(t, err)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodPut, mock.Calls[0].Method)
	assert.Equal(t, "/rmon/check/rabbitmq/5", mock.Calls[0].Path)
}

func TestService_Delete(t *testing.T) {
	mock := &testutil.MockRequester{}
	svc := New(mock)

	err := svc.Delete(context.Background(), 5)

	require.NoError(t, err)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodDelete, mock.Calls[0].Method)
	assert.Equal(t, "/rmon/check/rabbitmq/5", mock.Calls[0].Path)
}

func TestService_ErrorPropagation(t *testing.T) {
	errRemote := errors.New("connection refused")

	tests := []struct {
		name string
		call func(svc Service) error
	}{
		{"Create", func(svc Service) error { _, err := svc.Create(context.Background(), &Request{}); return err }},
		{"Get", func(svc Service) error { _, err := svc.Get(context.Background(), 1); return err }},
		{"Update", func(svc Service) error { return svc.Update(context.Background(), 1, &Request{}) }},
		{"Delete", func(svc Service) error { return svc.Delete(context.Background(), 1) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &testutil.MockRequester{}
			mock.RespondWithError(errRemote)
			svc := New(mock)

			err := tt.call(svc)
			require.Error(t, err)
			assert.ErrorIs(t, err, errRemote)
		})
	}
}
