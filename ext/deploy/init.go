package deploy

import (
	"fmt"
	"github.com/go-xiaohei/pugo/app/builder"
)

func Init() {
	builder.Before(Detect)
	builder.After(Action)
}

func Detect(ctx *builder.Context) {
	for _, m := range manager.tasks {
		fmt.Println(m.Detect(ctx))
	}
}

func Action(ctx *builder.Context) {
	println("action")
}
