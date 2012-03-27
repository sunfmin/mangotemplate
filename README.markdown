This package provide Layout and Render template feature as Middleware for mongo (https://github.com/paulbellamy/mango)

To use the Layout middleware, first create a layout middleware

	tpl, err := template.ParseGlob("test_templates/*.html")
	if err != nil {
		panic(err)
	}

	l := MakeLayout(tpl, "mainlayout", &header{"sunfmin"})

The Last parameter is the LayoutDataProvider interface, which is used to provide data for when rendering layout

next put the Layout middleware into your mango stack:

	s := new(Stack)
	s.Middleware(l)

And use the HandlerFunc of mango:

	mux := http.DefaultServeMux
	mux.HandleFunc("/home", s.HandlerFunc(func(env Env) (status Status, headers Headers, body Body) {
		return 200, Headers{}, Body("Hello, I am inside a layout")
	}))


For documentation:

http://gopkgdoc.appspot.com/pkg/github.com/sunfmin/mangotemplate

