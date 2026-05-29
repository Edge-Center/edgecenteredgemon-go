package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	edgecenterrmon "github.com/Edge-Center/edgecenteredgemon-go"
	"github.com/Edge-Center/edgecenteredgemon-go/channel"
	"github.com/Edge-Center/edgecenteredgemon-go/checkgroup"
	"github.com/Edge-Center/edgecenteredgemon-go/checks"
	"github.com/Edge-Center/edgecenteredgemon-go/checks/checkdns"
	"github.com/Edge-Center/edgecenteredgemon-go/checks/checkhttp"
	"github.com/Edge-Center/edgecenteredgemon-go/checks/checkping"
	"github.com/Edge-Center/edgecenteredgemon-go/checks/checkrabbitmq"
	"github.com/Edge-Center/edgecenteredgemon-go/checks/checksmtp"
	"github.com/Edge-Center/edgecenteredgemon-go/checks/checktcp"
	"github.com/Edge-Center/edgecenteredgemon-go/edgecenter"
	"github.com/Edge-Center/edgecenteredgemon-go/statuspage"
)

func TestIntegration_AuthAndUserAgentReachServer(t *testing.T) {
	client, srv := start(t)
	svc := edgecenterrmon.NewService(client)

	_, err := svc.CheckGroup().Create(context.Background(), &checkgroup.Request{Name: "auth-check"})
	require.NoError(t, err)

	calls := srv.calls()
	require.Len(t, calls, 1)
	assert.Equal(t, "APIKey "+apiKey, calls[0].Auth)
	assert.NotEmpty(t, calls[0].Body)
}

func TestIntegration_Channel_Lifecycle(t *testing.T) {
	client, srv := start(t)
	svc := edgecenterrmon.NewService(client).Channel()
	ctx := context.Background()

	created, err := svc.Create(ctx, "email", &channel.Request{Channel: "ops-email", Token: "tok-123"})
	require.NoError(t, err)
	require.NotNil(t, created)
	assert.NotZero(t, created.ID)
	assert.Equal(t, "ops-email", created.Channel)

	got, err := svc.Get(ctx, "email", created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, got.ID)
	assert.Equal(t, "ops-email", got.Channel)

	require.NoError(t, svc.Update(ctx, "email", created.ID, &channel.Request{Channel: "ops-email-2", Token: "tok-456"}))
	require.NoError(t, svc.Delete(ctx, "email", created.ID))

	assertPaths(t, srv.calls(), []call{
		{"POST", "/rmon/channel/email"},
		{"GET", "/rmon/channel/email/" + itoa(created.ID)},
		{"PUT", "/rmon/channel/email/" + itoa(created.ID)},
		{"DELETE", "/rmon/channel/email/" + itoa(created.ID)},
	})
}

func TestIntegration_CheckGroup_Lifecycle(t *testing.T) {
	client, srv := start(t)
	svc := edgecenterrmon.NewService(client).CheckGroup()
	ctx := context.Background()

	created, err := svc.Create(ctx, &checkgroup.Request{Name: "prod"})
	require.NoError(t, err)
	assert.NotZero(t, created.ID)
	assert.Equal(t, "prod", created.Name)

	got, err := svc.Get(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, "prod", got.Name)

	updated, err := svc.Update(ctx, created.ID, &checkgroup.Request{Name: "prod-eu"})
	require.NoError(t, err)
	assert.Equal(t, "prod-eu", updated.Name)

	require.NoError(t, svc.Delete(ctx, created.ID))

	assertPaths(t, srv.calls(), []call{
		{"POST", "/rmon/check-group"},
		{"GET", "/rmon/check-group/" + itoa(created.ID)},
		{"PUT", "/rmon/check-group/" + itoa(created.ID)},
		{"DELETE", "/rmon/check-group/" + itoa(created.ID)},
	})
}

