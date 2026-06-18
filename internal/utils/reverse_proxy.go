package utils

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func ProxyToService(targetBaseUrl string) *httputil.ReverseProxy {
	// parse the upstream service
	target, _ := url.Parse(targetBaseUrl)

	// create a reverse proxy instance
	proxy := httputil.NewSingleHostReverseProxy(target)

	// create a director to manipulate requests
	orgDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		orgDirector(req)

		req.Host = target.Host
		req.URL.Host = target.Host
		req.URL.Scheme = target.Scheme
	}

	return proxy
}
