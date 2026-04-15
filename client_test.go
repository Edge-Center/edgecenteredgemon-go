package edgecenterrmon

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Edge-Center/edgecenteredgemon-go/internal/testutil"
)

var _ ClientService = (*Service)(nil)

func TestNewService_NotNil(t *testing.T) {
	mock := &testutil.MockRequester{}
	svc := NewService(mock)
	require.NotNil(t, svc)
}

func TestService_Accessors(t *testing.T) {
	mock := &testutil.MockRequester{}
	svc := NewService(mock)

	tests := []struct {
		name     string
		accessor func() interface{}
	}{
		{"Channel", func() interface{} { return svc.Channel() }},
		{"StatusPage", func() interface{} { return svc.StatusPage() }},
		{"CheckGroup", func() interface{} { return svc.CheckGroup() }},
		{"CheckDNS", func() interface{} { return svc.CheckDNS() }},
		{"CheckHTTP", func() interface{} { return svc.CheckHTTP() }},
		{"CheckPing", func() interface{} { return svc.CheckPing() }},
		{"CheckRabbitMQ", func() interface{} { return svc.CheckRabbitMQ() }},
		{"CheckSMTP", func() interface{} { return svc.CheckSMTP() }},
		{"CheckTCP", func() interface{} { return svc.CheckTCP() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.accessor())
		})
	}
}
