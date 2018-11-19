package auth

import (
	"errors"
	authApi "github.com/wzt3309/k8sconsole/src/app/backend/auth/api"
	clientApi "github.com/wzt3309/k8sconsole/src/app/backend/client/api"
	kcErrors "github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"k8s.io/client-go/tools/clientcmd/api"
)

// Implements AuthManager interface
type authManager struct {
	tokenManager              	authApi.TokenManager
	clientManager             	clientApi.ClientManager
	authenticationModes       	authApi.AuthenticationModes
	authenticationSkippable 		bool
}

// Login implement AuthManager
func (self authManager) Login(spec *authApi.LoginSpec) (*authApi.AuthResponse, error) {
	authenticator, err := self.getAuthenticator(spec)
	if err != nil {
		return nil, err
	}

	authInfo, err := authenticator.GetAuthInfo()
	if err != nil {
		return nil, err
	}

	err = self.healthCheck(authInfo)
	nonCriticalErrors, criticalError := kcErrors.HandleError(err)
	if criticalError != nil || len(nonCriticalErrors) > 0 {
		return &authApi.AuthResponse{Errors: nonCriticalErrors}, criticalError
	}

	token, err := self.tokenManager.Generate(authInfo)
	if err != nil {
		return nil, err
	}

	return &authApi.AuthResponse{JWEToken: token, Errors: nonCriticalErrors}, nil
}

// Refresh implements AuthManager
func (self authManager) Refresh(jweToken string) (string, error) {
	return self.tokenManager.Refresh(jweToken)
}

// AuthenticationModes implements AuthManager
func (self authManager) AuthenticationModes() []authApi.AuthenticationMode {
	return self.authenticationModes.Array()
}

// AuthenticationSkippable implements AuthManager
func (self authManager) AuthenticationSkippable() bool {
	return self.authenticationSkippable
}

func (self authManager) getAuthenticator(spec *authApi.LoginSpec) (authApi.Authenticator, error) {
	if len(self.authenticationModes) == 0 {
		return nil, errors.New("All authentication options disabled.")
	}

	switch {
	case len(spec.Username) > 0 && len(spec.Password) > 0 && self.authenticationModes.IsEnabled(authApi.Basic):
		return NewBasicAuthenticator(spec), nil
	case len(spec.Token) > 0 && self.authenticationModes.IsEnabled(authApi.Token):
		return NewTokenAuthenticator(spec), nil
	}

	return nil, errors.New("No enough data to create supported authenticator.")
}

// Checks if user data extracted from provided AuthInfo is valid and user is correctly authenticated
// by k8s apiserver
func (self authManager) healthCheck(authInfo api.AuthInfo) error {
	return self.clientManager.HasAccess(authInfo)
}

func NewAuthManager(clientManager clientApi.ClientManager, tokenManager authApi.TokenManager,
	authenticationModes authApi.AuthenticationModes, authenticationSkippable bool) authApi.AuthManager {
	return &authManager{
		tokenManager:            tokenManager,
		clientManager:           clientManager,
		authenticationModes:     authenticationModes,
		authenticationSkippable: authenticationSkippable,
	}
}
