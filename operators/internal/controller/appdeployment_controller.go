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

	appsv1alpha1 "github.com/mohitrawat-ai/devflow-platform/operators/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AppDeploymentReconciler reconciles a AppDeployment object
type AppDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *AppDeploymentReconciler) createDeployment(ctx context.Context, appDeployment *appsv1alpha1.AppDeployment) error {
	log := logf.FromContext(ctx)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appDeployment.Name + "-deployment",
			Namespace: appDeployment.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: appDeployment.APIVersion,
					Kind:       appDeployment.Kind,
					Name:       appDeployment.Name,
					UID:        appDeployment.UID,
					Controller: &[]bool{true}[0],
				},
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &[]int32{1}[0],
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": appDeployment.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": appDeployment.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  appDeployment.Name,
						Image: appDeployment.Spec.Image,
					}},
				},
			},
		},
	}

	log.Info("Creating deployment", "name", deployment.Name)
	return r.Create(ctx, deployment)
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
	log := logf.FromContext(ctx)

	// TODO(user): your logic here
	log.Info("Starting reconcile", "AppDeployment", req.NamespacedName)

	// Fetch the AppDeployment instance
	var appDeployment appsv1alpha1.AppDeployment
	if err := r.Get(ctx, req.NamespacedName, &appDeployment); err != nil {
		log.Error(err, "unable to fetch AppDeployment")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// CREATE THE DEPLOYMENT
	if err := r.createDeployment(ctx, &appDeployment); err != nil {
		log.Error(err, "Failed to create deployment")
		return ctrl.Result{}, err
	}

	// ADD THIS LOG:
	log.Info("Found AppDeployment", "name", appDeployment.Name, "image", appDeployment.Spec.Image)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AppDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.AppDeployment{}).
		Named("appdeployment").
		Complete(r)
}
