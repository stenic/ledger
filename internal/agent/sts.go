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

type LedgerStatefulsetReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder

	RefreshRate time.Duration
}

func (r *LedgerStatefulsetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var statefulset v1apps.StatefulSet
	if err := r.Get(ctx, req.NamespacedName, &statefulset); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	val, hasAnnotation := statefulset.Annotations[ledgerAnnotationProcessedGen]
	currGen := strconv.FormatInt(statefulset.Generation, 10)
	if !hasAnnotation || val != currGen {
		statefulset.Annotations[ledgerAnnotationProcessedGen] = strconv.FormatInt(statefulset.Generation+1, 10)
		if err := r.Update(ctx, &statefulset); err != nil {
			return ctrl.Result{}, err
		}

		if causedByScaling(ctx, r.Client, &statefulset) {
			return ctrl.Result{}, nil
		}

		if err := processObject(statefulset.ObjectMeta, statefulset.Spec.Template); err != nil {
			logrus.Error(err)
			r.Recorder.Event(&statefulset, "Warning", "Notified", fmt.Sprintf("Ledger notified failed - %s", err))
		} else {
			r.Recorder.Event(&statefulset, "Normal", "Notified", "Ledger was notified")
		}
	}

	return ctrl.Result{}, nil
}

func (r *LedgerStatefulsetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1apps.StatefulSet{}).
		Complete(r)
}
