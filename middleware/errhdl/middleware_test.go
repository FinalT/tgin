package errhdl

import (
	"bytes"
	web "github.com/FinalT/tgin"
	"html/template"
	"testing"
)

func TestNewMiddlewareBuilder(t *testing.T) {
	s := web.NewHTTPServer()
	s.Get("/user", func(ctx *web.Context) {
		ctx.RespData = []byte("hello, world")
	})
	page := `
<html>
	<h1>404 NOT FOUND</h1>
</html>
`
	tpl, err := template.New("404").Parse(page)
	if err != nil {
		t.Fatal(err)
	}
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, nil)
	if err != nil {
		t.Fatal(err)
	}
	s.Use(NewMiddlewareBuilder().
		RegisterError(404, buffer.Bytes()).Build())

	s.Start(":8081")
}
