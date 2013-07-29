package mangotemplate

import (
	"bytes"
	. "github.com/paulbellamy/mango"
	"html/template"
	"log"
	"strings"
)

type wrapperData struct {
	Layout interface{}
	Body   template.HTML
}

type LayoutDataProvider interface {
	LayoutData(env Env) interface{}
}

func MakeLayout(tpl *template.Template, name string, ldp LayoutDataProvider) Middleware {
	defaultTemplate := tpl

	return func(env Env, app App) (status Status, headers Headers, body Body) {
		status, headers, body = app(env)

		if status == 0 {
			status = 200
		}

		if status != 200 {
			return
		}

		tpl := getTemplateFromEnv(env, defaultTemplate)
		b := bytes.NewBuffer([]byte{})
		tempName := name
		if IsFromMobile(env.Request().UserAgent()) && !strings.HasPrefix(tempName, "mobiles_layout/") {
			if tempName != "main" {
				tempName = "home"
			}
			tempName = "mobiles_layout/" + tempName

		}
		err := Render(tpl, b, tempName, &wrapperData{ldp.LayoutData(env), template.HTML(body)})
		if err != nil {
			log.Printf("mangotemplate: layout %s failed, %s", name, err)
		}
		return status, headers, Body(strings.TrimSpace(b.String()))
	}
}
