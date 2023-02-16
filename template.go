package tgin

import (
	"bytes"
	"context"
	"html/template"
	"io/fs"
)

type TemplateEngine interface {
	// Render 渲染页面
	// data 是渲染页面所需要的数据
	Render(ctx context.Context, tplName string, data any) ([]byte, error)
}

type GoTemplateEngine struct {
	T *template.Template
	// 也可以考虑设计为 map[string]*template.Template
	// 但是其实没太大必要，因为 template.Template 本身就提供了按名索引的功能
}

func (g *GoTemplateEngine) Render(ctx context.Context,
	tplName string, data any) ([]byte, error) {
	res := &bytes.Buffer{}
	err := g.T.ExecuteTemplate(res, tplName, data)
	return res.Bytes(), err
}

// 以下这三个方法，可以加可以不加，看你是什么风格的设计者

func (g *GoTemplateEngine) LoadFromGlob(pattern string) error {
	var err error
	g.T, err = template.ParseGlob(pattern)
	return err
}

func (g *GoTemplateEngine) LoadFromFiles(filenames ...string) error {
	var err error
	g.T, err = template.ParseFiles(filenames...)
	return err
}

func (g *GoTemplateEngine) LoadFromFS(fs fs.FS, patterns ...string) error {
	var err error
	g.T, err = template.ParseFS(fs, patterns...)
	return err
}
