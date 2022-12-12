package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stenic/ledger/internal/pkg/utils/env"
)

var (
	logLevel  string
	logFormat string
	logCaller bool
)

var rootCommand = &cobra.Command{
	Use:   "ledger",
	Short: "Ledger",
	Run: func(c *cobra.Command, args []string) {
		c.Help()
	},
}

func init() {
	rootCommand.CompletionOptions.DisableDefaultCmd = true

	rootCommand.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		lvl, err := logrus.ParseLevel(logLevel)
		if err != nil {
			return err
		}
		logrus.SetLevel(lvl)

		switch logFormat {
		case "json":
			logrus.SetFormatter(&logrus.JSONFormatter{})
		case "text":
			logrus.SetFormatter(&logrus.TextFormatter{})
		}

		if logCaller {
			logrus.SetReportCaller(true)
		}

		return nil
	}

	rootCommand.PersistentFlags().StringVar(&logLevel, "loglevel", env.GetString("LOG_LEVEL", logrus.InfoLevel.String()), "Log level (debug, info, warn, error, fatal, panic)")
	rootCommand.PersistentFlags().StringVar(&logFormat, "logformat", env.GetString("LOG_FORMAT", "text"), "Log format (json, text)")
	rootCommand.PersistentFlags().BoolVar(&logCaller, "logcaller", env.GetBool("LOG_CALLER", false), "Log caller is printed")
}

func main() {
	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
