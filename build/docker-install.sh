#!/usr/bin/env bash

# install docker version 17.03.02-ce
curl -fsSL https://raw.githubusercontent.com/rancher/install-docker/master/17.03.2.sh | bash -s docker --mirror Aliyun
sudo usermod -aG docker
newgrp docker
docker version

# install the newest docker version
# curl -fsSL https://get.docker.com | bash -s docker --mirror Aliyun