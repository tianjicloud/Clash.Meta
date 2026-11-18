package outboundgroup

import (
	"net/http"
	"net/http/httptest"
	"testing"

	C "github.com/metacubex/mihomo/constant"
	P "github.com/metacubex/mihomo/constant/provider"
	"github.com/stretchr/testify/require"
)

type urlTestProxy struct {
	C.Proxy
	delay uint16
}

func (p *urlTestProxy) AliveForTestUrl(string) bool       { return true }
func (p *urlTestProxy) LastDelayForTestUrl(string) uint16 { return p.delay }

func proxyWithDelay(name string, delay uint16) C.Proxy {
	return &urlTestProxy{Proxy: namedProxy(name), delay: delay}
}

func newURLTestForDelays(t *testing.T, tolerance uint16, proxies ...C.Proxy) *URLTest {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(server.Close)
	group, err := NewURLTest(
		GroupCommonOption{Name: "auto", URL: server.URL},
		URLTestOption{Tolerance: tolerance},
		namedProxy("COMPATIBLE"),
		[]P.ProxyProvider{&fakeProvider{name: "p", proxies: proxies}},
	)
	require.NoError(t, err)
	return group
}

func TestURLTestSwitchesWhenFastNodeIsNotFirstMember(t *testing.T) {
	slower := proxyWithDelay("slower", 80)
	fastest := proxyWithDelay("fastest", 20)
	group := newURLTestForDelays(t, 10, slower, fastest)
	group.fastNode = slower

	require.Equal(t, "fastest", group.Now())
}

func TestURLTestToleranceDoesNotOverflowUint16(t *testing.T) {
	current := proxyWithDelay("current", 1000)
	faster := proxyWithDelay("faster", 900)
	group := newURLTestForDelays(t, 200, current, faster)
	group.fastNode = current

	require.Equal(t, "current", group.Now())
	group.fastSingle.Reset()
	group.tolerance = 50
	require.Equal(t, "faster", group.Now())
}
