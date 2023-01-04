package main

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stenic/ledger/internal/pkg/utils/env"
	"github.com/stenic/ledger/internal/pkg/utils/errors"
	"github.com/stenic/ledger/internal/server"
)

func NewServerCmd() *cobra.Command {
	var (
		listenAddr string
	)

	var opts server.ServerOpts

	serverCommand := &cobra.Command{
		Use:   "server",
		Short: "Run the Ledger server",
		Long:  "This command runs Ledger server in the foreground.  It can be configured by following options.",
		Run: func(c *cobra.Command, args []string) {
			logrus.WithFields(logrus.Fields{
				"addr": listenAddr,
			}).Info("Ledger server started")
			errors.CheckError(server.NewServer(opts).Listen(listenAddr))
		},
	}

	serverCommand.Flags().StringVar(&listenAddr, "addr", env.GetString("PORT", ":8080"), "Listen on given port")
	serverCommand.Flags().StringVar(&opts.StaticAssetPath, "statisassetpath", env.GetString("STATIC_ASSET_PATH", "./ui/build"), "")
	serverCommand.Flags().StringVar(&opts.OidcIssuerURL, "oidc-issuer-url", env.GetString("OIDC_ISSUER_URL", ""), "")
	serverCommand.Flags().StringVar(&opts.OidcClientID, "oidc-client-id", env.GetString("OIDC_CLIENT_ID", ""), "")
	serverCommand.Flags().StringArrayVar(&opts.OidcAudience, "oidc-audience", strings.Split(env.GetString("OIDC_AUDIENCE", ""), ","), "")
	serverCommand.Flags().StringVar(&opts.DiscoveryNamespace, "discovery-namespace", env.GetString("DISCOVERY_NAMESPACE", ""), "")

	return serverCommand
}

func init() {
	rootCommand.AddCommand(NewServerCmd())
}
