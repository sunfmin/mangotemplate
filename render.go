package mangotemplate

import (
	"bytes"
	. "github.com/paulbellamy/mango"
	"html/template"
	"log"
)

const (
	TEMPLATE_NAME_KEY = "mangotemplate.name"
	TEMPLATE_DATA_KEY = "mangotemplate.data"
)

func templateName(env Env) (name string) {
	name, _ = env[TEMPLATE_NAME_KEY].(string)
	return
}

func templateData(env Env) (r interface{}) {
	r = env[TEMPLATE_DATA_KEY]
	return
}

func ForRender(env Env, name string, data interface{}) {
	env[TEMPLATE_NAME_KEY] = name
	env[TEMPLATE_DATA_KEY] = data
}

func MakeRenderer(tpl *template.Template) Middleware {
	return func(env Env, app App) (status Status, headers Headers, body Body) {
		status, headers, body = app(env)

		if status == 0 {
			status = 200
		}

		if status != 200 {
			return
		}

		name := templateName(env)
		if name == "" {
			return
		}

		data := templateData(env)
		b := bytes.NewBuffer([]byte{})
		err := tpl.ExecuteTemplate(b, name, data)
		if err != nil {
			log.Printf("mangotemplate: render %s failed, %s", name, err)
		}

		log.Printf("Rendered %s", name)
		return status, headers, Body(b.String())
	}
}
