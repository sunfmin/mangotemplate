package mangotemplate

import (
	"html/template"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
)

var (
	AutoReload      = false
	TemplatePath    = "templates/"
	TemplateSuffix  = ".html"
	templateKeyword = regexp.MustCompile(`\{\{\ *template\ +\"(.*)\".*}\}`)
)

func Render(preloadedTpl *template.Template, wr io.Writer, name string, data interface{}) (err error) {
	if !AutoReload {
		err = preloadedTpl.ExecuteTemplate(wr, name, data)
		return
	}

	templatesToBeParsed := []string{name}
	for len(templatesToBeParsed) != 0 {
		name := templatesToBeParsed[0]
		templatesToBeParsed = templatesToBeParsed[1:]

		parsed := false
		for _, t := range tpl.Templates() {
			if tpl.Name() != name && t.Name() == name {
				parsed = true
				break
			}
		}

		if !parsed {

			content := readTemplate(name)
			_, err = tpl.Parse(content)
			check(err)

			for _, matched := range templateKeyword.FindAllStringSubmatch(content, -1) {
				templatesToBeParsed = append(templatesToBeParsed, matched[1])
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

func readTemplate(name string) string {
	paths := templatePaths(name)
	for _, path := range paths {
		content, err := ioutil.ReadFile(path)
		if err == nil {
			return string(content)
		}
	}
	panic("Could not find template [" + name + "] from: " + strings.Join(paths, ", "))
	return ""
}

// Example
// input: index/menu
// output:
//    [templates/index/menu.html,
//		 templates/index/_menu.html,
//		 templates/layout/index/menu.html,
//		 templates/layout/index/_menu.html]
func templatePaths(name string) []string {
	names := strings.Split(name, "/")
	partialName := strings.Join(names[:len(names)-1], "/") + "/_" + names[len(names)-1]

	return []string{
		TemplatePath + name + TemplateSuffix,
		TemplatePath + partialName + TemplateSuffix,
		TemplatePath + "layout/" + name + TemplateSuffix,
		TemplatePath + "layout/" + partialName + TemplateSuffix,
	}
}
