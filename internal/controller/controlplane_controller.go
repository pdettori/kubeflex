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

package controller

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	clog "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	tenancyv1alpha1 "mcc.ibm.org/kubeflex/api/v1alpha1"
)

// ControlPlaneReconciler reconciles a ControlPlane object
type ControlPlaneReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tenancy.mcc.ibm.org,resources=controlplanes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tenancy.mcc.ibm.org,resources=controlplanes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tenancy.mcc.ibm.org,resources=controlplanes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ControlPlane object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *ControlPlaneReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := clog.FromContext(ctx)

	log.Info("Got ControlPlane event!")

	// Fetch the hostedControlPlane instance
	hostedControlPlane := &tenancyv1alpha1.ControlPlane{}
	err := r.Client.Get(ctx, req.NamespacedName, hostedControlPlane)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	hcp := hostedControlPlane.DeepCopy()

	if err = r.ReconcileNamespace(ctx, hcp.Name, &v1.OwnerReference{
		Kind: hcp.Kind, APIVersion: hcp.APIVersion, Name: hcp.Name, UID: hcp.UID}); err != nil {
		return ctrl.Result{}, err
	}

	if err = r.ReconcileAPIServerDeployment(ctx, hcp.Name, &v1.OwnerReference{
		Kind: hcp.Kind, APIVersion: hcp.APIVersion, Name: hcp.Name, UID: hcp.UID}); err != nil {
		return ctrl.Result{}, err
	}

	log.Info("Hosted control plane", "my-name-is", hcp.Name)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ControlPlaneReconciler) SetupWithManager(mgr ctrl.Manager) error {
	b := ctrl.NewControllerManagedBy(mgr).
		For(&tenancyv1alpha1.ControlPlane{}).
		WithOptions(controller.Options{
			RateLimiter: workqueue.NewItemExponentialFailureRateLimiter(1*time.Second, 10*time.Second),
		})
	for _, handler := range r.eventHandlers() {
		b.Watches(&source.Kind{Type: handler.obj}, handler.handler)
	}
	if _, err := b.Build(r); err != nil {
		return fmt.Errorf("failed setting up with a controller manager %w", err)
	}
	return nil
}

type eventHandler struct {
	obj     client.Object
	handler handler.EventHandler
}

func (r *ControlPlaneReconciler) eventHandlers() []eventHandler {
	handlers := []eventHandler{
		{obj: &corev1.Service{}, handler: &handler.EnqueueRequestForOwner{OwnerType: &tenancyv1alpha1.ControlPlane{}}},
		{obj: &appsv1.Deployment{}, handler: &handler.EnqueueRequestForOwner{OwnerType: &tenancyv1alpha1.ControlPlane{}}},
		{obj: &appsv1.StatefulSet{}, handler: &handler.EnqueueRequestForOwner{OwnerType: &tenancyv1alpha1.ControlPlane{}}},
		{obj: &corev1.Secret{}, handler: &handler.EnqueueRequestForOwner{OwnerType: &tenancyv1alpha1.ControlPlane{}}},
		{obj: &corev1.ConfigMap{}, handler: &handler.EnqueueRequestForOwner{OwnerType: &tenancyv1alpha1.ControlPlane{}}},
		{obj: &corev1.ServiceAccount{}, handler: &handler.EnqueueRequestForOwner{OwnerType: &tenancyv1alpha1.ControlPlane{}}},
	}
	return handlers
}