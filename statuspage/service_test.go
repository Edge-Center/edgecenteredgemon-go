package statuspage

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
	mock.RespondWith(CreateResponse{ID: 42})

	svc := New(mock)
	resp, err := svc.Create(context.Background(), &Request{
		Base:   Base{Name: "status", Slug: "s1"},
		Checks: []int{1, 2, 3},
	})

	require.NoError(t, err)
	assert.Equal(t, 42, resp.ID)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodPost, mock.Calls[0].Method)
	assert.Equal(t, "/rmon/status-page", mock.Calls[0].Path)
}

func TestService_Get(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(Response{
		ID:     42,
		Base:   Base{Name: "status", Slug: "s1"},
		Checks: []Checks{{CheckID: 1}, {CheckID: 2}},
	})

	svc := New(mock)
	resp, err := svc.Get(context.Background(), 42)

	require.NoError(t, err)
	assert.Equal(t, 42, resp.ID)
	assert.Equal(t, "status", resp.Name)
	assert.Len(t, resp.Checks, 2)
	assert.Equal(t, 1, resp.Checks[0].CheckID)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodGet, mock.Calls[0].Method)
	assert.Equal(t, "/rmon/status-page/42", mock.Calls[0].Path)
}

func TestService_Update(t *testing.T) {
	mock := &testutil.MockRequester{}
	svc := New(mock)

	err := svc.Update(context.Background(), 42, &Request{
		Base:   Base{Name: "updated"},
		Checks: []int{4},
	})

	require.NoError(t, err)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodPut, mock.Calls[0].Method)
	assert.Equal(t, "/rmon/status-page/42", mock.Calls[0].Path)
}

func TestService_Delete(t *testing.T) {
	mock := &testutil.MockRequester{}
	svc := New(mock)

	err := svc.Delete(context.Background(), 42)

	require.NoError(t, err)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, http.MethodDelete, mock.Calls[0].Method)
	assert.Equal(t, "/rmon/status-page/42", mock.Calls[0].Path)
}

func TestService_CreateReturnsCreateResponse(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(CreateResponse{ID: 99})

	svc := New(mock)
	resp, err := svc.Create(context.Background(), &Request{
		Base:   Base{Name: "page"},
		Checks: []int{1},
	})

	require.NoError(t, err)
	assert.IsType(t, &CreateResponse{}, resp)
	assert.Equal(t, 99, resp.ID)
}

func TestService_GetReturnsFullResponse(t *testing.T) {
	mock := &testutil.MockRequester{}
	mock.RespondWith(Response{
		ID:     99,
		Base:   Base{Name: "full", Slug: "full-slug", Description: "desc"},
		Checks: []Checks{{CheckID: 10}},
	})

	svc := New(mock)
	resp, err := svc.Get(context.Background(), 99)

	require.NoError(t, err)
	assert.IsType(t, &Response{}, resp)
	assert.Equal(t, "full", resp.Name)
	assert.Equal(t, "full-slug", resp.Slug)
	assert.Len(t, resp.Checks, 1)
}

func TestService_ErrorPropagation(t *testing.T) {
	errRemote := errors.New("forbidden")

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
