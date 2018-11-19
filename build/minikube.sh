#!/usr/bin/env bash

BASE_DIR=.tools/k8s
MINIKUBE_VERSION=v0.28.0
MINIKUBE_K8S_VERSION=v1.10.0
MINIKUBE_BIN=${BASE_DIR}/minikube-${MINIKUBE_VERSION}

echo "Making sure that ${BASE_DIR} directory exists"
mkdir -p ${BASE_DIR}

echo "Downloading minikube ${MINIKUBE_VERSION} if it is not cached"
curl -Lo minikube http://kubernetes.oss-cn-hangzhou.aliyuncs.com/minikube/releases/${MINIKUBE_VERSION}/minikube-linux-amd64 \
&& chmod +x minikube && sudo mv minikube ${MINIKUBE_BIN}

echo "Making sure that kubeconfig file exists and will be used by Dashboard"
mkdir -p $HOME/.kube
touch $HOME/.kube/config

echo "Starting minikube"
export MINIKUBE_WANTUPDATENOTIFICATION=false
export MINIKUBE_WANTREPORTERRORPROMPT=false
export MINIKUBE_HOME=${HOME}
export CHANGE_MINIKUBE_NONE_USER=true
sudo -E ${MINIKUBE_BIN} start --registry-mirror=https://registry.docker-cn.com \
--vm-driver=none \
--kubernetes-version ${MINIKUBE_K8S_VERSION}

