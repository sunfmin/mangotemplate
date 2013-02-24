package mangotemplate

import (
	"html/template"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
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

	tpl := template.New(name)

	templatesToBeParsed := make([]string, 1)
	templatesToBeParsed[0] = TemplatePath + name + TemplateSuffix

	for len(templatesToBeParsed) != 0 {
		content, err := ioutil.ReadFile(templatesToBeParsed[0])
		templatesToBeParsed = templatesToBeParsed[1:]
		check(err)

		_, err = tpl.Parse(string(content))
		check(err)

		for _, matched := range regexp.MustCompile(`\{\{\ *template\ +\"(.*)\"\ *\}\}`).FindAllStringSubmatch(string(content), -1) {
			partialName := matched[1]
			partialParsed := false
			for _, t := range tpl.Templates() {
				if t.Name() == partialName {
					partialParsed = true
					break
				}
			}

			if !partialParsed {
				names := strings.Split(partialName, "/")
				partialFilePath := TemplatePath + strings.Join(names[:len(names)-1], "/") + "/_" + names[len(names)-1] + TemplateSuffix
				templatesToBeParsed = append(templatesToBeParsed, partialFilePath)
			}
		}
	}

	err = tpl.Execute(wr, data)
	check(err)

	return
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
