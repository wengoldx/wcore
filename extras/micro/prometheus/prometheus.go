package prometheus

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/wengoldx/wing/logger"
	"net/http"
)

// NewProme service monitoring tools
func NewProme(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			logger.E("ListenAndServe: addr is :", addr, "err :", err)
		}
	}()
}
