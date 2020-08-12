package sdk

import "context"

//Implementing PerRPCCredentials interface
type TokenAuth struct {
	authCtx   map[string]string
	tokenType string
	token     string
}

func (t *TokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	r := make(map[string]string)
	count := len(t.authCtx)
	if count == 0 {
		r["authorization"] = t.tokenType + " " + t.token
	} else {
		r["authorization"] = t.tokenType + " " + t.token + " " + map2String(t.authCtx)
	}
	return r, nil
}

func (t *TokenAuth) RequireTransportSecurity() bool {
	return true
}
