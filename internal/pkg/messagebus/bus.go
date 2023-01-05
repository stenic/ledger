package messagebus

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/buraksezer/olric"
	"github.com/buraksezer/olric-cloud-plugin/lib"
	"github.com/buraksezer/olric/config"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type bus struct {
	ps     *olric.PubSub
	client *olric.EmbeddedClient
}

func New(opts Options) MessageBus {
	logger := logrus.WithField("scope", "messagebus")
	ctx, cancel := context.WithCancel(context.Background())

	cfg := config.New("lan")
	cfg.LogOutput = logger.WriterLevel(logrus.DebugLevel)
	cfg.Logger = log.New(logger.Writer(), "Olric: ", 0)

	if opts.DiscoveryNamespace != "" {
		cfg.ServiceDiscovery = map[string]interface{}{
			"plugin":   &lib.CloudDiscovery{},
			"provider": "k8s",
			"args":     fmt.Sprintf("namespace=%s label_selector=\"%s\"", opts.DiscoveryNamespace, opts.DiscoveryLabelSelector),
		}
	}
	cfg.Started = func() {
		defer cancel()
		logger.Info("Messagebus started")
	}
	db, err := olric.New(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	go func() {
		// Call Start at background. It's a blocker call.
		err = db.Start()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	<-ctx.Done()
	client := db.NewEmbeddedClient()
	ps, err := client.NewPubSub()
	if err != nil {
		logger.Error(err)
	}

	return &bus{
		ps:     ps,
		client: client,
	}
}

func (b bus) Publish(channel, msg string) error {
	_, err := b.ps.Publish(context.Background(), channel, msg)
	return err
}

func (b bus) Consume(channel string) <-chan *redis.Message {
	rps := b.ps.Subscribe(context.Background(), channel)
	if rps == nil {
		return nil
	}

	return rps.Channel()
}

func (b bus) Close() {
	// Don't forget the call Shutdown when you want to leave the cluster.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := b.client.Close(ctx)
	if err != nil {
		logrus.Errorf("Failed to close EmbeddedClient: %v", err)
	}
}
