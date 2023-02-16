package tgin

import (
	"bytes"
	"fmt"
	"html/template"
	"testing"
)

// 这里放着端到端测试的代码

func TestServer(t *testing.T) {
	s := NewHTTPServer()
	s.Get("/", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})
	s.Get("/user", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, user"))
	})

	s.Post("/form", func(ctx *Context) {
		err := ctx.Req.ParseForm()
		if err != nil {
			fmt.Println(err)
		}
	})

	s.Get("/login", func(ctx *Context) {
		tpl := template.New("login")
		tpl, err := tpl.Parse(`
<html>
	<body>
		<form>
			// 在这里继续写页面
		<form>
	</body>
</html>
`)
		if err != nil {
			t.Fatal(err)
		}

		page := &bytes.Buffer{}
		err = tpl.Execute(page, nil)
		if err != nil {
			t.Fatal(err)
		}

		ctx.RespStatusCode = 200
		ctx.RespData = page.Bytes()
	})

	s.Start(":8081")
}

func TestServerWithRenderEngine(t *testing.T) {
	tpl, err := template.ParseGlob("testdata/tpls/*.gohtml")
	if err != nil {
		t.Fatal(err)
	}
	s := NewHTTPServer(ServerWithTemplateEngine(&GoTemplateEngine{T: tpl}))
	s.Get("/login", func(ctx *Context) {
		er := ctx.Render("login.gohtml", nil)
		if er != nil {
			t.Fatal(er)
		}
	})
	err = s.Start(":8081")
	if err != nil {
		t.Fatal(err)
	}
}
