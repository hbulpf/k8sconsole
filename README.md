# k8sconsole
[![Build Status](https://travis-ci.org/wzt3309/k8sconsole.svg?branch=master)](https://travis-ci.org/wzt3309/k8sconsole)
A web ui which extends kubernetes dashboard.

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

### Use vscode
In vscode open file `./example/k8sconsole-api.http` and it will auto install plugin [vscode-restclient](https://github.com/Huachao/vscode-restclient)
for test apis.

#### Use browser
For some get apis you can use browser to directly access them.
