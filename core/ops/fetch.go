package ops

import (
	quickjs "github.com/abdfnx/qjs"
	"github.com/abdfnx/renio/tools"

	"github.com/imroc/req"
)

func Fetch(ctx *quickjs.Context, url quickjs.Value) quickjs.Value {
	r, err := req.Get(url.String())
	tools.Check(err)
	resp, err := r.ToString()
	tools.Check(err)
	return ctx.String(resp)
}
