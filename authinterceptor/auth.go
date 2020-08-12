package authinterceptor

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	AuthorizationHeadString = "authorization"
)

var (
	//含TOKEN 的认证协议类型集合
	tokenTypeList []string = []string{"Bearer"}
)

//Registered interceptor
func NewUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	newCtx, err := auth(ctx)
	if err != nil {
		return nil, err
	}
	// 继续处理请求
	return handler(newCtx, req)
}

func auth(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok { ///不存在MD
		return ctx, nil
	}
	authValue, ok := md[AuthorizationHeadString]
	if !ok {
		return ctx, nil
	}
	if len(authValue) != 1 {
		return ctx, status.Error(codes.Unauthenticated, "authorization head md size error")
	}
	auth := authValue[0]
	vs := strings.Split(auth, " ")
	if len(vs) < 2 {
		return ctx, status.Error(codes.Unauthenticated, "authorization head check error")
	}
	authType := strings.TrimSpace(vs[0])
	info := CallBackInfo{}
	authParam := ""
	if hasTokenInAuthProto(authType) {
		info.TokenType = authType
		info.Token = vs[1]
		if len(vs) > 2 {
			authParam = vs[2]
		}
	} else {
		info.TokenType = authType
		authParam = vs[1]
	}
	if authParam == "" && info.Token == "" {
		return ctx, status.Error(codes.Unauthenticated, "authorization head data check error")
	}
	info.TokenCtx = string2Map(authParam)
	result, err := callProvider(info)
	if err != nil {
		return ctx, err
	}
	if len(result) == 0 {
		return ctx, status.Error(codes.Unauthenticated, "authorization head check auth error")
	}
	/// 填充结果到 CTX MD ，作为认证结果
	for k, v := range result {
		md.Append(k, v)
	}
	/// 清除原有认证请求信息
	delete(md, AuthorizationHeadString)

	return metadata.NewIncomingContext(ctx, md), nil
}

func hasTokenInAuthProto(proto string) bool {
	for i := 0; i < len(tokenTypeList); i++ {
		if tokenTypeList[i] == proto {
			return true
		}
	}
	return false
}

func string2Map(str string) map[string]string {
	ctx := make(map[string]string)
	if len(str) > 0 {
		params := strings.Split(str, ",")
		for i := 0; i < len(params); i++ {
			kvs := strings.Split(params[i], "=")
			if len(kvs) == 1 {
				ctx[kvs[0]] = " "
			} else {
				ctx[kvs[0]] = kvs[1]
			}
		}
	}
	return ctx
}
