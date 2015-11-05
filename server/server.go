package server

import (
	"github.com/lunny/tango"
)

type Server struct {
	addr        string
	tango       *tango.Tango
	static      *Static
	buildHelper *Helper
}

func NewServer(addr string, static *Static, helper *Helper) *Server {
	return &Server{
		addr:        addr,
		tango:       tango.Classic(),
		static:      static,
		buildHelper: helper,
	}
}

func (s *Server) Static() *Static {
	return s.static
}

func (s *Server) Run() {
	s.tango.Use(s.buildHelper)
	s.tango.Use(s.static)
	s.tango.Run(s.addr)
}
