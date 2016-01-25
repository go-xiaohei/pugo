package deploy

import "github.com/go-xiaohei/pugo/app/builder"

type Git struct {
}

func (g *Git) Name() string {
	return "Git"
}

func (g *Git) Detect(*builder.Context) (Task, bool) {
	return nil, false
}

func (g *Git) Action(*builder.Context) error {
	return nil
}
