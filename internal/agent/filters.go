package agent

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	v1core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func causedByScaling(ctx context.Context, k8sClient client.Client, obj client.Object) bool {
	events := v1core.EventList{}
	k8sClient.List(ctx, &events, &client.ListOptions{
		FieldSelector: fields.SelectorFromSet(fields.Set{
			"involvedObject.kind":      obj.GetObjectKind().GroupVersionKind().Kind,
			"involvedObject.name":      obj.GetName(),
			"involvedObject.namespace": obj.GetNamespace(),
			"reason":                   "ScalingReplicaSet",
		}),
	})
	for _, e := range events.Items {
		if e.LastTimestamp.Time.After(time.Now().Add(-10 * time.Second)) {
			logrus.WithFields(logrus.Fields{
				"namespace": obj.GetNamespace(),
				"name":      obj.GetName(),
			}).Debug("Skipping, caused by ScalingReplicaSet")
			return true
		}
	}

	return false
}
