package agent

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stenic/ledger/internal/client"
	"github.com/stenic/ledger/internal/pkg/utils/env"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/bombsimon/logrusr/v4"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
)

var (
	scheme   = runtime.NewScheme()
	location = ""
	lc       client.LedgerClient
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
}

type Options struct {
	Endpoint   string
	Namespaces string
	Location   string

	LeaderElection          bool
	LeaderElectionNamespace string
}

func Run(opts Options) {
	location = opts.Location
	ctrl.SetLogger(logrusr.New(logrus.StandardLogger()))
	mgrOpts := ctrl.Options{
		Scheme:                  scheme,
		MetricsBindAddress:      ":8082",
		Port:                    9443,
		HealthProbeBindAddress:  ":8081",
		Namespace:               opts.Namespaces,
		LeaderElection:          opts.LeaderElection,
		LeaderElectionID:        "ledger.stenic.io",
		LeaderElectionNamespace: opts.LeaderElectionNamespace,
	}
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), mgrOpts)
	if err != nil {
		logrus.Fatal(err)
	}

	tkn := env.GetString("TOKEN", "")
	if tkn == "" {
		logrus.Fatal("Please provide a TOKEN environment variable")
	}
	lc = client.LedgerClient{
		Endpoint: opts.Endpoint + "/query",
		Token:    tkn,
	}

	interval := 5 * time.Second

	if err = (&LedgerStatefulsetReconciler{
		Client:      mgr.GetClient(),
		Scheme:      mgr.GetScheme(),
		RefreshRate: interval,
		Recorder:    mgr.GetEventRecorderFor("ledger"),
	}).SetupWithManager(mgr); err != nil {
		logrus.Fatal(err)
	}

	if err = (&LedgerDeploymemtReconciler{
		Client:      mgr.GetClient(),
		Scheme:      mgr.GetScheme(),
		RefreshRate: interval,
		Recorder:    mgr.GetEventRecorderFor("ledger"),
	}).SetupWithManager(mgr); err != nil {
		logrus.Fatal(err)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		logrus.Fatal(err, "unable to set up health check")
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		logrus.Fatal(err, "unable to set up ready check")
	}

	logrus.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		logrus.Fatal(err, "problem running manager")
	}
}

func notify(events []Event) error {
	for _, e := range events {
		if e.Location == "" {
			e.Location = location
		}
		logrus.WithFields(logrus.Fields{
			"application": e.Application,
			"location":    e.Location,
			"environment": e.Environment,
			"version":     e.Version,
		}).Info("Notifying ledger")
		if err := lc.PostNewVersion(e.Application, e.Location, e.Environment, e.Version); err != nil {
			logrus.Error(err)
			return err
		}
	}

	return nil
}
