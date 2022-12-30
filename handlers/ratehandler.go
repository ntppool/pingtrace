package handlers

import (
	"log"
	"net/http"
)

// RateHandler wraps a handler in rate limiting
type RateHandler struct {
	limiter Limiter
	h       http.Handler
}

// NewRateHandler returns h wrapped in the Limiter
func NewRateHandler(h http.Handler, limiter Limiter) *RateHandler {
	return &RateHandler{limiter: limiter, h: h}
}

// ServeHTTP implements http.Server interface
func (rh *RateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	rh.rated(w, req, rh.h.ServeHTTP)
}

func (rh *RateHandler) Wrap(wrapped http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rh.rated(w, r, wrapped)
	}
}

func (rh *RateHandler) rated(w http.ResponseWriter, req *http.Request, fn http.HandlerFunc) {
	if !rh.limiter.Allowed() {
		log.Println("queue is full, return 503")
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Queue is full, try again later\n\n"))
		return
	}
	defer func() {
		rh.limiter.Done()
	}()
	fn(w, req)
}

func (rh *RateHandler) CheckQueue(w http.ResponseWriter, req *http.Request) {
	if !rh.limiter.Check() {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	log.Println("queue is full, return 503")
	w.WriteHeader(http.StatusServiceUnavailable)
}
