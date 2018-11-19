package jwe

import (
	"errors"
	authApi "github.com/wzt3309/k8sconsole/src/app/backend/auth/api"
	kcErrors "github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"k8s.io/client-go/tools/clientcmd/api"
	"reflect"
	"testing"
	"time"
)

func getTokenManager() authApi.TokenManager {
	holder := NewRSAKeyHolder()
	return NewJWETokenManager(holder)
}

func asSameError(err1, err2 error) bool {
	return (err1 != nil && err2 != nil && err1.Error() == err2.Error()) ||
		(err1 == nil && err2 == nil)
}

func TestJweTokenManager_Generate(t *testing.T) {
	cases := []struct {
		info        string
		authInfo    api.AuthInfo
		expectedErr error
	}{
		{
			"Should generate encrytped token",
			api.AuthInfo{Token: "test-token"},
			nil,
		},
	}

	for _, c := range cases {
		tokenManager := getTokenManager()
		token, err := tokenManager.Generate(c.authInfo)

		if !asSameError(c.expectedErr, err) {
			t.Errorf("Test case: %s. Expected error to be: %v, but got %v.",
				c.info, c.expectedErr, err)
		}

		if len(token) == 0 {
			t.Errorf("Test case: %s. Expected token not to be empty.", c.info)
		}
	}
}

func TestJweTokenManager_Decrypt(t *testing.T) {
	cases := []struct {
		info        string
		authInfo    api.AuthInfo
		expected    *api.AuthInfo
		expectedErr error
	}{
		{
			"Should decrypt encrypted token",
			api.AuthInfo{Token: "test-token"},
			&api.AuthInfo{Token: "test-token"},
			nil,
		},
	}

	for _, c := range cases {
		tokenManager := getTokenManager()
		token, _ := tokenManager.Generate(c.authInfo)
		authInfo, err := tokenManager.Decrypt(token)

		if !asSameError(c.expectedErr, err) {
			t.Errorf("Test case: %s. Expected error to be: %v, but got %v",
				c.info, c.expectedErr, err)
		}

		if !reflect.DeepEqual(authInfo, c.expected) {
			t.Errorf("Test case: %s. Expected: %v, but got %v.", c.info, c.expected, authInfo)
		}
	}
}

func TestJweTokenManager_Refresh(t *testing.T) {
	cases := []struct {
		info        string
		authInfo    api.AuthInfo
		shouldSleep bool
		expected    bool
		expectedErr error
	}{
		{
			"Shoule refresh valid token",
			api.AuthInfo{Token: "test-token"},
			false,
			true,
			nil,
		},
		{
			info:        "Should return error when no token provided",
			authInfo:    api.AuthInfo{},
			shouldSleep: false,
			expected:    false,
			expectedErr: errors.New("Can not refresh token. No token provided."),
		},
		{
			info:        "Should return error when token has expired",
			authInfo:    api.AuthInfo{Token: "test-token"},
			shouldSleep: true,
			expected:    false,
			expectedErr: errors.New(kcErrors.MSG_TOKEN_EXPIRED_ERROR),
		},
	}

	for _, c := range cases {
		tokenManager := getTokenManager()
		tokenManager.SetTokenTTL(1)
		token, _ := tokenManager.Generate(c.authInfo)

		if len(c.authInfo.Token) == 0 {
			token = ""
		}

		if c.shouldSleep {
			time.Sleep(2 * time.Second)
		}

		refreshToken, err := tokenManager.Refresh(token)

		if !asSameError(c.expectedErr, err) {
			t.Errorf("Test case: %s. Excepted error to be: %v, but got %v.",
				c.info, c.expectedErr, err)
		}

		if (c.expected && len(refreshToken) == 0) || (!c.expected && len(refreshToken) > 0) {
			t.Errorf("Test Case: %s. Expected new token to be generated: %t", c.info, c.expected)
		}
	}
}
