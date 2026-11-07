package statistic

import C "github.com/metacubex/mihomo/constant"

type RequestNotify func(c Tracker)

var DefaultRequestNotify RequestNotify

// TrafficAccount receives the bytes a tracker moves. The Core installs the
// factory once; it runs after the tracker's metadata and rule are known and
// before the manager publishes the tracker, so the returned closure may keep
// only its category and accumulator.
type TrafficAccount func(upload, download int64)

var NewTrafficAccount func(metadata *C.Metadata, rule C.Rule, chain C.Chain) TrafficAccount

// IsDirect reports whether the named outbound reaches the destination without
// a proxy. A direct proxy may carry any name, so the tunnel, which knows each
// proxy's type, replaces this name check.
var IsDirect = func(name string) bool { return name == "DIRECT" }

func (m *Manager) TotalTraffic(onlyProxy bool) (up, down int64) {
	if onlyProxy {
		return m.proxyUploadTotal.Load(), m.proxyDownloadTotal.Load()
	}
	return m.uploadTotal.Load(), m.downloadTotal.Load()
}

func (m *Manager) NowTraffic(onlyProxy bool) (up, down int64) {
	if onlyProxy {
		return m.proxyUploadBlip.Load(), m.proxyDownloadBlip.Load()
	}
	return m.uploadBlip.Load(), m.downloadBlip.Load()
}
