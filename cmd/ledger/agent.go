package main

import (
	"github.com/spf13/cobra"
	"github.com/stenic/ledger/internal/agent"
	"github.com/stenic/ledger/internal/pkg/utils/env"
)

func NewAgentCmd() *cobra.Command {
	var (
		endpoint string
	)

	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Ledger cluster agent",
		Run: func(c *cobra.Command, args []string) {
			agent.Run(endpoint)
		},
	}

	cmd.Flags().StringVar(&endpoint, "endpoint", env.GetString("LEDGER_ENDPOINT", "http://127.0.0.1:8080"), "Ledger endpoint url")

	return cmd
}

func init() {
	rootCommand.AddCommand(NewAgentCmd())
}
