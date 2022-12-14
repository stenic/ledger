package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stenic/ledger/internal/client"
	"github.com/stenic/ledger/internal/pkg/utils/env"
)

func NewClientCmd() *cobra.Command {
	clientCommand := &cobra.Command{
		Use:   "client",
		Short: "Ledger client",
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}

	clientCommand.AddCommand(
		NewClientAddVersionCmd(),
	)

	return clientCommand
}

func NewClientAddVersionCmd() *cobra.Command {
	var (
		endpoint string
		location string
	)
	cmd := &cobra.Command{
		Use:   "new-version APPLICATION ENVIRONMENT VERSION",
		Short: "Send a version",
		Run: func(c *cobra.Command, args []string) {
			c.Println(os.Executable())

			tkn := env.GetString("TOKEN", "")
			if tkn == "" {
				c.PrintErrln("Please provide a TOKEN environment variable")
				os.Exit(1)
			}

			lc := client.LedgerClient{
				Endpoint: endpoint + "/query",
				Token:    tkn,
			}
			logrus.Infof("Sending data to %s", lc.Endpoint)
			if err := lc.PostNewVersion(args[0], location, args[1], args[2]); err != nil {
				logrus.Error(err)
			}

			logrus.Info("Version created in ledger")
		},
		Args: cobra.ExactArgs(3),
	}

	cmd.Flags().StringVar(&endpoint, "endpoint", env.GetString("LEDGER_ENDPOINT", "http://127.0.0.1:8080"), "Ledger endpoint url")
	cmd.Flags().StringVar(&location, "location", env.GetString("LEDGER_LOCATION", ""), "Location")

	return cmd
}

func init() {
	rootCommand.AddCommand(NewClientCmd())
}
