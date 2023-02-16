package accesslog

import (
	"encoding/json"
	web "github.com/FinalT/tgin"
	"log"
)

type MiddlewareBuilder struct {
	logFunc func(accessLog string)
}

func (b *MiddlewareBuilder) LogFunc(logFunc func(accessLog string)) *MiddlewareBuilder {
	b.logFunc = logFunc
	return b
}

func NewBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		logFunc: func(accessLog string) {
			log.Println(accessLog)
		},
	}
}

type accessLog struct {
	Host       string
	Route      string
	HTTPMethod string `json:"http_method"`
	Path       string
}

func (b *MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			defer func() {
				l := accessLog{
					Host:       ctx.Req.Host,
					Route:      ctx.MatchedRoute,
					Path:       ctx.Req.URL.Path,
					HTTPMethod: ctx.Req.Method,
				}
				val, _ := json.Marshal(l)
				b.logFunc(string(val))
			}()
			next(ctx)
		}
	}
}
