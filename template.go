package mangotemplate

import (
	"html/template"
	"io"
)

var (
	AutoReload     = false
	TemplatePath   = "templates/"
	TemplateSuffix = ".html"
)

func Render(preloadedTpl *template.Template, wr io.Writer, name string, data interface{}) (err error) {
	if !AutoReload {
		err = preloadedTpl.ExecuteTemplate(wr, name, data)
		return
	}

	templateFilePath := TemplatePath + name + TemplateSuffix
	var tpl *template.Template
	tpl, err = template.ParseFiles(templateFilePath)
	if err != nil {
		panic(err)
	}

	// tpl.Funcs(funcMap)
	err = tpl.ExecuteTemplate(wr, name, data)
	return
}
