package server

import (
	"github.com/Unknwon/com"
	"gopkg.in/inconshreveable/log15.v2"
	"net/http"
	"path"
	"strings"
	"time"
)

type Server struct {
	dstDir string
	prefix string
}

func New(dstDir string) *Server {
	s := &Server{
		dstDir: dstDir,
	}
	s.SetPrefix("")
	return s
}

// SetPrefix sets prefix to trim url
func (s *Server) SetPrefix(prefix string) {
	if prefix == "" {
		prefix = "/"
	}
	s.prefix = prefix
}

func (s *Server) GetPrefix() string {
	return s.prefix
}

func (s *Server) serveFile(w http.ResponseWriter, r *http.Request, file string) bool {
	if com.IsFile(file) {
		log15.Debug("Dest.File.[" + file + "]")
		http.ServeFile(w, r, file)
		return true
	}
	return false
}

func (s *Server) serveFiles(w http.ResponseWriter, r *http.Request, param string) bool {
	ext := path.Ext(param)
	if ext == "" || ext == "." {
		if !strings.HasSuffix(param, "/") {
			if s.serveFile(w, r, path.Join(s.dstDir, param+".html")) {
				return true
			}
		}
		if s.serveFile(w, r, path.Join(s.dstDir, param, "index.html")) {
			return true
		}
	}
	if s.serveFile(w, r, path.Join(s.dstDir, param)) {
		return true
	}
	return false
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	w := &responseWriter{
		ResponseWriter: rw,
		startTime:      time.Now(),
	}

	defer func() {
		if err := recover(); err != nil {
			w.error = err
			if w.status == 0 {
				http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			}
		}
		logger(w, r)
	}()

	param := r.URL.Path
	if param == "favicon.ico" || param == "robots.txt" {
		if !s.serveFiles(w, r, param) {
			http.NotFound(w, r)
		}
		return
	}
	if !strings.HasPrefix(param, s.prefix) {
		http.Redirect(w, r, s.prefix, 302)
		return
	}
	param = strings.TrimPrefix(param, s.prefix)
	s.serveFiles(w, r, param)
}

func (s *Server) Run(addr string) {
	log15.Info("Server.Start." + addr)
	http.ListenAndServe(addr, s)
}

type responseWriter struct {
	http.ResponseWriter
	status    int
	startTime time.Time
	error     interface{}
}

func (r *responseWriter) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
