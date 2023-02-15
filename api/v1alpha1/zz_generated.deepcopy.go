//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BothInstallerConfig) DeepCopyInto(out *BothInstallerConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BothInstallerConfig.
func (in *BothInstallerConfig) DeepCopy() *BothInstallerConfig {
	if in == nil {
		return nil
	}
	out := new(BothInstallerConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BothInstallerConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BothInstallerConfigList) DeepCopyInto(out *BothInstallerConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]BothInstallerConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BothInstallerConfigList.
func (in *BothInstallerConfigList) DeepCopy() *BothInstallerConfigList {
	if in == nil {
		return nil
	}
	out := new(BothInstallerConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BothInstallerConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BothInstallerConfigSpec) DeepCopyInto(out *BothInstallerConfigSpec) {
	*out = *in
	if in.InstallTemplate != nil {
		in, out := &in.InstallTemplate, &out.InstallTemplate
		*out = new(string)
		**out = **in
	}
	if in.UninstallTemplate != nil {
		in, out := &in.UninstallTemplate, &out.UninstallTemplate
		*out = new(string)
		**out = **in
	}
	if in.Repository != nil {
		in, out := &in.Repository, &out.Repository
		*out = new(string)
		**out = **in
	}
	if in.TagTemplate != nil {
		in, out := &in.TagTemplate, &out.TagTemplate
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BothInstallerConfigSpec.
func (in *BothInstallerConfigSpec) DeepCopy() *BothInstallerConfigSpec {
	if in == nil {
		return nil
	}
	out := new(BothInstallerConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BothInstallerConfigStatus) DeepCopyInto(out *BothInstallerConfigStatus) {
	*out = *in
	if in.InstallationSecret != nil {
		in, out := &in.InstallationSecret, &out.InstallationSecret
		*out = new(v1.ObjectReference)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BothInstallerConfigStatus.
func (in *BothInstallerConfigStatus) DeepCopy() *BothInstallerConfigStatus {
	if in == nil {
		return nil
	}
	out := new(BothInstallerConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BothInstallerConfigTemplate) DeepCopyInto(out *BothInstallerConfigTemplate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BothInstallerConfigTemplate.
func (in *BothInstallerConfigTemplate) DeepCopy() *BothInstallerConfigTemplate {
	if in == nil {
		return nil
	}
	out := new(BothInstallerConfigTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BothInstallerConfigTemplate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BothInstallerConfigTemplateList) DeepCopyInto(out *BothInstallerConfigTemplateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]BothInstallerConfigTemplate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BothInstallerConfigTemplateList.
func (in *BothInstallerConfigTemplateList) DeepCopy() *BothInstallerConfigTemplateList {
	if in == nil {
		return nil
	}
	out := new(BothInstallerConfigTemplateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BothInstallerConfigTemplateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BothInstallerConfigTemplateResource) DeepCopyInto(out *BothInstallerConfigTemplateResource) {
	*out = *in
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BothInstallerConfigTemplateResource.
func (in *BothInstallerConfigTemplateResource) DeepCopy() *BothInstallerConfigTemplateResource {
	if in == nil {
		return nil
	}
	out := new(BothInstallerConfigTemplateResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BothInstallerConfigTemplateSpec) DeepCopyInto(out *BothInstallerConfigTemplateSpec) {
	*out = *in
	in.Template.DeepCopyInto(&out.Template)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BothInstallerConfigTemplateSpec.
func (in *BothInstallerConfigTemplateSpec) DeepCopy() *BothInstallerConfigTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(BothInstallerConfigTemplateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BothInstallerConfigTemplateStatus) DeepCopyInto(out *BothInstallerConfigTemplateStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BothInstallerConfigTemplateStatus.
func (in *BothInstallerConfigTemplateStatus) DeepCopy() *BothInstallerConfigTemplateStatus {
	if in == nil {
		return nil
	}
	out := new(BothInstallerConfigTemplateStatus)
	in.DeepCopyInto(out)
	return out
}
