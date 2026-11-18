package adapter

import "github.com/metacubex/mihomo/adapter/outbound"

// CustomProxyFactory is installed once by the embedding application before parsing profiles.
var CustomProxyFactory func(map[string]any, outbound.BasicOption) (outbound.ProxyAdapter, error)

// timedOut follows hub/route's getProxyDelay, which answers 504 rather than 503
// once the test's context is done.
type UrlTestCheck func(url string, name string, delay uint16, timedOut bool)

var UrlTestHook UrlTestCheck
