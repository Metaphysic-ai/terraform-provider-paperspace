package psclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Returns number of each found machine event type (name).
func (c *Client) GetMachineEventsStateStat() (map[string]int, error) {
	allItems := []Event{}
	stat := map[string]int{}
	params := map[string]string{}

	err := fetchAllItems(c, &allItems, "machine-events", params)
	if err != nil {
		return nil, err
	}

	for _, event := range allItems {
		stat[event.Name]++
	}

	stat["_totalEventsProcessed"] = len(allItems)

	return stat, nil
}

func (c *Client) waitForEvent(eventID string) error {
	// The "state": "error" typically indicates an issue during the execution of a specific action on the machine (like a resource update).
	// However, if the "error" field is "null", it could mean the API did not capture any specific error details.
	// In these cases, "state": "error" may be reporting a transient issue that did not prevent successful completion,
	// or it may reflect a system state that was corrected automatically.
	//
	// The presence of a "dtFinished" timestamp suggests that the operation concluded,
	// which is an additional hint that the operation likely completed despite the error state.
	// This detail can help confirm that the event doesnt need further attention.
	//
	// Given these points, the following logic should be implemented:
	//  - Ignore "state": "error" when "error": null and the action completed as expected without user intervention.
	//  - Monitor for recurring "state": "error" entries, as this could indicate underlying issues even if individual operations appear to succeed.

	url := fmt.Sprintf("%s/machine-events/%s", c.HostURL, eventID)

	var event Event
	pollInterval := 5 * time.Second
	timeout := 30 * time.Minute

	// Create a timer for the timeout
	timer := time.After(timeout)

	// Poll with interval until event is done or error occurs
	for {
		// Make the request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
		resp, err := c.doRequest(req)
		if err != nil {
			return err
		}

		// Parse the response
		err = json.Unmarshal(resp, &event)
		if err != nil {
			return err
		}

		// Check the event state
		if event.State == "done" {
			return nil
		}

		if event.Error != nil {
			return fmt.Errorf("error during event processing: %s", *event.Error)
		}

		if event.DtFinished != nil {
			return nil
		}

		// There is a case, when {"state": "error", "error": null, "dtFinished": null}
		if event.State == "error" {
			return fmt.Errorf("unknown error during event %s processing", eventID)
		}

		// Wait for either the polling interval or the timeout
		select {
		case <-time.After(pollInterval):
		case <-timer:
			return fmt.Errorf("timeout reached while waiting for event %s", eventID)
		}
	}
}

func (c *Client) waitForMachineEvents(machineID string) error {
	// TODO: Add option to choose which events to ignore or wait
	allItems := []Event{}
	params := map[string]string{"machineId": machineID}

	err := fetchAllItems(c, &allItems, "machine-events", params)
	if err != nil {
		return err
	}

	for _, event := range allItems {
		// Skip 'finished' events, so nothing to wait
		if event.DtFinished != nil {
			continue
		}

		tflog.Info(*c.Context, fmt.Sprintf("Waiting for machine event '%s' to complete, event id: %s", event.Name, event.ID))

		err := c.waitForEvent(event.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
