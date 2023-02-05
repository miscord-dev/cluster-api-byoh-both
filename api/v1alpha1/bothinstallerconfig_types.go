/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BothInstallerConfigSpec defines the desired state of BothInstallerConfig
type BothInstallerConfigSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	InstallTemplate   *string `json:"installTemplate"`
	UninstallTemplate *string `json:"uninstallTemplate"`
	Repository        *string `json:"repository"`
	TagTemplate       *string `json:"tagNameTemplate,omitempty"`
}

// BothInstallerConfigStatus defines the observed state of BothInstallerConfig
type BothInstallerConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Ready indicates the InstallationSecret field is ready to be consumed
	// +optional
	Ready bool `json:"ready,omitempty"`

	// InstallationSecret is an optional reference to a generated installation secret by K8sInstallerConfig controller
	// +optional
	InstallationSecret *corev1.ObjectReference `json:"installationSecret,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// BothInstallerConfig is the Schema for the bothinstallerconfigs API
type BothInstallerConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BothInstallerConfigSpec   `json:"spec,omitempty"`
	Status BothInstallerConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BothInstallerConfigList contains a list of BothInstallerConfig
type BothInstallerConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BothInstallerConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BothInstallerConfig{}, &BothInstallerConfigList{})
}

const (
	DefaultTagTemplate = `{{ .K8sVersion }}-
{{- if (contains (toLower .OSImage) "ubuntu") -}}
	ubuntu
{{- else -}}
	unknown
{{- end -}}
	-{{ .Arch }}`

	DefaultInstallTemplate = `
set -euox pipefail

BUNDLE_DOWNLOAD_PATH={{.BundleDownloadPath}}
BUNDLE_ADDR={{.ImageTag}}
IMGPKG_VERSION=.v0.27.0
ARCH={{.Arch}}
BUNDLE_PATH=$BUNDLE_DOWNLOAD_PATH/$BUNDLE_ADDR

if ! command -v imgpkg >>/dev/null; then
	echo "installing imgpkg"	
	
	if command -v wget >>/dev/null; then
		dl_bin="wget -nv -O-"
	elif command -v curl >>/dev/null; then
		dl_bin="curl -s -L"
	else
		echo "installing curl"
		apt-get install -y curl
		dl_bin="curl -s -L"
	fi
	
	$dl_bin github.com/vmware-tanzu/carvel-imgpkg/releases/download/$IMGPKG_VERSION/imgpkg-linux-$ARCH > /tmp/imgpkg
	mv /tmp/imgpkg /usr/local/bin/imgpkg
	chmod +x /usr/local/bin/imgpkg
fi

echo "downloading bundle"
mkdir -p $BUNDLE_PATH
imgpkg pull -r -i $BUNDLE_ADDR -o $BUNDLE_PATH


## disable swap
swapoff -a && sed -ri '/\sswap\s/s/^#?/#/' /etc/fstab

## disable firewall
if command -v ufw >>/dev/null; then
	ufw disable
fi

## load kernal modules
modprobe overlay && modprobe br_netfilter

## adding os configuration
tar -C / -xvf "$BUNDLE_PATH/conf.tar" && sysctl --system 

## installing deb packages
for pkg in cri-tools kubernetes-cni kubectl kubelet kubeadm; do
	dpkg --install "$BUNDLE_PATH/$pkg.deb" && apt-mark hold $pkg
done

## intalling containerd
tar -C / -xvf "$BUNDLE_PATH/containerd.tar"

## starting containerd service
systemctl daemon-reload && systemctl enable containerd && systemctl start containerd
`

	DefaultUninstallTemplate = `
set -euox pipefail

BUNDLE_DOWNLOAD_PATH={{.BundleDownloadPath}}
BUNDLE_ADDR={{.ImageTag}}
BUNDLE_PATH=$BUNDLE_DOWNLOAD_PATH/$BUNDLE_ADDR

## disabling containerd service
systemctl stop containerd && systemctl disable containerd && systemctl daemon-reload

## removing containerd configurations and cni plugins
rm -rf /opt/cni/ && rm -rf /opt/containerd/ &&  tar tf "$BUNDLE_PATH/containerd.tar" | xargs -n 1 echo '/' | sed 's/ //g'  | grep -e '[^/]$' | xargs rm -f

## removing deb packages
for pkg in kubeadm kubelet kubectl kubernetes-cni cri-tools; do
	dpkg --purge $pkg
done

## removing os configuration
tar tf "$BUNDLE_PATH/conf.tar" | xargs -n 1 echo '/' | sed 's/ //g' | grep -e "[^/]$" | xargs rm -f

## remove kernal modules
modprobe -rq overlay && modprobe -r br_netfilter

## enable firewall
if command -v ufw >>/dev/null; then
	ufw enable
fi

## enable swap
swapon -a && sed -ri '/\sswap\s/s/^#?//' /etc/fstab

rm -rf $BUNDLE_PATH
`
)
