package psclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	MachineStateReady string = "ready"
	MachineStateOff   string = "off"
)

func (c *Client) CreateMachine(machineCreateConfig MachineCreateConfig) (*Machine, error) {
	// TODO: Handle "Get machine availability"
	// https://docs.digitalocean.com/reference/paperspace/pspace/api-reference/#operation/machineAvailability-list
	// Throw an error if machine is not available

	rb, err := json.Marshal(machineCreateConfig)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/machines", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)

	tflog.Info(*c.Context, "POST response body: "+string(res))

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

	tflog.Info(*c.Context, fmt.Sprintf("Waiting for machine event '%s' to complete, event id: %s", mashineResponse.Event.Name, mashineResponse.Event.ID))
	err = c.waitForEvent(mashineResponse.Event.ID)
	if err != nil {
		return nil, err
	}

	if machineCreateConfig.StartOnCreate {
		// Wait for machine to start
		tflog.Info(*c.Context, fmt.Sprintf("Waiting for machine '%s' to start", mashineResponse.Data.ID))

		err = c.waitForMachineState(mashineResponse.Data.ID, "ready", 30*time.Minute, 10*time.Second)
		if err != nil {
			// TODO: Handle situation when machine is created but could not start

			// log.Fatalf("Error waiting for machine state: %v", err)
			return nil, err
		}
	}

	// Fetch and return the created machine
	machine, err := c.GetMachine(mashineResponse.Data.ID)
	if err != nil {
		return nil, err
	}

	return machine, nil
}

func (c *Client) UpdateMachine(machineID string, machineUpdateConfig MachineUpdateConfig) error {
	rb, err := json.Marshal(machineUpdateConfig)
	if err != nil {
		return err
	}

	isEmpty, err := isJSONEmptyObject(string(rb))
	if err != nil {
		return fmt.Errorf("Invalid JSON: %v", err)
	} else if isEmpty {
		tflog.Info(*c.Context, "PUT request body is empty, nothing to update")
		return nil
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/machines/%s", c.HostURL, machineID), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	res, err := c.doRequest(req)
	tflog.Info(*c.Context, "PUT response body: "+string(res))
	if err != nil {
		return err
	}

	// Declare var and fill struct values with defaults
	mashineResponse := MashineResponse{}
	err = json.Unmarshal(res, &mashineResponse)
	if err != nil {
		return err
	}

	// Wait for events to finish
	tflog.Info(*c.Context, "Waiting for machine events to complete...")
	err = c.waitForMachineEvents(machineID)
	if err != nil {
		return err
	}

	return nil
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

func (c *Client) DeleteMachine(machineID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/machines/%s", c.HostURL, machineID), nil)
	if err != nil {
		return err
	}

	res, err := c.doRequest(req)

	// Check if the error is related to 404 status code
	if err != nil {
		if strings.Contains(err.Error(), "status: 404") {
			tflog.Info(*c.Context, fmt.Sprintf("Machine %s not found, assuming already deleted", machineID))
			return nil
		}
		// Other errors that are not 404
		return err
	}

	tflog.Info(*c.Context, "DELETE response body: "+string(res))

	// Check periodically if the machine still exists
	checkInterval := 10 * time.Second
	maxAttempts := 30 // Max number of attempts before giving up
	totalWaitTime := checkInterval * time.Duration(maxAttempts)
	for i := 0; i < maxAttempts; i++ {
		_, err := c.GetMachine(machineID)
		if err != nil {
			return err
		}

		// Machine still exists, wait for the next check
		time.Sleep(checkInterval)
	}

	// If we reach here, the machine still exists after the wait limit
	return fmt.Errorf("machine %s was not deleted after %f seconds", machineID, totalWaitTime.Seconds())
}

func (c *Client) ManageMachineState(machineID string, targetState string) error {
	var action string

	switch targetState {
	case MachineStateReady:
		action = "start"
	case MachineStateOff:
		action = "stop"
	default:
		return fmt.Errorf("invalid action: %s", targetState)
	}

	// Get current state
	machine, err := c.GetMachine(machineID)
	if err != nil {
		return fmt.Errorf("failed to get machine %s: %v", machineID, err)
	}

	if machine.State == targetState {
		tflog.Info(*c.Context, fmt.Sprintf("Machine '%s' is already '%s'", machineID, targetState))
		return nil
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/machines/%s/%s", c.HostURL, machineID, action), nil)
	if err != nil {
		return err
	}

	res, err := c.doRequest(req)
	tflog.Info(*c.Context, "PATCH response body: "+string(res))
	if err != nil {
		return err
	}

	// Wait for action to complete
	tflog.Info(*c.Context, fmt.Sprintf("Waiting for machine '%s' to %s", machineID, action))
	err = c.waitForMachineState(machineID, targetState, 30*time.Minute, 10*time.Second)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) waitForMachineState(machineID string, desiredState string, timeout time.Duration, pollInterval time.Duration) error {
	// Create a ticker for polling and a timeout channel
	ticker := time.NewTicker(pollInterval)
	timeoutChan := time.After(timeout)

	defer ticker.Stop() // Ensure ticker is stopped after the function exits

	for {
		select {
		case <-timeoutChan:
			return fmt.Errorf("timeout reached waiting for machine %s to reach state: %s", machineID, desiredState)

		case <-ticker.C:
			machine, err := c.GetMachine(machineID)
			if err != nil {
				return fmt.Errorf("failed to get machine %s: %v", machineID, err)
			}

			if machine.State == desiredState {
				return nil // Return if the desired state is reached
			}
		}
	}
}

func isJSONEmptyObject(jsonStr string) (bool, error) {
	var obj map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return false, err // Return error if JSON is invalid
	}
	return len(obj) == 0, nil // Check if the map is empty
}
