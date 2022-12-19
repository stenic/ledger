package agent

import (
	"strings"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var cache = map[string]int64{}

func processObject(objectMeta metav1.ObjectMeta, podTemplate v1.PodTemplateSpec) error {
	val, ok := cache[string(objectMeta.UID)]
	if ok && val == objectMeta.Generation {
		// Skip processing if the object was processed in the last 10 seconds.
		logrus.WithFields(logrus.Fields{
			"namespace":  objectMeta.Namespace,
			"name":       objectMeta.Name,
			"generation": objectMeta.Generation,
		}).Trace("Skipping, recently processed")
		return nil
	}
	cache[string(objectMeta.UID)] = objectMeta.Generation

	if objectMeta.DeletionTimestamp != nil {
		// Skip if we are deleting the resource
		logrus.WithFields(logrus.Fields{
			"namespace":  objectMeta.Namespace,
			"name":       objectMeta.Name,
			"generation": objectMeta.Generation,
		}).Debug("Skipping MOD, assuming deletion")

		return nil
	}
	return notify(collectObjectInfo(objectMeta, podTemplate))
}

func collectObjectInfo(objectMeta metav1.ObjectMeta, podTemplate v1.PodTemplateSpec) (events []Event) {
	app := objectMeta.Name
	env := objectMeta.Namespace
	loc := ""

	if val, ok := objectMeta.Annotations[LedgerAnnotationApplication]; ok {
		app = val
	}
	if val, ok := objectMeta.Annotations[LedgerAnnotationEnvironment]; ok {
		env = val
	}
	if val, ok := objectMeta.Annotations[LedgerAnnotationLocation]; ok {
		loc = val
	}

	for _, img := range getImages(podTemplate) {
		events = append(events, Event{
			Application: app,
			Environment: env,
			Location:    loc,
			Version:     strings.Split(img, ":")[1],
		})
	}

	return events
}

func getImages(spec v1.PodTemplateSpec) (images []string) {
	for _, container := range spec.Spec.Containers {
		images = append(images, container.Image)
	}

	return images
}
