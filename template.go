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
	templateKeyword = regexp.MustCompile(`\{\{\ *template\ +\"([^\}]*)\"[^\}]*\}\}`)
	templateFuncMap = map[*template.Template]template.FuncMap{}
)

func Funcs(tpl *template.Template, funcmap template.FuncMap) {
	tpl.Funcs(funcmap)
	templateFuncMap[tpl] = funcmap
}

func Render(preloadedTpl *template.Template, wr io.Writer, name string, data interface{}) (err error) {
	if !AutoReload {
		err = preloadedTpl.ExecuteTemplate(wr, name, data)
		return
	}

	tpl := template.New(name)
	addFunc(tpl, preloadedTpl)
	parseTemplates(tpl)

	err = tpl.Execute(wr, data)
	check(err)

	return
}

func addFunc(tpl *template.Template, preloadedTpl *template.Template) {
	funcmap, ok := templateFuncMap[preloadedTpl]
	if ok {
		tpl.Funcs(funcmap)
	}
}

func parseTemplates(tpl *template.Template) {
	templatesToBeParsed := []string{tpl.Name()}

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
			_, err := tpl.Parse(content)
			check(err)

			for _, matched := range templateKeyword.FindAllStringSubmatch(content, -1) {
				templatesToBeParsed = append(templatesToBeParsed, matched[1])
			}
		}
	}
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
