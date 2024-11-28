package psclient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetPrivateNetwork(id string) (*PrivateNetwork, error) {
	url := fmt.Sprintf("%s/private-networks/%s", c.HostURL, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resBody, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	// Handle "null" response (network does not exist, or it's non-private)
	if string(resBody) == "null" {
		return nil, nil
	}

	privateNetwork := PrivateNetwork{}
	err = json.Unmarshal(resBody, &privateNetwork)
	if err != nil {
		return nil, err
	}

	return &privateNetwork, nil
}
