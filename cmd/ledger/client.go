package main

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/cobra"
	"github.com/stenic/ledger/internal/auth"
	"github.com/stenic/ledger/internal/pkg/client"
	"github.com/stenic/ledger/internal/storage"
)

func NewClientCmd() *cobra.Command {
	clientCommand := &cobra.Command{
		Use:   "client",
		Short: "Manage clients",
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}

	clientCommand.AddCommand(
		NewClientAddTokenCmd(),
	)

	return clientCommand
}

func NewClientAddTokenCmd() *cobra.Command {
	var (
		ttlDays int
	)

	cmd := &cobra.Command{
		Use:   "new-token 'username'",
		Short: "Create a new token",
		Run: func(c *cobra.Command, args []string) {
			engine := storage.Database{}
			engine.InitDB()
			defer engine.CloseDB()
			engine.Migrate()

			newClient, err := client.CreateClient(args[0])
			if err != nil {
				c.PrintErr(err)
			}

			tkn, err := auth.GenerateToken(newClient.Username,
				func(claims *auth.Claims) {
					claims.RegisteredClaims.ID = newClient.ID
					claims.RegisteredClaims.Issuer = client.LocalClientIssuer
					claims.RegisteredClaims.NotBefore = jwt.NewNumericDate(time.Now())
					claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(365 * 24 * time.Hour))
				})
			if err != nil {
				c.PrintErr(err)
			}

			c.Println(tkn)
		},
		Args: cobra.ExactArgs(1),
	}

	cmd.Flags().IntVar(&ttlDays, "ttldays", 365, "TTL in days for the token")

	return cmd
}

func init() {
	rootCommand.AddCommand(NewClientCmd())
}
