package checks

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Edge-Center/edgecenteredgemon-go/internal/testutil"
)

type testRequest struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}

type testResponse struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func TestService_Create_UsesCheckTypeAndPayload(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(CreateResponse{ID: 42})

	svc := NewService[testRequest, testResponse](mock, "custom")
	resp, err := svc.Create(context.Background(), &testRequest{
		Name: "test-check",
		Port: 443,
	})

	require.NoError(t, err)
	assert.Equal(t, 42, resp.ID)
	require.Len(t, mock.Calls, 1)

	call := mock.Calls[0]
	assert.Equal(t, http.MethodPost, call.Method)
	assert.Equal(t, "/rmon/check/custom", call.Path)

	req, ok := call.Payload.(*testRequest)
	require.True(t, ok)
	assert.Equal(t, "test-check", req.Name)
	assert.Equal(t, 443, req.Port)
}

func TestService_Get_DecodesGenericResponse(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(testResponse{
		Name:   "decoded-check",
		Status: "ok",
	})

	svc := NewService[testRequest, testResponse](mock, "custom")
	resp, err := svc.Get(context.Background(), 5)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "decoded-check", resp.Name)
	assert.Equal(t, "ok", resp.Status)
	require.Len(t, mock.Calls, 1)

	call := mock.Calls[0]
	assert.Equal(t, http.MethodGet, call.Method)
	assert.Equal(t, "/rmon/check/custom/5", call.Path)
	assert.Nil(t, call.Payload)
}

func TestService_Methods_BuildExpectedPaths(t *testing.T) {
	tests := []struct {
		name       string
		call       func(Service[testRequest, testResponse]) error
		wantMethod string
		wantPath   string
		wantNilReq bool
	}{
		{
			name: "Update",
			call: func(svc Service[testRequest, testResponse]) error {
				return svc.Update(context.Background(), 7, &testRequest{Name: "updated", Port: 80})
			},
			wantMethod: http.MethodPut,
			wantPath:   "/rmon/check/custom/7",
			wantNilReq: false,
		},
		{
			name: "Delete",
			call: func(svc Service[testRequest, testResponse]) error {
				return svc.Delete(context.Background(), 7)
			},
			wantMethod: http.MethodDelete,
			wantPath:   "/rmon/check/custom/7",
			wantNilReq: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &testutil.MockRequester{}
			svc := NewService[testRequest, testResponse](mock, "custom")

			err := tt.call(svc)

			require.NoError(t, err)
			require.Len(t, mock.Calls, 1)

			call := mock.Calls[0]
			assert.Equal(t, tt.wantMethod, call.Method)
			assert.Equal(t, tt.wantPath, call.Path)

			if tt.wantNilReq {
				assert.Nil(t, call.Payload)
			} else {
				assert.NotNil(t, call.Payload)
			}
		})
	}
}

func TestService_ErrorPropagation(t *testing.T) {
	errRemote := errors.New("connection refused")

	tests := []struct {
		name        string
		call        func(Service[testRequest, testResponse]) error
		wantMessage string
	}{
		{
			name: "Create",
			call: func(svc Service[testRequest, testResponse]) error {
				_, err := svc.Create(context.Background(), &testRequest{})
				return err
			},
			wantMessage: "request: connection refused",
		},
		{
			name: "Get",
			call: func(svc Service[testRequest, testResponse]) error {
				_, err := svc.Get(context.Background(), 1)
				return err
			},
			wantMessage: "request: connection refused",
		},
		{
			name: "Update",
			call: func(svc Service[testRequest, testResponse]) error {
				return svc.Update(context.Background(), 1, &testRequest{})
			},
			wantMessage: "request: connection refused",
		},
		{
			name: "Delete",
			call: func(svc Service[testRequest, testResponse]) error {
				return svc.Delete(context.Background(), 1)
			},
			wantMessage: "request: connection refused",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &testutil.MockRequester{}
			mock.RespondWithError(errRemote)
			svc := NewService[testRequest, testResponse](mock, "custom")

			err := tt.call(svc)

			require.Error(t, err)
			assert.ErrorIs(t, err, errRemote)
			assert.EqualError(t, err, tt.wantMessage)
		})
	}
}
