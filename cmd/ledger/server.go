package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stenic/ledger/internal/pkg/utils/env"
	"github.com/stenic/ledger/internal/pkg/utils/errors"
	"github.com/stenic/ledger/internal/server"
)

func NewServerCmd() *cobra.Command {
	var (
		listenAddr      string
		staticAssetPath string
		oidcIssuerURL   string
		oidcClientID    string
	)

	serverCommand := &cobra.Command{
		Use:   "server",
		Short: "Run the Ledger server",
		Long:  "This command runs Ledger server in the foreground.  It can be configured by following options.",
		Run: func(c *cobra.Command, args []string) {
			logrus.WithFields(logrus.Fields{
				"addr": listenAddr,
			}).Info("Ledger server started")

			opts := server.ServerOpts{
				StaticAssetPath: staticAssetPath,
				OidcIssuerURL:   oidcIssuerURL,
				OidcClientID:    oidcClientID,
			}

			errors.CheckError(server.NewServer(opts).Listen(listenAddr))
		},
	}

	serverCommand.Flags().StringVar(&staticAssetPath, "statisassetpath", env.GetString("STATIC_ASSET_PATH", "./ui/build"), "")
	serverCommand.Flags().StringVar(&listenAddr, "addr", env.GetString("PORT", ":8080"), "Listen on given port")
	serverCommand.Flags().StringVar(&oidcIssuerURL, "oidc-issuer-url", env.GetString("OIDC_ISSUER_URL", ""), "")
	serverCommand.Flags().StringVar(&oidcClientID, "oidc-client-id", env.GetString("OIDC_CLIENT_ID", ""), "")

	return serverCommand
}

func init() {
	rootCommand.AddCommand(NewServerCmd())
}
