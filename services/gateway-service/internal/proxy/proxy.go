package proxy

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func ReverseProxy(targetURL string) gin.HandlerFunc {
	target, err := url.Parse(targetURL)
	if err != nil {
		panic("invalid target URL: " + targetURL)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Transport = &http.Transport{
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"error":"service unavailable","code":"PROXY_ERROR"}`))
	}

	return func(c *gin.Context) {
		c.Request.URL.Host = target.Host
		c.Request.URL.Scheme = target.Scheme
		c.Request.Header.Set("X-Forwarded-Host", c.Request.Host)

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
