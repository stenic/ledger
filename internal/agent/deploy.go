package agent

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	v1apps "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type LedgerDeploymemtReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder

	RefreshRate time.Duration
}

func (r *LedgerDeploymemtReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var deploy v1apps.Deployment
	if err := r.Get(ctx, req.NamespacedName, &deploy); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	val, hasAnnotation := deploy.Annotations[ledgerAnnotationProcessedGen]
	currGen := strconv.FormatInt(deploy.Generation, 10)
	if !hasAnnotation || val != currGen {
		deploy.Annotations[ledgerAnnotationProcessedGen] = strconv.FormatInt(deploy.Generation+1, 10)
		if err := r.Update(ctx, &deploy); err != nil {
			return ctrl.Result{}, err
		}

		if err := processObject(deploy.ObjectMeta, deploy.Spec.Template); err != nil {
			logrus.Error(err)
			r.Recorder.Event(&deploy, "Warning", "Notified", fmt.Sprintf("Ledger notified failed - %s", err))
		} else {
			r.Recorder.Event(&deploy, "Normal", "Notified", "Ledger was notified")
		}
	}

	return ctrl.Result{}, nil
}

func (r *LedgerDeploymemtReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1apps.Deployment{}).
		Complete(r)
}
