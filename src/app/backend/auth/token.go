package auth

import (
	authApi "github.com/wzt3309/k8sconsole/src/app/backend/auth/api"
	"k8s.io/client-go/tools/clientcmd/api"
)

// Implements Authenticator interface.
type tokenAuthenticator struct {
	token string
}

// GetAuthInfo implements Authenticator interface.
func (self tokenAuthenticator) GetAuthInfo() (api.AuthInfo, error) {
	return api.AuthInfo{
		Token: self.token,
	}, nil
}

// NewTokenAuthenticator returns Authenticator based on LoginSpec.
func NewTokenAuthenticator(spec *authApi.LoginSpec) authApi.Authenticator {
	return &tokenAuthenticator{
		token: spec.Token,
	}
}
