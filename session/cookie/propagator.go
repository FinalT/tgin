package cookie

import (
	"net/http"
)

type PropagatorOption func(propagator *Propagator)

func WithCookieOption(opt func(c *http.Cookie)) PropagatorOption {
	return func(propagator *Propagator) {
		propagator.cookieOpt = opt
	}
}

type Propagator struct {
	cookieName string
	cookieOpt  func(c *http.Cookie)
}

func NewPropagator(cookieName string, opts ...PropagatorOption) *Propagator {
	res := &Propagator{
		cookieName: cookieName,
		cookieOpt:  func(c *http.Cookie) {},
	}
	return res
}

func (c *Propagator) Inject(id string, writer http.ResponseWriter) error {
	cookie := &http.Cookie{
		Name:  c.cookieName,
		Value: id,
	}
	c.cookieOpt(cookie)
	http.SetCookie(writer, cookie)
	return nil
}

func (c *Propagator) Extract(req *http.Request) (string, error) {
	cookie, err := req.Cookie(c.cookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (c *Propagator) Remove(writer http.ResponseWriter) error {
	cookie := &http.Cookie{
		Name:   c.cookieName,
		MaxAge: -1,
	}
	c.cookieOpt(cookie)
	http.SetCookie(writer, cookie)
	return nil
}
