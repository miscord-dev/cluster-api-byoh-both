/*
Copyright 2023 VMware, Inc., miscord-dev

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

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	infrav1 "github.com/vmware-tanzu/cluster-api-provider-bringyourownhost/apis/infrastructure/v1beta1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	byohmiscordwinv1alpha1 "github.com/miscord-dev/cluster-api-byoh-both/api/v1alpha1"
	"github.com/miscord-dev/cluster-api-byoh-both/pkg/installer"
)

// BothInstallerConfigReconciler reconciles a BothInstallerConfig object
type BothInstallerConfigReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	Installer installer.Installer
}

// bothInstallerConfigScope defines a scope defined around a BothInstallerConfig and its ByoMachine
type bothInstallerConfigScope struct {
	Client     client.Client
	Logger     logr.Logger
	Cluster    *clusterv1.Cluster
	ByoMachine *infrav1.ByoMachine
	Config     *byohmiscordwinv1alpha1.BothInstallerConfig
}

//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=bothinstallerconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=bothinstallerconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=bothinstallerconfigs/finalizers,verbs=update
//+kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=byomachines,verbs=get;list;watch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=byomachines/status,verbs=get
//+kubebuilder:rbac:groups="",resources=secrets;events,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the BothInstallerConfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *BothInstallerConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconcile request received")

	// Fetch the BothInstallerConfig instance
	config := &byohmiscordwinv1alpha1.BothInstallerConfig{}
	err := r.Client.Get(ctx, req.NamespacedName, config)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get BothInstallerConfig")
		return ctrl.Result{}, err
	}

	// Create the BothInstallerConfig scope
	scope := &bothInstallerConfigScope{
		Client: r.Client,
		Logger: logger.WithValues("bothInstallerConfig", config.Name),
		Config: config,
	}

	// Fetch the ByoMachine
	byoMachine, err := GetOwnerByoMachine(ctx, r.Client, &config.ObjectMeta)
	if err != nil && !apierrors.IsNotFound(err) {
		logger.Error(err, "failed to get Owner ByoMachine")
		return ctrl.Result{}, err
	}

	helper, err := patch.NewHelper(config, r.Client)
	if err != nil {
		logger.Error(err, "unable to create helper")
		return ctrl.Result{}, err
	}
	defer func() {
		if err = helper.Patch(ctx, config); err != nil && reterr == nil {
			logger.Error(err, "failed to patch BothInstallerConfig")
			reterr = err
		}
	}()

	// Handle deleted BothInstallerConfig
	if !config.ObjectMeta.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, nil
	}

	if byoMachine == nil {
		logger.Info("Waiting for ByoMachine Controller to set OwnerRef on InstallerConfig")
		return ctrl.Result{}, nil
	}
	scope.ByoMachine = byoMachine
	logger = logger.WithValues("byoMachine", byoMachine.Name, "namespace", byoMachine.Namespace)
	logger.Info("byoMachine found")

	// Fetch the Cluster
	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, byoMachine.ObjectMeta)
	if err != nil {
		logger.Error(err, "ByoMachine owner Machine is missing cluster label or cluster does not exist")
		return ctrl.Result{}, err
	}
	logger = logger.WithValues("cluster", cluster.Name)
	scope.Cluster = cluster
	scope.Logger = logger

	if annotations.IsPaused(cluster, config) {
		logger.Info("Reconciliation is paused for this object")
		return ctrl.Result{}, nil
	}

	switch {
	// waiting for ByoMachine to updating it's ByoHostReady condition to false for reason InstallationSecretNotAvailableReason
	case conditions.GetReason(byoMachine, infrav1.BYOHostReady) != infrav1.InstallationSecretNotAvailableReason:
		logger.Info("ByoMachine is not waiting for InstallationSecret", "reason", conditions.GetReason(byoMachine, infrav1.BYOHostReady))
		return ctrl.Result{}, nil
	// Status is ready means a config has been generated.
	case config.Status.Ready:
		logger.Info("BothInstallerConfig is ready")
		return ctrl.Result{}, nil
	}

	return r.reconcileNormal(ctx, scope)
}

func (r *BothInstallerConfigReconciler) reconcileNormal(ctx context.Context, scope *bothInstallerConfigScope) (reconcile.Result, error) {
	logger := scope.Logger
	logger.Info("Reconciling BothInstallerConfig")

	configSpec := scope.Config.Spec
	k8sVersion := scope.Config.GetAnnotations()[infrav1.K8sVersionAnnotation]

	installTemplate := stringWithDefault(configSpec.InstallTemplate, byohmiscordwinv1alpha1.DefaultInstallTemplate)
	uninstallTemplate := stringWithDefault(configSpec.UninstallTemplate, byohmiscordwinv1alpha1.DefaultUninstallTemplate)
	repo := stringWithDefault(configSpec.TagTemplate, "")
	tagTemplate := stringWithDefault(configSpec.TagTemplate, byohmiscordwinv1alpha1.DefaultTagTemplate)

	install, uninstall, err := r.Installer.Generate(installer.InstallerConfig{
		NodeInfo: installer.NodeInfo{
			OS:         scope.ByoMachine.Status.HostInfo.OSName,
			OSImage:    scope.ByoMachine.Status.HostInfo.OSImage,
			Arch:       scope.ByoMachine.Status.HostInfo.Architecture,
			K8sVersion: k8sVersion,
		},
		InstallTemplate:   installTemplate,
		UninstallTemplate: uninstallTemplate,
		Repository:        repo,
		TagTemplate:       tagTemplate,
	})

	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to generate install/uninstall script: %w", err)
	}

	// creating installation secret
	if err := r.storeInstallationData(ctx, scope, install, uninstall); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// storeInstallationData creates a new secret with the install and unstall data passed in as input,
// sets the reference in the configuration status and ready to true.
func (r *BothInstallerConfigReconciler) storeInstallationData(ctx context.Context, scope *bothInstallerConfigScope, install, uninstall string) error {
	logger := scope.Logger
	logger.Info("creating installation secret")

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      scope.Config.Name,
			Namespace: scope.Config.Namespace,
			Labels: map[string]string{
				clusterv1.ClusterLabelName: scope.Cluster.Name,
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: infrav1.GroupVersion.String(),
					Kind:       scope.Config.Kind,
					Name:       scope.Config.Name,
					UID:        scope.Config.UID,
					Controller: pointer.Bool(true),
				},
			},
		},
		Data: map[string][]byte{
			"install":   []byte(install),
			"uninstall": []byte(uninstall),
		},
		Type: clusterv1.ClusterSecretType,
	}

	// as secret creation and scope.Config status patch are not atomic operations
	// it is possible that secret creation happens but the config.Status patches are not applied
	if err := r.Client.Create(ctx, secret); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrapf(err, "failed to create installation secret for BothInstallerConfig %s/%s", scope.Config.Namespace, scope.Config.Name)
		}
		logger.Info("installation secret for BothInstallerConfig already exists, updating", "secret", secret.Name, "BothInstallerConfig", scope.Config.Name)
		if err := r.Client.Update(ctx, secret); err != nil {
			return errors.Wrapf(err, "failed to update installation secret for BothInstallerConfig %s/%s", scope.Config.Namespace, scope.Config.Name)
		}
	}
	scope.Config.Status.InstallationSecret = &corev1.ObjectReference{
		Kind:      secret.Kind,
		Namespace: secret.Namespace,
		Name:      secret.Name,
	}
	scope.Config.Status.Ready = true
	logger.Info("created installation secret")
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BothInstallerConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&byohmiscordwinv1alpha1.BothInstallerConfig{}).
		Watches(
			&source.Kind{Type: &infrav1.ByoMachine{}},
			handler.EnqueueRequestsFromMapFunc(r.ByoMachineToBothInstallerConfigMapFunc),
		).
		Complete(r)
}

// BothInstallerConfigReconciler is a handler.ToRequestsFunc to be used to enqeue
// request for reconciliation of BothInstallerConfig.
func (r *BothInstallerConfigReconciler) ByoMachineToBothInstallerConfigMapFunc(o client.Object) []ctrl.Request {
	ctx := context.TODO()
	logger := log.FromContext(ctx)

	m, ok := o.(*infrav1.ByoMachine)
	if !ok {
		panic(fmt.Sprintf("Expected a ByoMachine but got a %T", o))
	}
	m.GetObjectKind().SetGroupVersionKind(infrav1.GroupVersion.WithKind("ByoMachine"))

	result := []ctrl.Request{}
	if m.Spec.InstallerRef != nil && m.Spec.InstallerRef.GroupVersionKind() == byohmiscordwinv1alpha1.GroupVersion.WithKind("BothInstallerConfigTemplate") {
		configList := &byohmiscordwinv1alpha1.BothInstallerConfigList{}
		if err := r.Client.List(ctx, configList, client.InNamespace(m.Namespace)); err != nil {
			logger.Error(err, "failed to list BothInstallerConfig")
			return result
		}
		for idx := range configList.Items {
			config := &configList.Items[idx]
			if hasOwnerReferenceFrom(config, m) {
				name := client.ObjectKey{Namespace: config.Namespace, Name: config.Name}
				result = append(result, ctrl.Request{NamespacedName: name})
			}
		}
	}
	return result
}

// hasOwnerReferenceFrom will check if object have owner reference of the given owner
func hasOwnerReferenceFrom(obj, owner client.Object) bool {
	for _, o := range obj.GetOwnerReferences() {
		if o.Kind == owner.GetObjectKind().GroupVersionKind().Kind && o.Name == owner.GetName() {
			return true
		}
	}
	return false
}

func stringWithDefault(ptr *string, defaultValue string) string {
	if ptr == nil {
		return defaultValue
	}

	return *ptr
}
