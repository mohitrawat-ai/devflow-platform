/*
Copyright 2025.

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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	appsv1 "github.com/mohitrawat-ai/devflow-platform/operators/api/v1"
)

// AppDeploymentReconciler reconciles a AppDeployment object
type AppDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.devflow.io,resources=appdeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.devflow.io,resources=appdeployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.devflow.io,resources=appdeployments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AppDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *AppDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log = logf.FromContext(ctx)

	// TODO(user): your logic here
	log.Info("Reconciling AppDeployment : ", req.NamespacedName)

	// Fetch the AppDeployment instance
	var appDeployment appsv1.AppDeployment
	if err := r.Get(ctx, req.NamespacedName, &appDeployment); err != nil {
		log.Error(err, "unable to fetch AppDeployment")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	
	// ADD THIS LOG:
	log.Info("Found AppDeployment", "name", appDeployment.Name, "image", appDeployment.Spec.Image)
	

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AppDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.AppDeployment{}).
		Named("appdeployment").
		Complete(r)
}
