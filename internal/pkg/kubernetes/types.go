package kubernetes

// https://github.com/dtan4/k8s-pod-notifier/blob/master/kubernetes/client.go

type Event struct {
	Environment string
	Application string
	Version     string
}

// NotifyFunc represents callback function for Pod event
type NotifyFunc func(events []Event) error
