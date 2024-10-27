package ppclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Returns number of each found machine event type (name)
func (c *Client) GetMachineEventsStateStat() (error, map[string]int) {
	allItems := []Event{}
	stat := map[string]int{}
	params := map[string]string{}

	err := fetchAllItems(c, &allItems, "machine-events", params)
	if err != nil {
		return err, nil
	}

	for _, event := range allItems {
		stat[event.Name]++
	}

	stat["_totalEventsProcessed"] = len(allItems)

	return nil, stat
}

func (c *Client) waitForEvent(eventID string) (*Event, error) {
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
			return nil, err
		}
		resp, err := c.doRequest(req)
		if err != nil {
			return nil, err
		}

		// Parse the response
		err = json.Unmarshal(resp, &event)
		if err != nil {
			return nil, err
		}

		// Check the event state
		if event.State == "done" {
			return &event, nil
		}
		if event.Error != nil {
			return nil, fmt.Errorf("error during event processing: %s", *event.Error)
		}

		// There is a case, when {"state": "error", "error": null}
		if event.State == "error" {
			return nil, fmt.Errorf("unknown error during event %s processing", eventID)
		}

		// Wait for either the polling interval or the timeout
		select {
		case <-time.After(pollInterval):
		case <-timer:
			return nil, fmt.Errorf("timeout reached while waiting for event %s", eventID)
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
		// Skip 'done' events, so nothing to wait
		if event.State == "done" {
			continue
		}

		tflog.Info(*c.Context, fmt.Sprintf("Waiting for machine event '%s' to complete, event id: %s", event.Name, event.ID))

		_, err := c.waitForEvent(event.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
