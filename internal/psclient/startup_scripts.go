package psclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (c *Client) CreateStartupScript(startupScriptCreateConfig StartupScriptCreateConfig) (*StartupScript, error) {
	rb, err := json.Marshal(startupScriptCreateConfig)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/startup-scripts", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)

	tflog.Info(*c.Context, "POST response body: "+string(res))

	if err != nil {
		return nil, err
	}

	// Declare var and fill struct values with defaults
	startupScript := StartupScript{}
	err = json.Unmarshal(res, &startupScript)
	if err != nil {
		return nil, err
	}

	return &startupScript, nil
}

func (c *Client) GetStartupScript(id string) (*StartupScript, error) {
	url := fmt.Sprintf("%s/startup-scripts/%s", c.HostURL, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	startupScript := StartupScript{}
	err = json.Unmarshal(body, &startupScript)
	if err != nil {
		return nil, err
	}

	return &startupScript, nil
}

func (c *Client) DeleteStartupScript(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/startup-scripts/%s", c.HostURL, id), nil)
	if err != nil {
		return err
	}

	res, err := c.doRequest(req)

	// Check if the error is related to 404 status code
	if err != nil {
		if strings.Contains(err.Error(), "status: 404") {
			tflog.Info(*c.Context, fmt.Sprintf("Startup script %s not found, assuming already deleted", id))
			return nil
		}
		// Other errors that are not 404
		return err
	}

	tflog.Info(*c.Context, "DELETE response body: "+string(res))

	// Check periodically if the resource still exists
	checkInterval := 10 * time.Second
	maxAttempts := 18 // Max number of attempts before giving up
	totalWaitTime := checkInterval * time.Duration(maxAttempts)
	for i := 0; i < maxAttempts; i++ {
		_, err := c.GetStartupScript(id)
		if err != nil {
			return err
		}

		// Resource still exists, wait for the next check
		time.Sleep(checkInterval)
	}

	// If we reach here, the resource still exists after the wait limit
	return fmt.Errorf("startup script %s was not deleted after %f seconds", id, totalWaitTime.Seconds())
}
