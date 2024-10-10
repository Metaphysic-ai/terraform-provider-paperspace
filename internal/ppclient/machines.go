package ppclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// TODO: Wait until is started if 'create on start' is true ("state": "starting")
// TODO: Find better way to pass 'ctx' to the client for logging
func (c *Client) CreateMachine(machineConfig MachineConfig, ctx context.Context) (*Machine, error) {
	rb, err := json.Marshal(machineConfig)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/machines", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)

	tflog.Info(ctx, "POST response body: "+string(res))

	if err != nil {
		return nil, err
	}

	// Declare var and fill struct values with defaults
	mashineResponse := MashineResponse{}
	err = json.Unmarshal(res, &mashineResponse)
	if err != nil {
		return nil, err
	}

	// WaitForEvent Section

	tflog.Info(ctx, fmt.Sprintf("Waiting for machine event '%s' to complete, event id: %s", mashineResponse.Event.Name, mashineResponse.Event.ID))

	event, err := c.WaitForEvent(mashineResponse.Event.ID)
	if err != nil {
		return nil, err
	}

	if event.State != "done" {
		return nil, fmt.Errorf("event not completed successfully: %v", event.State)
	}

	// Fetch and return the created machine
	machine, err := c.GetMachine(event.MachineID)
	if err != nil {
		return nil, err
	}

	return machine, nil
}

func (c *Client) GetMachine(machineID string) (*Machine, error) {
	url := fmt.Sprintf("%s/machines/%s", c.HostURL, machineID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	machine := Machine{}
	err = json.Unmarshal(body, &machine)
	if err != nil {
		return nil, err
	}

	return &machine, nil
}

// TODO: Implement update function

func (c *Client) DeleteMachine(machineID string, ctx context.Context) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/machines/%s", c.HostURL, machineID), nil)
	if err != nil {
		return err
	}

	res, err := c.doRequest(req)

	// Check if the error is related to 404 status code
	if err != nil {
		if strings.Contains(err.Error(), "status: 404") {
			tflog.Info(ctx, fmt.Sprintf("Machine %s not found, assuming already deleted", machineID))
			return nil
		}
		// Other errors that are not 404
		return err
	}

	tflog.Info(ctx, "DELETE response body: "+string(res))

	// Check periodically if the machine still exists
	checkInterval := 10 * time.Second
	maxAttempts := 30 // Max number of attempts before giving up
	totalWaitTime := checkInterval * time.Duration(maxAttempts)
	for i := 0; i < maxAttempts; i++ {
		_, err := c.GetMachine(machineID)
		if err != nil {
			tflog.Info(ctx, fmt.Sprintf("Machine %s not found, assuming has been deleted successfully", machineID))
			return nil
		}

		// Machine still exists, wait for the next check
		time.Sleep(checkInterval)
	}

	// If we reach here, the machine still exists after the wait limit
	return fmt.Errorf("machine %s was not deleted after %s seconds", machineID, totalWaitTime.Seconds())
}
