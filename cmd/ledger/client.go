package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

			ep := endpoint + "/query"
			logrus.Infof("Sending data to %s", ep)

			jsonData := map[string]string{
				"query": fmt.Sprintf(`
					mutation { 
						createVersion(input: {
							application: "%s",
							environment: "%s",
							version: "%s"
						}) {
							id
						}
					}
				`, args[0], args[1], args[2]),
			}
			jsonValue, _ := json.Marshal(jsonData)
			logrus.Debugf("Payload: %s", jsonValue)

			request, err := http.NewRequest("POST", ep, bytes.NewBuffer(jsonValue))
			request.Header.Add("Authorization", "Bearer "+tkn)
			request.Header.Add("Content-Type", "application/json")
			if err != nil {
				c.PrintErr(err)
				os.Exit(1)
			}

			client := &http.Client{Timeout: time.Second * 10}
			response, err := client.Do(request)
			if err != nil {
				if err != nil {
					c.PrintErr("The HTTP request failed with error: " + err.Error())
					os.Exit(1)
				}
			}
			defer response.Body.Close()

			data, _ := ioutil.ReadAll(response.Body)
			logrus.Debug("Response: %s", data)
			logrus.Info("Version created in ledger")
		},
		Args: cobra.ExactArgs(3),
	}

	cmd.Flags().StringVar(&endpoint, "endpoint", env.GetString("LEDGER_ENDPOINT", "http://127.0.0.1:8080"), "Ledger endpoint url")

	return cmd
}

func init() {
	rootCommand.AddCommand(NewClientCmd())
}
