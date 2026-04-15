package channel

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
	receivers := []string{"email", "slack", "telegram", "mattermost", "pagerduty"}

	for _, recv := range receivers {
		t.Run(recv, func(t *testing.T) {
			mock := &testutil.MockRequester{}
			mock.RespondWith(Response{ID: 1, Channel: "ch"})

			svc := New(mock)
			resp, err := svc.Create(context.Background(), recv, &Request{Channel: "ch", Token: "tok"})

			require.NoError(t, err)
			assert.Equal(t, 1, resp.ID)
			require.Len(t, mock.Calls, 1)
			assert.Equal(t, http.MethodPost, mock.Calls[0].Method)
			assert.Equal(t, "/rmon/channel/"+recv, mock.Calls[0].Path)
			assert.NotNil(t, mock.Calls[0].Payload)
		})
	}
}

func TestService_Get(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(Response{ID: 5, Channel: "test-ch"})

	svc := New(mock)
	resp, err := svc.Get(context.Background(), "email", 5)

	require.NoError(t, err)
	assert.Equal(t, 5, resp.ID)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodGet, mock.Calls[0].Method)
	assert.Equal(t, "/rmon/channel/email/5", mock.Calls[0].Path)
}

func TestService_Update(t *testing.T) {
	mock := &testutil.MockRequester{}
	svc := New(mock)

	err := svc.Update(context.Background(), "slack", 3, &Request{Channel: "upd"})

	require.NoError(t, err)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodPut, mock.Calls[0].Method)
	assert.Equal(t, "/rmon/channel/slack/3", mock.Calls[0].Path)
	assert.NotNil(t, mock.Calls[0].Payload)
}

func TestService_Delete(t *testing.T) {
	mock := &testutil.MockRequester{}
	svc := New(mock)

	err := svc.Delete(context.Background(), "telegram", 7)

	require.NoError(t, err)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodDelete, mock.Calls[0].Method)
	assert.Equal(t, "/rmon/channel/telegram/7", mock.Calls[0].Path)
}

func TestService_ErrorPropagation(t *testing.T) {
	errRemote := errors.New("connection refused")

	tests := []struct {
		name string
		call func(svc Service) error
	}{
		{"Create", func(svc Service) error { _, err := svc.Create(context.Background(), "email", &Request{}); return err }},
		{"Get", func(svc Service) error { _, err := svc.Get(context.Background(), "email", 1); return err }},
		{"Update", func(svc Service) error { return svc.Update(context.Background(), "email", 1, &Request{}) }},
		{"Delete", func(svc Service) error { return svc.Delete(context.Background(), "email", 1) }},
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
