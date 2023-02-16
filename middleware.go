package tgin

type Middleware func(next HandleFunc) HandleFunc
