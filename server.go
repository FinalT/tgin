package tgin

import (
	"log"
	"net/http"
)

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler
	// Start 启动服务器
	// addr 是监听地址。如果只指定端口，可以使用 ":8081"
	// 或者 "localhost:8082"
	Start(addr string) error

	// addRoute 注册一个路由
	// method 是 HTTP 方法
	addRoute(method string, path string, handler HandleFunc)
	// 我们并不采取这种设计方案
	// addRoute(method string, path string, handlers... HandleFunc)
}

// 确保 HTTPServer 肯定实现了 Server 接口
var _ Server = &HTTPServer{}

type ServerOption func(server *HTTPServer)

type HTTPServer struct {
	router
	mdls      []Middleware
	tplEngine TemplateEngine
}

func NewHTTPServer(opts ...ServerOption) *HTTPServer {
	s := &HTTPServer{
		router: newRouter(),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func ServerWithTemplateEngine(engine TemplateEngine) ServerOption {
	return func(server *HTTPServer) {
		server.tplEngine = engine
	}
}

func (s *HTTPServer) Use(mdls ...Middleware) {
	if s.mdls == nil {
		s.mdls = mdls
		return
	}
	s.mdls = append(s.mdls, mdls...)
}

// UseV1 会执行路由匹配，只有匹配上了的 mdls 才会生效
// 这个只需要稍微改造一下路由树就可以实现
func (s *HTTPServer) UseV1(path string, mdls ...Middleware) {
	panic("implement me")
}

// ServeHTTP HTTPServer 处理请求的入口
func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:       request,
		Resp:      writer,
		tplEngine: s.tplEngine,
	}
	// 最后一个应该是 HTTPServer 执行路由匹配，执行用户代码
	root := s.serve
	// 从后往前组装
	for i := len(s.mdls) - 1; i >= 0; i-- {
		root = s.mdls[i](root)
	}
	// 第一个应该是回写响应的
	// 因为它在调用next之后才回写响应，
	// 所以实际上 flashResp 是最后一个步骤
	var m Middleware = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
			s.flashResp(ctx)
		}
	}
	root = m(root)
	root(ctx)
}

// Start 启动服务器
func (s *HTTPServer) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *HTTPServer) Post(path string, handler HandleFunc) {
	s.addRoute(http.MethodPost, path, handler)
}

func (s *HTTPServer) Get(path string, handler HandleFunc) {
	s.addRoute(http.MethodGet, path, handler)
}

func (s *HTTPServer) serve(ctx *Context) {
	mi, ok := s.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || mi.n == nil || mi.n.handler == nil {
		ctx.RespStatusCode = 404
		return
	}
	ctx.PathParams = mi.pathParams
	ctx.MatchedRoute = mi.n.route
	mi.n.handler(ctx)
}

func (s *HTTPServer) flashResp(ctx *Context) {
	if ctx.RespStatusCode > 0 {
		ctx.Resp.WriteHeader(ctx.RespStatusCode)
	}
	_, err := ctx.Resp.Write(ctx.RespData)
	if err != nil {
		log.Fatalln("回写响应失败", err)
	}
}
