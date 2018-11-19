package api

import "strings"

// ToAuthenticationModes convert array of authentication mode strings to valid Authentication type.
func ToAuthenticationModes(modes []string) AuthenticationModes {
	res := AuthenticationModes{}
	modesMap := map[string]bool{}

	for _, mode := range []AuthenticationMode{Basic, Token} {
		modesMap[mode.String()] = true
	}

	for _, mode := range modes {
		if _, exists := modesMap[mode]; exists {
			res.Add(AuthenticationMode(mode))
		}
	}

	return res
}

// ShouldRejectRequest returns true if url contains name and namespace of resource that should be filtered out
func ShouldRejectRquest(url string) bool {
	return strings.Contains(url, EncrytionKeyHolderName) && strings.Contains(url, EncryptionKeyHolderNamespace)
}
