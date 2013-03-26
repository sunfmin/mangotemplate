package mangotemplate

import (
	"fmt"
	"github.com/paulbellamy/mango"
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

const EnvTemplateName = "MANGOTEMPLATE_TEMPLATE"

var (
	AutoReload      = false
	TemplatePath    = "templates/"
	TemplateSuffix  = ".html"
	templateKeyword = regexp.MustCompile(`\{\{\ *template\ +\"([^\}]*)\"[^\}]*\}\}`)
	templateFuncMap = map[*template.Template]template.FuncMap{}
	// Mapping for the template name and its corresponding file path. Checkout the ParseGlob for how it works.
	templateMapPaths = map[string]string{}
)

func Funcs(tpl *template.Template, funcmap template.FuncMap) {
	tpl.Funcs(funcmap)
	templateFuncMap[tpl] = funcmap
}

// Pase files match the [pattern] and store file path of each template name in templateMapPaths
func ParseGlob(templates *template.Template, pattern string) (*template.Template, error) {
	filePaths, err := filepath.Glob(pattern)
	if err != nil {
		return templates, err
	}

	if len(filePaths) == 0 {
		return templates, fmt.Errorf("mangotemplate.ParseGlob: pattern matches no files: %#q", pattern)
	}

	for _, path := range filePaths {
		_, err := templates.ParseFiles(path)
		if err != nil {
			return templates, err
		}

		for _, parsedTemplate := range templates.Templates() {
			if _, ok := templateMapPaths[parsedTemplate.Name()]; ok {
				continue
			}
			templateMapPaths[parsedTemplate.Name()] = path
		}
	}

	return templates, nil
}

func Render(preloadedTpl *template.Template, wr io.Writer, name string, data interface{}) (err error) {
	if !AutoReload {
		err = preloadedTpl.ExecuteTemplate(wr, name, data)
		return
	}

	tpl := template.New(name)
	addFunc(tpl, preloadedTpl)
	err = parseTemplates(tpl)
	if err != nil {
		return fmt.Errorf("\n==== Could not reload templates:\n%s\n\n\n", err.Error())
	}

	if err = tpl.Execute(wr, data); err != nil {
		panic(err)
	}

	return
}

func addFunc(tpl *template.Template, preloadedTpl *template.Template) {
	funcmap, ok := templateFuncMap[preloadedTpl]
	if ok {
		tpl.Funcs(funcmap)
	}
}

func parseTemplates(tpl *template.Template) (err error) {
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
			content, path, err := readTemplate(name)
			if err != nil {
				return err
			}

			if _, err = tpl.New(path).Parse(content); err != nil {
				if strings.Contains(err.Error(), "redefinition of template") {
					return fmt.Errorf("A redefinition error raised while parsing [%s]. Did you happen to put the HTML comments outside the {{define}} block?)", name)
				}
				return err
			}

			for _, matched := range templateKeyword.FindAllStringSubmatch(content, -1) {
				templatesToBeParsed = append(templatesToBeParsed, matched[1])
			}
		}
	}

	return
}

func readTemplate(name string) (result string, path string, err error) {
	var found bool
	path, found = templateMapPaths[name]
	if !found {
		err = fmt.Errorf("Unable to locate path of template [%s].", name)
		return
	}

	var content []byte
	content, err = ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("Uable to read template content from [%s].", path)
		return
	}

	result = string(content)
	return
}

func getTemplateFromEnv(env mango.Env, defaultTemplate *template.Template) (tpl *template.Template) {
	tpl = defaultTemplate
	if _, ok := env[EnvTemplateName]; ok {
		tpl = env[EnvTemplateName].(*template.Template)
	}
	return
}
