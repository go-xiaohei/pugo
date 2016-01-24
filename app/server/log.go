package server

import (
	"time"

	"net/http"

	"gopkg.in/inconshreveable/log15.v2"
)

var (
	logFormat = "Http|%s|%d|%s|%s|%.1fms" //
)

// log middleware handler
func logger(w *responseWriter, r *http.Request) {
	p := r.URL.Path
	if len(r.URL.RawQuery) > 0 {
		p = p + "?" + r.URL.RawQuery
	}

	if w.status == 0 {
		http.NotFound(w, r)
	}

	// skip static files
	if r.Method == "GET" {
		if p == "/favicon.ico" || p == "/robots.txt" {
			return
		}
	}

	statusCode := w.status
	if statusCode >= 200 && statusCode < 400 {
		log15.Info(logFormat, r.Method, statusCode, p, r.RemoteAddr, time.Since(w.startTime).Seconds()*1000)
		return
	}
	if statusCode < 500 {
		log15.Warn(logFormat, r.Method, statusCode, p, r.RemoteAddr, time.Since(w.startTime).Seconds()*1000)
		return
	}
	log15.Info(logFormat+"|%v", r.Method, statusCode, p, r.RemoteAddr, time.Since(w.startTime).Seconds()*1000, w.error)

}
