package messagebus

type Options struct {
	EnableDiscovery    bool
	DiscoveryNamespace string
	// app.kubernetes.io/name=%s,app.kubernetes.io/instance=%s
	DiscoveryLabelSelector string
}
