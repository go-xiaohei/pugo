package server

import (
	"github.com/go-xiaohei/pugo-static/builder"
	"github.com/lunny/tango"
	"gopkg.in/inconshreveable/log15.v2"
)

type Server struct {
	addr         string
	tango        *tango.Tango
	builder      *builder.Builder
	Static       []*Static
	Helper       *Helper
	ErrorHandler tango.HandlerFunc

	refreshTime int64
}

func NewServer(addr string, b *builder.Builder) *Server {
	return &Server{
		addr: addr,
		tango: tango.New([]tango.Handler{
			tango.Return(),
			tango.Param(),
			tango.Contexts(),
			tango.Recovery(true),
		}...),
		builder: b,
	}
}

func (s *Server) Run() {
	if s.ErrorHandler != nil {
		s.tango.ErrHandler = s.ErrorHandler
	}
	s.tango.Use(tango.HandlerFunc(s.refresh))
	s.tango.Use(logger())
	if len(s.Static) > 0 {
		for _, ss := range s.Static {
			if base := s.builder.Context().Meta.Base; base != "" {
				ss.setBase(base)
			}
			s.tango.Use(ss)
		}
	}
	if s.Helper != nil {
		s.tango.Use(s.Helper)
	}
	s.tango.Run(s.addr)
}

func (s *Server) refresh(ctx *tango.Context) {
	// if build time changed, refresh server setting
	if s.refreshTime == s.builder.Report().BeginTime.Unix() {
		ctx.Next()
		return
	}
	log15.Debug("Server.Refresh.AfterBuild")

	base := s.builder.Context().Meta.Base
	for _, ss := range s.Static {
		if base != "" {
			ss.setBase(base)
			continue
		}
		if ss.Prefix != ss.realPrefix {
			ss.setBase(base)
		}
	}
	logBase = base

	s.refreshTime = s.builder.Report().BeginTime.Unix()
	ctx.Next()
}
