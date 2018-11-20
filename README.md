# K8sConsole / [中文版](./README_cn.md)
> A web ui which extends kubernetes dashboard

## How To Start 
1. [Get binary file](https://github.com/wzt3309/k8sconsole/releases)
2. Start [Kubernetes](https://github.com/kubernetes/kubernetes) or [MiniKube](https://github.com/kubernetes/minikube)
3. Start [K8sConsole](https://github.com/hbulpf/k8sconsole)
```
./k8sconsole-darwin-amd64 \
	--apiserver-host {api url of Kubernetes cluster} \
	--insecure-port {a local port not in use }
./k8sconsole-darwin-amd64 --apiserver-host http://127.0.0.1:61987 --insecure-port 63453
```