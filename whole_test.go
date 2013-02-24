package mangotemplate

import (
	. "github.com/paulbellamy/mango"
	"github.com/shaoshing/gotest"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"
)

func home(env Env) (status Status, headers Headers, body Body) {
	r := RenderToString("index", []string{"44444", "55555"})
	ForRender(env, "index", []string{"11111", "22222", "33333", r})
	return 200, Headers{}, Body("")
}

type header struct {
	Username string
}

func (h *header) LayoutData(env Env) interface{} {
	return h
}

func mux() *http.ServeMux {
	s := new(Stack)
	tpl, err := template.ParseGlob("test_templates/*.html")
	if err != nil {
		panic(err)
	}

	l := MakeLayout(tpl, "layout", &header{"sunfmin"})
	rdr := MakeRenderer(tpl)

	s.Middleware(l, rdr)

	mux := http.DefaultServeMux
	mux.HandleFunc("/home", s.HandlerFunc(home))
	return mux
}

var ts = httptest.NewServer(mux())

func get(url string) string {

	res, _ := http.Get(ts.URL + url)
	b, _ := ioutil.ReadAll(res.Body)

	return string(b)
}

func TestLayout(t *testing.T) {
	assert.Test = t

	body := get("/home")

	assert.Contain("sunfmin", body)
	assert.Contain("<li>22222</li>", body)
	assert.Contain("44444", body)
}

func TestAutoReload(t *testing.T) {
	assert.Test = t

	preBody := get("/home")

	exec.Command("cp", "test_templates/index.html.reload", "test_templates/index.html").Run()
	defer exec.Command("git", "checkout", "test_templates").Run()

	AutoReload = true
	TemplatePath = "test_templates/"
	defer func() {
		AutoReload = false
	}()

	body := get("/home")

	assert.NotContain("reload index", preBody)
	assert.Contain("reload index", body)
}
