package main

import (
	"github.com/spf13/cobra"
	"github.com/stenic/ledger/internal/agent"
	"github.com/stenic/ledger/internal/pkg/utils/env"
)

var opts agent.Options

func NewAgentCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Ledger cluster agent",
		Run: func(c *cobra.Command, args []string) {
			agent.Run(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Endpoint, "endpoint", env.GetString("LEDGER_ENDPOINT", "http://127.0.0.1:8080"), "Ledger endpoint url")
	cmd.Flags().StringVar(&opts.Namespaces, "namespace", env.GetString("LEDGER_NAMESPACE", ""), "Ledger namespace, empty for all")
	cmd.Flags().StringVar(&opts.Location, "location", env.GetString("LEDGER_LOCATION", ""), "Location")
	cmd.Flags().BoolVar(&opts.LeaderElection, "leader-election", env.GetBool("LEDGER_LEADER_ELECTION", false), "Leader election")
	cmd.Flags().StringVar(&opts.LeaderElectionNamespace, "leader-election-namespace", env.GetString("LEDGER_LEADER_ELECTION_NAMESPACE", "ledger"), "Leader election namespace")

	return cmd
}

func init() {
	rootCommand.AddCommand(NewAgentCmd())
}
