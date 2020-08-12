package authinterceptor

type CallBackInfo struct {
	TokenType string
	Token     string            ///可以为空
	TokenCtx  map[string]string ///从Authorization 解开的 MAP
}

//AuthProvider do auth callback , support unit test,grpc...
type AuthProvider interface {
	//CheckAuth 输入一些认证参数，返回认证的结果，如认证主体等， error 为 GRPC 标准错误
	//返回的结果K,V 会输出到后续 RPC 调用链的 CTX 的 MD 里面
	CheckAuth(info CallBackInfo) (map[string]string, error)
}

func callProvider(info CallBackInfo) (map[string]string, error) {
	return defaultAuthProvider.CheckAuth(info)
}

type TestProvider struct {
}

func (t *TestProvider) CheckAuth(info CallBackInfo) (map[string]string, error) {
	m := make(map[string]string, 0)
	return m, nil
}

var (
	defaultAuthProvider AuthProvider = new(TestProvider)
)

func RegisterDefaultAuthProvider(p AuthProvider) {
	defaultAuthProvider = p
	if defaultAuthProvider == nil {
		defaultAuthProvider = new(TestProvider)
	}
}
