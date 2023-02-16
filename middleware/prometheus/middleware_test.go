package prometheus

import (
	web "github.com/FinalT/tgin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"testing"
	"time"
)

// 启动之后，访问一下 localhost:8081/user
// 然后再访问一下 localhost:2112/metrics
// 就能看到类似的输出，注意找一下
// # HELP web_http_request 这是测试例子
// # TYPE web_http_request summary
// web_http_request_sum{instance_id="1234567",method="GET",pattern="/user",status="0"} 1000
// web_http_request_count{instance_id="1234567",method="GET",pattern="/user",status="0"} 1
// web_http_request_sum{instance_id="1234567",method="GET",pattern="unknown",status="404"} 0
// web_http_request_count{instance_id="1234567",method="GET",pattern="unknown",status="404"} 1
// 如果你启动了 prometheus 服务器，那么就配置它来采集这个 2112 端口和 /metrics 路径
func TestMiddlewareBuilder_Build(t *testing.T) {
	s := web.NewHTTPServer()
	s.Get("/", func(ctx *web.Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})
	s.Get("/user", func(ctx *web.Context) {
		time.Sleep(time.Second)
	})

	s.Use((&MiddlewareBuilder{
		Subsystem: "web",
		Name:      "http_request",
		Help:      "这是测试例子",
		ConstLabels: map[string]string{
			"instance_id": "1234567",
		},
	}).Build())
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		// 一般来说，在实际中我们都会单独准备一个端口给这种监控
		http.ListenAndServe(":2112", nil)
	}()
	s.Start(":8081")
}
