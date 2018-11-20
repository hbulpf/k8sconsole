# k8s控制台 / [Englis](./README.md)
> 基于容器的弹性大数据平台的k8s控制台

## 如何部署
1. [下载二进制文件](https://github.com/wzt3309/k8sconsole/releases)
2. 启动 [Kubernetes](https://github.com/kubernetes/kubernetes) 集群或者本地启动 [MiniKube](https://github.com/kubernetes/minikube)
3. Start [K8sConsole](https://github.com/hbulpf/k8sconsole)
```
./k8sconsole-darwin-amd64 \
	--apiserver-host {Kubernetes集群的api url} \
	--insecure-port {未被占用的本地端口}
./k8sconsole-darwin-amd64 --apiserver-host http://127.0.0.1:61987 --insecure-port 63453
```