func TestIntegration_StatusPage_Lifecycle(t *testing.T) {
	client, srv := start(t)
	svc := edgecenterrmon.NewService(client).StatusPage()
	ctx := context.Background()

	req := &statuspage.Request{
		Base:   statuspage.Base{Name: "Public", Slug: "public", Description: "uptime"},
		Checks: []int{1, 2, 3},
	}

	created, err := svc.Create(ctx, req)
	require.NoError(t, err)
	assert.NotZero(t, created.ID)

	got, err := svc.Get(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, "Public", got.Name)
	assert.Equal(t, "public", got.Slug)
	require.Len(t, got.Checks, 3)
	assert.Equal(t, 1, got.Checks[0].CheckID)

	require.NoError(t, svc.Update(ctx, created.ID, req))
	require.NoError(t, svc.Delete(ctx, created.ID))

	assertPaths(t, srv.calls(), []call{
		{"POST", "/rmon/status-page"},
		{"GET", "/rmon/status-page/" + itoa(created.ID)},
		{"PUT", "/rmon/status-page/" + itoa(created.ID)},
		{"DELETE", "/rmon/status-page/" + itoa(created.ID)},
	})
}

func TestIntegration_Checks_Lifecycle(t *testing.T) {
	t.Run("dns", func(t *testing.T) {
		runCheckLifecycle(t, checkdns.New, "dns",
			&checkdns.Request{Name: "dns-check", IP: "1.1.1.1", Resolver: "8.8.8.8", RecordType: "A", Place: "all"})
	})
	t.Run("http", func(t *testing.T) {
		runCheckLifecycle(t, checkhttp.New, "http",
			&checkhttp.Request{Name: "http-check", URL: "https://example.com", Method: "GET", Place: "all", AcceptedStatusCodes: []int{200}})
	})
	t.Run("ping", func(t *testing.T) {
		runCheckLifecycle(t, checkping.New, "ping",
			&checkping.Request{Name: "ping-check", IP: "1.1.1.1", Place: "all"})
	})
	t.Run("rabbitmq", func(t *testing.T) {
		runCheckLifecycle(t, checkrabbitmq.New, "rabbitmq",
			&checkrabbitmq.Request{Name: "rmq-check", IP: "10.0.0.1", Port: 5672, Username: "guest", Password: "guest", Vhost: "/", Place: "all"})
	})
	t.Run("smtp", func(t *testing.T) {
		runCheckLifecycle(t, checksmtp.New, "smtp",
			&checksmtp.Request{Name: "smtp-check", IP: "10.0.0.2", Port: 587, Username: "u", Password: "p", Place: "all"})
	})
	t.Run("tcp", func(t *testing.T) {
		runCheckLifecycle(t, checktcp.New, "tcp",
			&checktcp.Request{Name: "tcp-check", IP: "10.0.0.3", Port: 443, Priority: "high", Place: "all"})
	})
}

func runCheckLifecycle[Req any, Resp any](t *testing.T, newSvc func(edgecenter.Requester) checks.Service[Req, Resp], checkType string, req *Req) {
	t.Helper()

	client, srv := start(t)
	svc := newSvc(client)
	ctx := context.Background()

	created, err := svc.Create(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, created)
	assert.NotZero(t, created.ID)

	got, err := svc.Get(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, got)

	require.NoError(t, svc.Update(ctx, created.ID, req))
	require.NoError(t, svc.Delete(ctx, created.ID))

	base := "/rmon/check/" + checkType
	assertPaths(t, srv.calls(), []call{
		{"POST", base},
		{"GET", base + "/" + itoa(created.ID)},
		{"PUT", base + "/" + itoa(created.ID)},
		{"DELETE", base + "/" + itoa(created.ID)},
	})
}

func TestIntegration_ErrorMapping_NotFound(t *testing.T) {
	client, _ := start(t)
	svc := edgecenterrmon.NewService(client)
	ctx := context.Background()

	_, err := svc.CheckGroup().Get(ctx, 999)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "404")
}
