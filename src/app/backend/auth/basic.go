package auth

import (
	authApi "github.com/wzt3309/k8sconsole/src/app/backend/auth/api"
	"k8s.io/client-go/tools/clientcmd/api"
)

// Implements Authenticator interface.
type basicAuthenticator struct {
	username string
	password string
}

// GetAuthInfo implements Authenticator interface.
func (self *basicAuthenticator) GetAuthInfo() (api.AuthInfo, error) {
	return api.AuthInfo{
		Username: self.username,
		Password: self.password,
	}, nil
}

// NewBasicAuthenticator returns Authenticator based on LoginSpec.
func NewBasicAuthenticator(spec *authApi.LoginSpec) authApi.Authenticator {
	return &basicAuthenticator{
		username: spec.Username,
		password: spec.Password,
	}
}
