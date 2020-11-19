package httpproxy

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

var (
	DefaultTransport = http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			Deadline:  time.Now().Add(time.Duration(31) * time.Second),
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
)

type HttpProxy struct {
	// proxy target url
	target *url.URL

	transport *http.Transport
}

func NewHttpProxy(target *url.URL, transport *http.Transport) *HttpProxy {
	httpProxy := HttpProxy{target: target}

	if transport == nil {
		httpProxy.transport = &DefaultTransport
	}

	return &httpProxy
}

// With go-gin, use ServeHTTP(c.Wirter, c.Request), c is gin.Context.
func (hp *HttpProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	proxy := httputil.NewSingleHostReverseProxy(hp.target)
	proxy.Transport = hp.transport

	proxy.ServeHTTP(rw, req)
}
