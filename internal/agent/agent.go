package agent

import (
	"context"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/stenic/ledger/internal/client"
	"github.com/stenic/ledger/internal/pkg/kubernetes"
	"github.com/stenic/ledger/internal/pkg/utils/env"
	"k8s.io/client-go/util/homedir"
)

// https://github.com/dtan4/k8s-pod-notifier/blob/master/kubernetes/client.go

func Run(endpoint, namespace, location string) {
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")

	k8sClient, err := kubernetes.NewClient(kubeconfigPath)
	if err != nil {
		logrus.Fatal(err)
	}

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Make sure it's called to release resources even if no errors

	tkn := env.GetString("TOKEN", "")
	if tkn == "" {
		logrus.Fatal("Please provide a TOKEN environment variable")
	}

	notify := getNotifyFunc(client.LedgerClient{
		Endpoint: endpoint + "/query",
		Token:    tkn,
	}, location)

	wg.Add(1)
	go func() {
		logrus.Info("Watching for deployment changes")
		if err := k8sClient.WatchDeploymentEvents(ctx, namespace, notify); err != nil {
			logrus.Error(err)
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		logrus.Info("Watching for statefulset changes")
		if err := k8sClient.WatchStatefulsetEvents(ctx, namespace, notify); err != nil {
			logrus.Error(err)
			cancel()
		}
	}()

	wg.Wait()
}

func getNotifyFunc(lc client.LedgerClient, location string) kubernetes.NotifyFunc {
	return func(events []kubernetes.Event) error {
		for _, e := range events {
			logrus.WithFields(logrus.Fields{
				"application": e.Application,
				"environment": e.Environment,
				"version":     e.Version,
			}).Info("Notifying ledger")
			if e.Location == "" {
				e.Location = location
			}
			if err := lc.PostNewVersion(e.Application, e.Location, e.Environment, e.Version); err != nil {
				logrus.Error(err)
			}
		}
		return nil
	}

}
