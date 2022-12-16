package kubernetes

// https://github.com/dtan4/k8s-pod-notifier/blob/master/kubernetes/client.go

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"

	v1apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	clientset *kubernetes.Clientset
}

// NewClient creates Client object using local kubecfg
func NewClient(kubeconfig string) (*Client, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load clientset")
	}

	return &Client{
		clientset: clientset,
	}, nil
}

func (c *Client) WatchDeploymentEvents(ctx context.Context, namespace string, notifyFunc NotifyFunc) error {
	watcher, err := c.clientset.AppsV1().Deployments(namespace).Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "cannot create Deployment event watcher")
	}

	go func() {
		for {
			select {
			case e := <-watcher.ResultChan():
				if e.Object == nil {
					return
				}
				deployment, ok := e.Object.(*v1apps.Deployment)
				if !ok {
					continue
				}
				if ev := processObject(e, deployment.ObjectMeta, deployment.Spec.Template); ev != nil {
					notifyFunc(ev)
				}
			case <-ctx.Done():
				watcher.Stop()
				return
			}
		}
	}()

	return nil
}

func (c *Client) WatchStatefulsetEvents(ctx context.Context, namespace string, notifyFunc NotifyFunc) error {
	watcher, err := c.clientset.AppsV1().StatefulSets(namespace).Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "cannot create StatefulSet event watcher")
	}

	go func() {
		for {
			select {
			case e := <-watcher.ResultChan():
				if e.Object == nil {
					return
				}
				statefulset, ok := e.Object.(*v1apps.StatefulSet)
				if !ok {
					continue
				}
				if ev := processObject(e, statefulset.ObjectMeta, statefulset.Spec.Template); ev != nil {
					notifyFunc(ev)
				}
			case <-ctx.Done():
				watcher.Stop()
				return
			}
		}
	}()

	return nil
}

var cache = map[string]int64{}

func processObject(e watch.Event, objectMeta metav1.ObjectMeta, podTemplate v1.PodTemplateSpec) []Event {
	val, ok := cache[string(objectMeta.UID)]
	if ok && val == objectMeta.Generation {
		// Skip processing if the object was processed in the last 10 seconds.
		logrus.WithFields(logrus.Fields{
			"namespace": objectMeta.Namespace,
			"name":      objectMeta.Name,
		}).Trace("Skipping, recently processed")
		return nil
	}
	cache[string(objectMeta.UID)] = objectMeta.Generation

	switch e.Type {
	case watch.Added:
		if objectMeta.CreationTimestamp.Time.Before(time.Now().Add(-1 * time.Minute)) {
			// Skip processing if the object was created longer that 1 minute ago.
			// Most likely we are starting up and getting all the objects a first time.
			logrus.WithFields(logrus.Fields{
				"namespace": objectMeta.Namespace,
				"name":      objectMeta.Name,
			}).Debug("Skipping, assuming startup")
			return nil
		}
		return collectObjectInfo(objectMeta, podTemplate)

	case watch.Modified:
		if objectMeta.DeletionTimestamp != nil {
			// Skip if we are deleting the resource
			logrus.WithFields(logrus.Fields{
				"namespace": objectMeta.Namespace,
				"name":      objectMeta.Name,
			}).Debug("Skipping, assuming deletion")
			return nil
		}
		return collectObjectInfo(objectMeta, podTemplate)
	}

	return nil
}

func getImages(spec v1.PodTemplateSpec) (images []string) {
	for _, container := range spec.Spec.Containers {
		images = append(images, container.Image)
	}

	return images
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
