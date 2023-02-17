package session

import (
	"context"
	"net/http"
)

type Session interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, val string) error
	ID() string
}

// Store 管理 Session
// 从设计的角度来说，Generate 方法和 Refresh 在处理 Session 过期时间上有点关系
// 也就是说，如果 Generate 设计为接收一个 expiration 参数，
// 那么 Refresh 也应该接收一个 expiration 参数。
// 因为这意味着用户来管理过期时间
type Store interface {
	// Generate 生成一个 session
	Generate(ctx context.Context, id string) (Session, error)
	// Refresh 这种设计是一直用同一个 id 的
	// 如果想支持 Refresh 换 ID，那么可以重新生成一个，并移除原有的
	// 又或者 Refresh(ctx context.Context, id string) (Session, error)
	// 其中返回的是一个新的 Session
	Refresh(ctx context.Context, id string) error
	Remove(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (Session, error)
}

type Propagator interface {
	// Inject 将 session id 注入到里面
	// Inject 必须是幂等的
	Inject(id string, writer http.ResponseWriter) error
	// Extract 将 session id 从 http.Request 中提取出来
	// 例如从 cookie 中将 session id 提取出来
	Extract(req *http.Request) (string, error)

	// Remove 将 session id 从 http.ResponseWriter 中删除
	// 例如删除对应的 cookie
	Remove(writer http.ResponseWriter) error
}
