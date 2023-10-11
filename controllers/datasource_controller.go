/*
Copyright 2023 KubeAGI.

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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/go-logr/logr"
	arcadiav1alpha1 "github.com/kubeagi/arcadia/api/v1alpha1"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	rootUser     = "rootUser"
	rootPassword = "rootPassword"
)

// DatasourceReconciler reconciles a Datasource object
type DatasourceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=arcadia.kubeagi.k8s.com.cn,resources=datasources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=arcadia.kubeagi.k8s.com.cn,resources=datasources/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=arcadia.kubeagi.k8s.com.cn,resources=datasources/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Datasource object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *DatasourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Starting datasource reconcile")

	instance := &arcadiav1alpha1.Datasource{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if errors.IsNotFound(err) {
			// datasourcce has been deleted.
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if err := r.Checkdatasource(ctx, logger, instance); err != nil {
		logger.Error(err, "Failed to call LLM")
		// Update conditioned status
		return reconcile.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DatasourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&arcadiav1alpha1.Datasource{}).
		Complete(r)
}

// CheckLLM updates new LLM instance.
func (r *DatasourceReconciler) Checkdatasource(ctx context.Context, logger logr.Logger, instance *arcadiav1alpha1.Datasource) error {
	logger.Info("check datasource")
	endpoint := instance.Spec.URL
	secret := corev1.Secret{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: "kubebb-addons", Name: instance.Spec.AuthSecret}, &secret); err != nil {
		return err
	}
	accessKeyID := string(secret.Data[rootUser])
	secretAccessKey := string(secret.Data[rootPassword])
	useSSL := true

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return err
	}
	err = minioClient.MakeBucket(ctx, "mybucket", minio.MakeBucketOptions{Region: "local"})
	if err != nil {
		return err
	}
	location, err := minioClient.GetBucketLocation(ctx, "mybucket")
	logger.Info(location)
	if err != nil {
		return err

	}

	return nil
}
