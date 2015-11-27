package server

import (
	"path"

	"github.com/Unknwon/com"
	"github.com/lunny/log"
	"github.com/lunny/tango"
	"gopkg.in/inconshreveable/log15.v2"
)

type Server struct {
	Tango *tango.Tango

	dstDir string
	prefix string
}

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

func (s *Server) SetPrefix(prefix string) {
	s.prefix = prefix
}

func (s *Server) Run(addr string) {
	log15.Debug("Server.Start." + addr)
	s.Tango.Use(logger())
	s.Tango.Get("/", s.globalHandler)
	s.Tango.Get("/*name", s.globalHandler)
	s.Tango.Run(addr)
}

func (s *Server) globalHandler(ctx *tango.Context) {
	file := path.Join(s.dstDir, ctx.Param("*name", "index.html"))
	if com.IsFile(file) {
		ctx.ServeFile(file)
		return
	}
	ctx.Redirect("/error/404.html")
}
