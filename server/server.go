package server

import (
	"path"

	"github.com/Unknwon/com"
	"github.com/lunny/log"
	"github.com/lunny/tango"
	"gopkg.in/inconshreveable/log15.v2"
)

// simple built-in http server
type Server struct {
	Tango *tango.Tango // use tango

	dstDir string
	prefix string
}

// new server
// set dstDir to make sure read correct static file
func New(dstDir string) *Server {
	t := tango.New([]tango.Handler{
		tango.Return(),
		tango.Param(),
		tango.Contexts(),
		tango.Recovery(true),
	}...)
	t.Logger().(*log.Logger).SetOutputLevel(log.Lerror)
	return &Server{
		Tango:  t,
		dstDir: dstDir,
	}
}

// set prefix to trim url
func (s *Server) SetPrefix(prefix string) {
	s.prefix = prefix
}

// set run
func (s *Server) Run(addr string) {
	log15.Debug("Server.Start." + addr)
	s.Tango.Use(logger())
	s.Tango.Get("/", s.globalHandler)
	s.Tango.Get("/*name", s.globalHandler)
	s.Tango.Run(addr)
}

func (s *Server) globalHandler(ctx *tango.Context) {
	param := ctx.Param("*name")
	if path.Ext(param) == "" {
		// visit directory, file index.html
		param = path.Join(param, "index.html")
	}
	file := path.Join(s.dstDir, param)
	// try simple file
	if com.IsFile(file) {
		ctx.ServeFile(file)
		return
	}

	ctx.Redirect("/")
}
