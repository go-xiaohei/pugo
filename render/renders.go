package render

import (
	"errors"
	"github.com/Unknwon/com"
	"path"
)

var (
	ErrRenderDirMissing = errors.New("render-dir-missing")
)

type Renders struct {
	dir     string
	renders map[string]*Render
	current string
	reload  bool
}

func NewRenders(dir, current string, reload bool) (*Renders, error) {
	r := &Renders{
		dir:     dir,
		renders: make(map[string]*Render),
		current: current,
	}
	if _, err := r.NewRender(current, reload); err != nil {
		return nil, err
	}
	return r, nil
}

func (rs *Renders) NewRender(dir string, reload bool) (*Render, error) {
	if rs.renders[dir] == nil {
		fullDir := path.Join(rs.dir, dir)
		if !com.IsDir(fullDir) {
			return nil, ErrRenderDirMissing
		}
		rs.renders[dir] = NewRender(fullDir, reload)
	}
	return rs.renders[dir], nil
}

func (rs *Renders) Current() *Render {
	return rs.renders[rs.current]
}
