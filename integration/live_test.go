//go:build integration

package integration

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	edgecenterrmon "github.com/Edge-Center/edgecenteredgemon-go"
	"github.com/Edge-Center/edgecenteredgemon-go/channel"
	"github.com/Edge-Center/edgecenteredgemon-go/checkgroup"
	"github.com/Edge-Center/edgecenteredgemon-go/checks/checkping"
	"github.com/Edge-Center/edgecenteredgemon-go/edgecenter/provider"
	"github.com/Edge-Center/edgecenteredgemon-go/statuspage"
)

func liveClient(t *testing.T) edgecenterrmon.ClientService {
	t.Helper()

	baseURL := os.Getenv("EDGECENTER_RMON_BASE_URL")
	key := os.Getenv("EDGECENTER_RMON_API_KEY")

	if baseURL == "" || key == "" {
		t.Skip("set EDGECENTER_RMON_BASE_URL and EDGECENTER_RMON_API_KEY to run live integration tests")
	}

	signer := provider.WithSignerFunc(func(req *http.Request) error {
		for k, v := range provider.AuthenticatedHeaders(key) {
			req.Header.Set(k, v)
		}
		return nil
	})

	client := provider.NewClient(baseURL, signer,
		provider.WithUserAgent("edgecenteredgemon-go/integration"),
		provider.WithTimeout(30*time.Second),
	)

	return edgecenterrmon.NewService(client)
}

func liveCtx(t *testing.T) context.Context {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	t.Cleanup(cancel)

	return ctx
}

func requireEnv(t *testing.T, names ...string) []string {
	t.Helper()

	values := make([]string, len(names))
	for i, n := range names {
		v := os.Getenv(n)
		if v == "" {
			t.Skipf("set %s to run this live test", n)
		}
		values[i] = v
	}

	return values
}

func uniqueName(prefix string) string {
	return prefix + "-" + time.Now().UTC().Format("20060102-150405")
}

func TestLive_CheckGroup_Lifecycle(t *testing.T) {
	svc := liveClient(t).CheckGroup()
	ctx := liveCtx(t)

	name := uniqueName("sdk-it-grp")

	created, err := svc.Create(ctx, &checkgroup.Request{Name: name})
	require.NoError(t, err)
	require.NotZero(t, created.ID)

	t.Cleanup(func() {
		if err := svc.Delete(context.Background(), created.ID); err != nil {
			t.Logf("cleanup: failed to delete check group %d: %v", created.ID, err)
		}
	})

	got, err := svc.Get(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, name, got.Name)

	updated, err := svc.Update(ctx, created.ID, &checkgroup.Request{Name: name + "-upd"})
	require.NoError(t, err)
	assert.Equal(t, name+"-upd", updated.Name)
}

func TestLive_StatusPage_Lifecycle(t *testing.T) {
	svc := liveClient(t).StatusPage()
	ctx := liveCtx(t)

	slug := uniqueName("sdk-it-page")
	req := &statuspage.Request{
		Base:   statuspage.Base{Name: "SDK IT", Slug: slug, Description: "created by integration test"},
		Checks: []int{},
	}

	created, err := svc.Create(ctx, req)
	require.NoError(t, err)
	require.NotZero(t, created.ID)

	t.Cleanup(func() {
		if err := svc.Delete(context.Background(), created.ID); err != nil {
			t.Logf("cleanup: failed to delete status page %d: %v", created.ID, err)
		}
	})

	got, err := svc.Get(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, slug, got.Slug)

	req.Description = "updated by integration test"
	require.NoError(t, svc.Update(ctx, created.ID, req))
}

func TestLive_Channel_Lifecycle(t *testing.T) {
	env := requireEnv(t, "EDGECENTER_RMON_TEST_CHANNEL_RECEIVER", "EDGECENTER_RMON_TEST_CHANNEL_TOKEN")
	receiver, token := env[0], env[1]

	svc := liveClient(t).Channel()
	ctx := liveCtx(t)

	created, err := svc.Create(ctx, receiver, &channel.Request{Channel: uniqueName("sdk-it-ch"), Token: token})
	require.NoError(t, err)
	require.NotZero(t, created.ID)

	t.Cleanup(func() {
		if err := svc.Delete(context.Background(), receiver, created.ID); err != nil {
			t.Logf("cleanup: failed to delete channel %d: %v", created.ID, err)
		}
	})

	got, err := svc.Get(ctx, receiver, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, got.ID)

	require.NoError(t, svc.Update(ctx, receiver, created.ID, &channel.Request{Channel: uniqueName("sdk-it-ch2"), Token: token}))
}

func TestLive_CheckPing_Lifecycle(t *testing.T) {
	env := requireEnv(t, "EDGECENTER_RMON_TEST_PLACE", "EDGECENTER_RMON_TEST_ENTITY")
	place := env[0]

	entity, err := strconv.Atoi(env[1])
	require.NoError(t, err, "EDGECENTER_RMON_TEST_ENTITY must be an integer")

	svc := liveClient(t).CheckPing()
	ctx := liveCtx(t)

	req := &checkping.Request{
		Name:     uniqueName("sdk-it-ping"),
		Enabled:  0,
		Place:    place,
		Entities: []int{entity},
		IP:       "1.1.1.1",
	}

	created, err := svc.Create(ctx, req)
	require.NoError(t, err)
	require.NotZero(t, created.ID)

	t.Cleanup(func() {
		if err := svc.Delete(context.Background(), created.ID); err != nil {
			t.Logf("cleanup: failed to delete ping check %d: %v", created.ID, err)
		}
	})

	got, err := svc.Get(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, "1.1.1.1", got.IP)

	req.IP = "8.8.8.8"
	require.NoError(t, svc.Update(ctx, created.ID, req))
}
