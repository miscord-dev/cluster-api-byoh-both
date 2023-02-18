#!/bin/bash

# Copyright 2021 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -eux

echo  Update the apt package index and install packages needed to use the Kubernetes apt repository
apt-get update
apt-get install --no-install-recommends -y apt-transport-https ca-certificates curl

echo Download containerd
curl -LOJR https://github.com/containerd/containerd/releases/download/v${CONTAINERD_VERSION}/cri-containerd-cni-${CONTAINERD_VERSION}-linux-amd64.tar.gz

echo Download the Google Cloud public signing key
curl -fsSLo /usr/share/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg

echo Add the Kubernetes apt repository
echo "deb [signed-by=/usr/share/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | tee /etc/apt/sources.list.d/kubernetes.list

echo Update apt package index, install kubelet, kubeadm and kubectl
apt-get update || apt-get update
apt-get download {kubelet,kubeadm,kubectl}:$ARCH=$KUBERNETES_VERSION-00
apt-get download kubernetes-cni:$ARCH=$KUBERNETES_CNI_VERSION-00
apt-get download cri-tools:$ARCH=$KUBERNETES_CRI_VERSION-00

mv *containerd* containerd.tar
mv *kubeadm*.deb kubeadm.deb
mv *kubelet*.deb kubelet.deb
mv *kubectl*.deb kubectl.deb
mv *cri-tools*.deb cri-tools.deb
mv *kubernetes-cni*.deb kubernetes-cni.deb
