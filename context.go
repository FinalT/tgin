package tgin

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Context struct {
	Req *http.Request
	// Resp 原生的 ResponseWriter。当你直接使用 Resp 的时候，
	// 那么相当于你绕开了 RespStatusCode 和 RespData。
	// 响应数据直接被发送到前端，其它中间件将无法修改响应
	// 其实我们也可以考虑将这个做成私有的
	Resp http.ResponseWriter
	// 缓存的响应部分
	// 这部分数据会在最后刷新
	RespStatusCode int
	// RespData []byte
	RespData []byte

	PathParams map[string]string
	// 命中的路由
	MatchedRoute string

	// 万一将来有需求，可以考虑支持这个，但是需要复杂一点的机制
	// Body []byte 用户返回的响应
	// Err error 用户执行的 Error

	// 缓存的数据
	cacheQueryValues url.Values

	// 页面渲染的引擎
	tplEngine TemplateEngine
}

func (c *Context) BindJSON(val any) error {
	if c.Req.Body == nil {
		return errors.New("web: body 为 nil")
	}
	decoder := json.NewDecoder(c.Req.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(val)
}

func (c *Context) FormValue(key string) StringValue {
	if err := c.Req.ParseForm(); err != nil {
		return StringValue{err: err}
	}
	return StringValue{val: c.Req.FormValue(key)}
}

func (c *Context) QueryValue(key string) StringValue {
	if c.cacheQueryValues == nil {
		c.cacheQueryValues = c.Req.URL.Query()
	}
	vals, ok := c.cacheQueryValues[key]
	if !ok {
		return StringValue{err: errors.New("web: 找不到这个 key")}
	}
	return StringValue{val: vals[0]}
}

func (c *Context) PathValue(key string) StringValue {
	val, ok := c.PathParams[key]
	if !ok {
		return StringValue{err: errors.New("web: 找不到这个 key")}
	}
	return StringValue{val: val}
}

func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Resp, cookie)
}

func (c *Context) RespJSONOK(val any) error {
	return c.RespJSON(http.StatusOK, val)
}

func (c *Context) RespJSON(code int, val any) error {
	bs, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c.RespStatusCode = code
	c.RespData = bs
	return err
}

func (c *Context) Render(tpl string, data any) error {
	var err error
	c.RespData, err = c.tplEngine.Render(c.Req.Context(), tpl, data)
	c.RespStatusCode = 200
	if err != nil {
		c.RespStatusCode = 500
	}
	return err
}

// func (c *Context) QueryValueAsInt64(key string) (int64, error) {
// 	val, err := c.QueryValue(key)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return strconv.ParseInt(val, 10, 64)
// }

type StringValue struct {
	val string
	err error
}

func (s StringValue) String() (string, error) {
	return s.val, s.err
}

func (s StringValue) ToInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}

// 不能用泛型
// func (s StringValue) To[T any]() (T, error) {
//
// }
