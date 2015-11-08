package server

import (
	"github.com/lunny/tango"
)

type Server struct {
	addr         string
	tango        *tango.Tango
	Static       *Static
	Helper       *Helper
	ErrorHandler tango.HandlerFunc
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
		tango: tango.New([]tango.Handler{
			tango.Return(),
			tango.Param(),
			tango.Contexts(),
			tango.Recovery(true),
		}...),
	}
}

func (s *Server) Run() {
	if s.ErrorHandler != nil {
		s.tango.ErrHandler = s.ErrorHandler
	}
	s.tango.Use(logger())
	if s.Static != nil {
		s.tango.Use(s.Static)
	}
	if s.Helper != nil {
		s.tango.Use(s.Helper)
	}
	s.tango.Run(s.addr)
}
