package controllers

/*
Copyright 2023 VMware, Inc.

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

import (
	"context"

	infrav1 "github.com/vmware-tanzu/cluster-api-provider-bringyourownhost/apis/infrastructure/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetOwnerByoMachine returns the ByoMachine object owning the current resource.
func GetOwnerByoMachine(ctx context.Context, c client.Client, obj *metav1.ObjectMeta) (*infrav1.ByoMachine, error) {
	for _, ref := range obj.OwnerReferences {
		gv, err := schema.ParseGroupVersion(ref.APIVersion)
		if err != nil {
			return nil, err
		}
		if ref.Kind == "ByoMachine" && gv.Group == infrav1.GroupVersion.Group {
			return GetByoMachineByName(ctx, c, obj.Namespace, ref.Name)
		}
	}
	return nil, nil
}

// GetByoMachineByName finds and return a ByoMachine object using the specified params.
func GetByoMachineByName(ctx context.Context, c client.Client, namespace, name string) (*infrav1.ByoMachine, error) {
	m := &infrav1.ByoMachine{}
	key := client.ObjectKey{Name: name, Namespace: namespace}
	if err := c.Get(ctx, key, m); err != nil {
		return nil, err
	}
	return m, nil
}
