package args

import "net"

var Holder = &holder{}

// Argument holder structure. It holds all arguments values passed to k8sconsole
// Using singleton mode
type holder struct {
	insecurePort        int
	port                int
	insecureBindAddress net.IP
	bindAddress         net.IP

	apiServerHost  string
	kubeConfigFile string

	tokenTTL           int
	authenticationMode []string
	disableSkipButton  bool

	enableInsecureLogin bool
}

// GetInsecurePort 'insecure-port' argument of k8sconsole.
func (self *holder) GetInsecurePort() int {
	return self.insecurePort
}

// GetPort 'port' argument of k8sconsole.
func (self *holder) GetPort() int {
	return self.port
}

// GetInsecureBindAddress 'insecure-bind-address' argument of k8sconsole.
func (self *holder) GetInsecureBindAddress() net.IP {
	return self.insecureBindAddress
}

// GetBindAddress 'bind-address' argument of k8sconsole.
func (self *holder) GetBindAddress() net.IP {
	return self.bindAddress
}

// GetApiServerHost 'apiserver-host' argument of k8sconsole.
func (self *holder) GetApiServerHost() string {
	return self.apiServerHost
}

// GetKubeConfigFile 'kubeconfig' argument of k8sconsole.
func (self *holder) GetKubeConfigFile() string {
	return self.kubeConfigFile
}

// GetTokenTTL 'token-ttl' argument of k8sconsole.
func (self *holder) GetTokenTTL() int {
	return self.tokenTTL
}

// GetAuthenticationMode 'authentication-mode' argument of k8sconsole.
func (self *holder) GetAuthenticationMode() []string {
	return self.authenticationMode
}

// GetDisableSkipButton 'disable-settings-authorizer' argument of k8sconsole.
func (self *holder) GetDisableSkipButton() bool {
	return self.disableSkipButton
}

// GetEnableInsecureLogin 'enable-insecure-login' argument of k8sconsole.
func (self *holder) GetEnableInsecureLogin() bool {
	return self.enableInsecureLogin
}
