package sdk

import (
	"strings"

	"google.golang.org/grpc/credentials"
)

func NewTokenAuth(authType string, token string, authCtx map[string]string) credentials.PerRPCCredentials {
	tokenAuth := new(TokenAuth)
	tokenAuth.token = token
	tokenAuth.tokenType = authType
	tokenAuth.authCtx = authCtx
	return tokenAuth
}

func NewBearerAuth(token string) credentials.PerRPCCredentials {
	return NewTokenAuth("Bearer", token, nil)
}

func map2String(ctx map[string]string) string {
	count := len(ctx)
	valueStr := ""
	if count > 0 {
		nCount := 0
		for k, v := range ctx {
			nCount++
			if strings.TrimSpace(v) == "" {
				valueStr = valueStr + k
			} else {
				valueStr = valueStr + k + "=" + v
			}
			if nCount < count {
				valueStr = valueStr + ","
			}
		}
	}
	return valueStr
}
