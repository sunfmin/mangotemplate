package mangotemplate

import (
	"bytes"
	. "github.com/paulbellamy/mango"
	"html/template"
	"log"
	"strings"
)

const (
	TEMPLATE_NAME_KEY = "mangotemplate.name"
	TEMPLATE_DATA_KEY = "mangotemplate.data"
)

var MOBILE_KEYS = []string{"Android", "iPhone", "iPod", "Windows Phone"}

func IsFromMobile(userAgent string) bool {

	for _, key := range MOBILE_KEYS {
		if strings.Contains(userAgent, key) {
			return true
		}
	}
	return false
}

func templateName(env Env) (name string) {
	name, _ = env[TEMPLATE_NAME_KEY].(string)
	return
}

func templateData(env Env) (r interface{}) {
	r = env[TEMPLATE_DATA_KEY]
	return
}

var RenderToStringTemplate *template.Template

func RenderToString(name string, data interface{}) (r string) {
	if RenderToStringTemplate == nil {
		panic("you must use mangotemplate.RenderToStringTemplate = tpl to set the template.")
	}
	return RenderToStringT(name, RenderToStringTemplate, data)
}

func RenderToStringT(name string, tpl *template.Template, data interface{}) (r string) {
	b := bytes.NewBuffer([]byte{})
	err := Render(tpl, b, name, data)
	if err != nil {
		log.Printf("mangotemplate: RenderToString %s failed, %s", name, err)
		return
	}
	r = b.String()
	return
}

func ForRender(env Env, name string, data interface{}) {
	if IsFromMobile(env.Request().UserAgent()) && !strings.HasPrefix(name, "mobiles_layout/") {
		arr := strings.Split(name, "/")
		name = "mobiles_layout/" + arr[len(arr)-1]
	}
	env[TEMPLATE_NAME_KEY] = name
	env[TEMPLATE_DATA_KEY] = data
}

func MakeRenderer(tpl *template.Template) Middleware {
	if RenderToStringTemplate == nil {
		RenderToStringTemplate = tpl
	}

	defaultTemplate := tpl

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
		tpl := getTemplateFromEnv(env, defaultTemplate)
		b := bytes.NewBuffer([]byte{})
		err := Render(tpl, b, name, data)
		if err != nil {
			log.Printf("mangotemplate: render %s failed, %s", name, err)
			status = 500
			return
		}

		log.Printf("Rendered %s", name)
		return status, headers, Body(b.String())
	}
}
