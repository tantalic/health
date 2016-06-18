package health

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

type StatusCheckFunc func() bool

type Checker interface {
	IsLive() (live bool, err error)
	IsReady() (ready bool, err error)
}

type Service interface {
	Checker
	http.Handler
}

type healthHandler struct {
	mux     *http.ServeMux
	checker Checker
}

func NewHealthHandler(c Checker) http.Handler {
	h := &healthHandler{
		checker: c,
		mux:     http.NewServeMux(),
	}

	h.mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		live, err := h.checker.IsLive()

		if !live {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
		} else {
			w.Write([]byte("LIVE"))
		}
	})

	h.mux.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
		ready, err := h.checker.IsReady()

		if !ready {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
		} else {
			w.Write([]byte("READY"))
		}
	})

	h.mux.Handle("/metrics", prometheus.Handler())

	return h
}

func (h *healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}
