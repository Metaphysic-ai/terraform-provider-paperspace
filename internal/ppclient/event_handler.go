package ppclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (c *Client) WaitForEvent(eventID string) (*Event, error) {
	url := fmt.Sprintf("%s/machine-events/%s", c.HostURL, eventID)

	var event Event
	pollInterval := 5 * time.Second

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

		// Wait before polling again
		time.Sleep(pollInterval)
	}
}
