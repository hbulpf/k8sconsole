package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/spf13/pflag"
	"github.com/wzt3309/k8sconsole/src/app/backend/args"
	"github.com/wzt3309/k8sconsole/src/app/backend/auth"
	authApi "github.com/wzt3309/k8sconsole/src/app/backend/auth/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/auth/jwe"
	"github.com/wzt3309/k8sconsole/src/app/backend/client"
	clientApi "github.com/wzt3309/k8sconsole/src/app/backend/client/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/handler"
	"net"
	"net/http"
	"time"
)

var (
	argInsecurePort        = pflag.Int("insecure-port", 9090, "The port to listen to for incoming HTTP requests.")
	argPort                = pflag.Int("port", 9443, "The port to listen to for incoming HTTPS requests.")
	argInsecureBindAddress = pflag.IP("insecure-bind-address", net.IPv4(127, 0, 0, 1), "The IP address on which to serve the --insecure-port (set to 0.0.0.0 for all interfaces).")
	argBindAddress         = pflag.IP("bind-address", net.IPv4(0, 0, 0, 0), "The IP address on which to serve the --port (set to 0.0.0.0 for all interfaces).")

	argApiServerHost = pflag.String("apiserver-host", "", "The address of kubernetes apiserver to connect to in the format of protocol://address:port, e.g. http://localhost:8080."+
		"If not specified, the assumption is that k8sconsole binary runs inside a kubernetes cluster and local discovery is attempted")
	argKubeConfigFile = pflag.String("kubeconfig", "", "Path to kubeconfig file with kubernetes cluster authorization and location information."+
		"If not specified, the assumption is that k8sconsole binary runs inside a kubernetes cluster and local discovery is attempted")
	argTokenTTL           = pflag.Int("token-ttl", int(authApi.DefaultTokenTTL), "Expiration time (in seconds) of JWE tokens generated by k8sconsole. Default: 15 min. 0 - never expires.")
	argAuthenticationMode = pflag.StringSlice("authentication-mode", []string{authApi.Token.String()}, "Enables authentication options that will be reflected on login screen. Supported values: basic, token. Default: token."+
		"Note that basic option should only be used if apiserver has '--authorization-mode=ABAC' and '--basic-auth-file' flags set.")
	argDisableSkip         = pflag.Bool("disable-skip", false, "When enabled, the skip button on the login page will not be shown. Default: false.")
	argEnableInsecureLogin = pflag.Bool("enable-insecure-login", false, "When enabled, k8sconsole login view will also be shown when k8sconsole is not served over HTTPS. Default: false.")
)

func initArgHolder() {
	builder := args.GetHolderBuilder()
	builder.SetInsecurePort(*argInsecurePort)
	builder.SetPort(*argPort)
	builder.SetInsecureBindAddress(*argInsecureBindAddress)
	builder.SetBindAddress(*argBindAddress)
	builder.SetApiServerHost(*argApiServerHost)
	builder.SetKubeConfigFile(*argKubeConfigFile)
	builder.SetTokenTTL(*argTokenTTL)
	builder.SetAuthenticationMode(*argAuthenticationMode)
	builder.SetDisableSkipButton(*argDisableSkip)
	builder.SetEnableInsecureLogin(*argEnableInsecureLogin)
}

func initAuthManager(clientManager clientApi.ClientManager) authApi.AuthManager {
	// Init tokenManager
	keyHolder := jwe.NewRSAKeyHolder()
	tokenManager := jwe.NewJWETokenManager(keyHolder)
	tokenTTL := time.Duration(args.Holder.GetTokenTTL())
	if tokenTTL != authApi.DefaultTokenTTL {
		tokenManager.SetTokenTTL(tokenTTL)
	}

	// Set tokenManager for clientManager
	clientManager.SetTokenManager(tokenManager)

	// Convert auth modes string slice to AuthenticationModes
	authModes := authApi.ToAuthenticationModes(args.Holder.GetAuthenticationMode())
	if len(authModes) == 0 {
		authModes.Add(authApi.Token)
	}

	authenticationSkippable := !args.Holder.GetDisableSkipButton()

	return auth.NewAuthManager(clientManager, tokenManager, authModes, authenticationSkippable)
}

func handleFatalInitError(err error) {
	glog.Fatalf("Error while initializing connection to Kubernetes apiserver." +
		" This most likely means that cluster is misconfigured(e.g., it has" +
		" invalid apiserver certs or service account's configuration) or the" +
		" --apiserver-host param points to a apiserver that does not exist.", err)
}

func main() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	flag.CommandLine.Parse(make([]string, 0))	// Init for glog calls

	// Initializes k8sconsole arguments holder so we can read them in other package
	initArgHolder()

	if args.Holder.GetApiServerHost() != "" {
		glog.Infof("Using apisever-host location: %s", args.Holder.GetApiServerHost())
	}
	if args.Holder.GetKubeConfigFile() != "" {
		glog.Infof("Using kuberconfig file: %s", args.Holder.GetKubeConfigFile())
	}

	// Initialize clientManager
	clientManager := client.NewClientManager(args.Holder.GetKubeConfigFile(), args.Holder.GetApiServerHost())
	versionInfo, err := clientManager.InsecureClient().Discovery().ServerVersion()
	if err != nil {
		handleFatalInitError(err)
	}

	glog.Infof("Successful initial request to the apiserver, version: %s", versionInfo.String())

	// Initialize auth manager
	authManager := initAuthManager(clientManager)

	// Create apiHandler
	apiHandler, err := handler.CreateHTTPAPIHandler(clientManager, authManager)
	if err != nil {
		handleFatalInitError(err)
	}

	http.Handle("/api/", apiHandler)
	http.Handle("/api/sockjs/", handler.CreateAttachHandler("/api/sockjs"))

	// TODO(wzt3309) listening on https
	glog.Infof("Serving insecurely on HTTP port: %d", args.Holder.GetInsecurePort())
	addr := fmt.Sprintf("%s:%d", args.Holder.GetInsecureBindAddress(), args.Holder.GetInsecurePort())
	go func() {glog.Fatal(http.ListenAndServe(addr, nil))}()
	select {}
}
