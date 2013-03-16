package mangotemplate

import (
	"bytes"
	. "github.com/paulbellamy/mango"
	"github.com/shaoshing/gotest"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"
)

func home(env Env) (status Status, headers Headers, body Body) {
	r := RenderToString("index", []string{"44444", "55555"})
	ForRender(env, "index", []string{"11111", "22222", "33333", r})
	return 200, Headers{}, Body("")
}

func notification(env Env) (status Status, headers Headers, body Body) {
	ForRender(env, "notification", nil)
	return 200, Headers{}, Body("")
}

type header struct {
	Username string
}

func (h *header) LayoutData(env Env) interface{} {
	return h
}

func mux() *http.ServeMux {
	TemplatePath = "test_templates/"

	s := new(Stack)
	tpl := template.New("")
	_, err := ParseGlob(tpl, "test_templates/*.html")
	if err != nil {
		panic(err)
	}
	_, err = ParseGlob(tpl, "test_templates/*/*.html")
	if err != nil {
		panic(err)
	}

	l := MakeLayout(tpl, "layout", &header{"sunfmin"})
	rdr := MakeRenderer(tpl)

	s.Middleware(l, rdr)

	mux := http.DefaultServeMux
	mux.HandleFunc("/home", s.HandlerFunc(home))
	mux.HandleFunc("/notification", s.HandlerFunc(notification))
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
	assert.NotContain("reload index", preBody)
	preBody = get("/notification")
	assert.NotContain("reload notification 1", preBody)

	bash("cp test_templates/index.html /tmp/mangotemplate.index.html && cp test_templates/layout.html /tmp/mangotemplate.layout.html && cp test_templates/tips.html /tmp/mangotemplate.tips.html")
	bash("cp test_templates/index.html.reload test_templates/index.html && cp test_templates/layout.html.reload test_templates/layout.html && cp test_templates/tips.html.reload test_templates/tips.html")
	defer bash("cp /tmp/mangotemplate.index.html test_templates/index.html && cp /tmp/mangotemplate.layout.html test_templates/layout.html && cp /tmp/mangotemplate.tips.html test_templates/tips.html")

	AutoReload = true
	defer func() {
		AutoReload = false
	}()

	body := get("/home")

	assert.Contain("reload index", body)  // Body should be changed when template files were changed
	assert.Contain("reload layout", body) // Body should be changed when template files were changed

	assert.Contain("index partial", body)  // Make sure partials rendered inside template will work.
	assert.Contain("inline partial", body) // Make sure inline partials will work.
	assert.Contain("menu", body)           // Read partial without trailing "_"
	assert.Contain("footer", body)         // Read partial from layout folder
	assert.Contain("header", body)         // Read partial from layout folder

	body = get("/notification")
	assert.Contain("inline notification", body)
	assert.Contain("reload notification 1", body)
	assert.Contain("reload notification 2", body)
}

func bash(bash string) string {
	var b bytes.Buffer
	cmd := exec.Command("sh", "-c", bash)
	cmd.Stdout = &b
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	return strings.TrimRight(b.String(), "\n")
}
