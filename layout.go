package mangotemplate

import (
	"bytes"
	. "github.com/paulbellamy/mango"
	"html/template"
)

type wrapperData struct {
	Layout interface{}
	Body   template.HTML
}

type LayoutDataProvider interface {
	LayoutData(env Env) interface{}
}

func MakeLayout(tpl *template.Template, name string, ldp LayoutDataProvider) Middleware {
	return func(env Env, app App) (status Status, headers Headers, body Body) {
		status, headers, body = app(env)

		if status == 0 {
			status = 200
		}

		if status != 200 {
			return
		}

		b := bytes.NewBuffer([]byte{})
		tpl.ExecuteTemplate(b, name, &wrapperData{ldp.LayoutData(env), template.HTML(body)})
		return status, headers, Body(b.String())
	}
}
