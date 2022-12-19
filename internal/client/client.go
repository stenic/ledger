package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sirupsen/logrus"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

type LedgerClient struct {
	Endpoint string
	Token    string
}

func (c *LedgerClient) PostNewVersion(app, location, env, version string) error {
	jsonData := map[string]string{
		"query": fmt.Sprintf(`
			mutation { 
				createVersion(input: {
					application: "%s",
					environment: "%s",
					version: "%s",
					location: "%s",
				}) {
					id
				}
			}
		`, app, env, version, location),
	}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	r := regexp.MustCompile(`(\\[tn]|\s+)`)
	logrus.Tracef("Payload: %s", r.ReplaceAllString(string(jsonValue), ""))

	request, err := http.NewRequest("POST", c.Endpoint, bytes.NewBuffer(jsonValue))
	request.Header.Add("Authorization", "Bearer "+c.Token)
	request.Header.Add("Content-Type", "application/json")
	if err != nil {
		return err
	}

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10
	retryClient.Backoff = retryablehttp.DefaultBackoff
	retryClient.Logger = logrus.StandardLogger()
	httpClient := retryClient.StandardClient()
	httpClient.Timeout = time.Second * 10
	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	logrus.Tracef("Response: %s", string(body))

	var resp = gqlResponse{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return err
	}

	if len(resp.Errors) > 0 {
		return fmt.Errorf(resp.Errors.Error())
	}

	return err
}

type gqlResponse struct {
	Errors gqlerror.List `json:"errors"`
}
