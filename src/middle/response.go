package middle

import "github.com/lunny/tango"

var (
	_ IResponse = (*Responsor)(nil)
)

type IResponse interface {
	setContext(ctx *tango.Context)

	JSON(data interface{})
	JSONError(status int, err error)
}

type Responsor struct {
	ctx *tango.Context
}

func (r *Responsor) setContext(ctx *tango.Context) {
	r.ctx = ctx
}

type jsonResult struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func (r *Responsor) JSON(data interface{}) {
	if err := r.ctx.ServeJson(jsonResult{true, data, ""}); err != nil {
		panic(err)
	}
}

func (r *Responsor) JSONError(status int, myErr error) {
	r.ctx.WriteHeader(status)
	if err := r.ctx.ServeJson(jsonResult{false, nil, myErr.Error()}); err != nil {
		panic(err)
	}
}

func Responser() tango.HandlerFunc {
	return func(ctx *tango.Context) {
		if resp, ok := ctx.Action().(IResponse); ok {
			resp.setContext(ctx)
		}
		ctx.Next()
	}
}
