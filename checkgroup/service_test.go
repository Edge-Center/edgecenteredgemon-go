package checkgroup

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Edge-Center/edgecenteredgemon-go/internal/testutil"
)

func TestService_Create(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(Response{ID: 10, Name: "group1"})

	svc := New(mock)
	resp, err := svc.Create(context.Background(), &Request{Name: "group1"})

	require.NoError(t, err)
	assert.Equal(t, 10, resp.ID)
	assert.Equal(t, "group1", resp.Name)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodPost, mock.Calls[0].Method)
	assert.Equal(t, "/rmon/check-group", mock.Calls[0].Path)
}

func TestService_Get(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(Response{ID: 10, Name: "group1"})

	svc := New(mock)
	resp, err := svc.Get(context.Background(), 10)

	require.NoError(t, err)
	assert.Equal(t, 10, resp.ID)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodGet, mock.Calls[0].Method)
	assert.Equal(t, "/rmon/check-group/10", mock.Calls[0].Path)
}

func TestService_Update(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(Response{ID: 10, Name: "server-confirmed"})

	svc := New(mock)
	resp, err := svc.Update(context.Background(), 10, &Request{Name: "new-name"})

	require.NoError(t, err)
	assert.Equal(t, "server-confirmed", resp.Name)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodPut, mock.Calls[0].Method)
	assert.Equal(t, "/rmon/check-group/10", mock.Calls[0].Path)
}

func TestService_Delete(t *testing.T) {
	mock := &testutil.MockRequester{}
	svc := New(mock)

	err := svc.Delete(context.Background(), 10)

	require.NoError(t, err)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodDelete, mock.Calls[0].Method)
	assert.Equal(t, "/rmon/check-group/10", mock.Calls[0].Path)
}

func TestService_ErrorPropagation(t *testing.T) {
	errRemote := errors.New("timeout")

	tests := []struct {
		name string
		call func(svc Service) error
	}{
		{"Create", func(svc Service) error { _, err := svc.Create(context.Background(), &Request{}); return err }},
		{"Get", func(svc Service) error { _, err := svc.Get(context.Background(), 1); return err }},
		{"Update", func(svc Service) error { _, err := svc.Update(context.Background(), 1, &Request{}); return err }},
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
