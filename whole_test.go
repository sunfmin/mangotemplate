package mangotemplate

import (
	. "github.com/paulbellamy/mango"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func home(env Env) (status Status, headers Headers, body Body) {
	ForRender(env, "home/index", []string{"11111", "22222", "33333"})
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

	l := MakeLayout(tpl, "mainlayout", &header{"sunfmin"})
	rdr := MakeRenderer(tpl)

	s.Middleware(l, rdr)

	mux := http.DefaultServeMux
	mux.HandleFunc("/home", s.HandlerFunc(home))
	return mux
}

func TestLayout(t *testing.T) {
	ts := httptest.NewServer(mux())
	defer ts.Close()

	res, _ := http.Get(ts.URL + "/home")

	b, _ := ioutil.ReadAll(res.Body)

	body := string(b)

	if !strings.Contains(body, "sunfmin") {
		t.Errorf("%+v should contain \"sunfmin\"", body)

	}
	if !strings.Contains(body, "<li>22222</li>") {
		t.Errorf("%+v should contain \"11111\"", body)

	}

}
