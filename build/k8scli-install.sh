#!/usr/bin/env bash
# Use kubernetes version 1.10.0

# Debian / Ubuntu
sudo apt-get update && apt-get install -y apt-transport-https
curl https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | apt-key add -
cat <<EOF >/etc/apt/sources.list.d/kubernetes.list
deb https://mirrors.aliyun.com/kubernetes/apt/ kubernetes-xenial main
EOF

sudo apt-get update
## specified k8s version
sudo apt-get install -y kubelet=1.10.0-00 kubeadm=1.10.0-00 kubectl=1.10.4-00

# CentOS / RHEL / Fedora
#cat <<EOF > /etc/yum.repos.d/kubernetes.repo
#[kubernetes]
#name=Kubernetes
#baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64/
#enabled=1
#gpgcheck=1
#repo_gpgcheck=1
#gpgkey=https://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg https://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
#EOF

#setenforce 0
## need to set version
#yum install -y kubelet kubeadm kubectl
#systemctl enable kubelet && systemctl start kubelet