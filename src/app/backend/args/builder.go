package args

import "net"

var builder = &holderBuilder{holder: Holder}

// Builder structure for holder. Use builder to set value of holder.
// Like java Builder.
type holderBuilder struct {
	holder *holder
}

// SetInsecurePort 'insecure-port' argument of k8sconsole.
func (self *holderBuilder) SetInsecurePort(port int) *holderBuilder {
	self.holder.insecurePort = port
	return self
}

// SetInsecurePort 'port' argument of k8sconsole.
func (self *holderBuilder) SetPort(port int) *holderBuilder {
	self.holder.port = port
	return self
}

// SetInsecureBindAddress 'insecure-bind-address' argument of k8sconsole.
func (self *holderBuilder) SetInsecureBindAddress(ip net.IP) *holderBuilder {
	self.holder.insecureBindAddress = ip
	return self
}

// SetBindAddress 'bind-address' argument of k8sconsole.
func (self *holderBuilder) SetBindAddress(ip net.IP) *holderBuilder {
	self.holder.bindAddress = ip
	return self
}

// SetApiServerHost 'apiserver-host' argument of k8sconsole.
func (self *holderBuilder) SetApiServerHost(apiServerHost string) *holderBuilder {
	self.holder.apiServerHost = apiServerHost
	return self
}

// SetKubeConfigFile 'kubeconfig' argument of k8sconsole.
func (self *holderBuilder) SetKubeConfigFile(kubeConfigFile string) *holderBuilder {
	self.holder.kubeConfigFile = kubeConfigFile
	return self
}

// SetTokenTTL 'token-ttl' argument of k8sconsole.
func (self *holderBuilder) SetTokenTTL(ttl int) *holderBuilder {
	self.holder.tokenTTL = ttl
	return self
}

// SetAuthenticationMode 'authentication-mode' argument of k8sconsole.
func (self *holderBuilder) SetAuthenticationMode(authMode []string) *holderBuilder {
	self.holder.authenticationMode = authMode
	return self
}

// SetDisableSkipButton 'disable-settings-authorizer' argument of k8sconsole.
func (self *holderBuilder) SetDisableSkipButton(disableSkipButton bool) *holderBuilder {
	self.holder.disableSkipButton = disableSkipButton
	return self
}

// SetEnableInsecureLogin 'enable-insecure-login' argument of k8sconsole.
func (self *holderBuilder) SetEnableInsecureLogin(enableInsecureLogin bool) *holderBuilder {
	self.holder.enableInsecureLogin = enableInsecureLogin
	return self
}

// GetHolderBuilder returns singletone instance of argument holder builder.
func GetHolderBuilder() *holderBuilder {
	return builder
}
