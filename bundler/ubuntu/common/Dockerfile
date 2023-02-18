# Copyright 2021 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# Downloads bundle ingredients : containerd as tar, kubelet, kubeadm, kubectl as Debian packages
#
# Usage:
# 1. Mount a host path as /ingredients
# 2. Run the image
#

ARG BASE_IMAGE=ubuntu:20.04
ARG INGREDIENTS_PATH=/ingredients
FROM $BASE_IMAGE as build

# Override to download other version
ARG CONTAINERD_VERSION=1.6.8
ARG KUBERNETES_VERSION=1.24.2
ARG KUBERNETES_CNI_VERSION=1.2.0
ARG KUBERNETES_CRI_VERSION=1.26.0
ARG ARCH=amd64
ARG INGREDIENTS_PATH

WORKDIR /bundle-builder
COPY download.sh .
RUN chmod a+x download.sh
WORKDIR $INGREDIENTS_PATH
RUN /bundle-builder/download.sh

FROM $BASE_IMAGE as bundler

ARG INGREDIENTS_PATH

WORKDIR /bundler

# Mandatory
COPY --from=build $INGREDIENTS_PATH/containerd.tar ./containerd.tar
COPY --from=build $INGREDIENTS_PATH/kubeadm.deb ./kubeadm.deb
COPY --from=build $INGREDIENTS_PATH/kubelet.deb ./kubelet.deb
COPY --from=build $INGREDIENTS_PATH/kubectl.deb ./kubectl.deb
COPY --from=build $INGREDIENTS_PATH/cri-tools.deb cri-tools.deb
COPY --from=build $INGREDIENTS_PATH/kubernetes-cni.deb kubernetes-cni.deb

COPY ./config /config
RUN (cd /config && tar -cvf conf.tar *) && \
    cp /config/conf.tar .
