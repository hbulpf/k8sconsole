# K8sConsole / [中文版](./README_cn.md)
> A web ui which extends kubernetes dashboard

## Quick Start
1. [Get binary file](https://github.com/wzt3309/k8sconsole/releases)
2. Start [Kubernetes](https://github.com/kubernetes/kubernetes) or [MiniKube](https://github.com/kubernetes/minikube)
3. Start [K8sConsole](https://github.com/hbulpf/k8sconsole)

```
# When you use macOS
./k8sconsole-darwin-amd64 \
	--apiserver-host {api url of Kubernetes cluster} \
	--insecure-port {a local port not in use, default is 9090 }
./k8sconsole-darwin-amd64 --apiserver-host http://127.0.0.1:61987 --insecure-port 63453

# When you use linux
./k8sconsole-linux-amd64 --apiserver-host http://127.0.0.1:61987 --insecure-port 63453

# When you use windows
k8sconsole-windows-amd64.exe --apiserver-host http://127.0.0.1:61987 --insecure-port 63453
```

# k8sconsole API
## Online
See online api docs in [k8sconsole-go](https://app.swaggerhub.com/apis/ztwang/k8sconsole-go/0.0.1).

## Local test
### Step 1. Start a kubernetes cluster
We use minikube to start a local kubernetes cluster v1.10.0.
> Required. You need to install docker before `./build/docker-install.sh` (This will install docker 17.03.02-ce)

`./build/minikube.sh`

### Step 2. Start backend
> Install backend from [releases](https://github.com/wzt3309/k8sconsole/releases)

`./k8sconsole --apiserver-host=http://localhost:8080 --logtostderr`

The k8sconsole will listen on default insecure port 9090.

You can use `./k8sconsole --help` for more information.

### Step 3. Access rest api
#### Use Jetbrain(IDEA, WebStorm, ...)
In Jetbrain IDE open file `./example/k8sconsole-api.http`, you can use ide's "HTTP Client Tool"
to test rest apis.

#### Use curl
like `curl -X GET "http://localhost:9090/api/v1/node?filterBy=name%2Cminikube&sortBy=d%2Cname&itemsPerPage=1&page=1" -H "accept: application/json"`

#### Use vscode
In vscode open file `./example/k8sconsole-api.http` and it will auto install plugin [vscode-restclient](https://github.com/Huachao/vscode-restclient)
for test apis.

#### Use browser
For some get apis you can use browser to directly access them.
