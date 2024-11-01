package ppclient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ItemsResponse[T any] struct {
	Items    []T    `json:"items"`
	NextPage string `json:"nextPage"`
	HasMore  bool   `json:"hasMore"`
}

func fetchAllItems[T any](c *Client, allItems *[]T, path string, params map[string]string) error {
	nextPage := ""
	hasMore := true
	params["limit"] = "120"

	for hasMore {
		// Build the request URL, append nextPage if it's not empty
		url := fmt.Sprintf("%s/%s", c.HostURL, path)

		if nextPage != "" {
			params["after"] = nextPage
		}

		url = fmt.Sprintf("%s?%s", url, buildQueryString(params)) // Custom request arguments

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		// Execute the request and get the response body
		body, err := c.doRequest(req)
		if err != nil {
			return err
		}

		// Create a response struct for the given type
		res := &ItemsResponse[T]{}
		err = json.Unmarshal(body, res)
		if err != nil {
			return err
		}

		// Append items to allItems
		*allItems = append(*allItems, res.Items...)

		nextPage = res.NextPage
		hasMore = res.HasMore
	}

	return nil
}

// Helper function to build the query string from a map of parameters.
func buildQueryString(params map[string]string) string {
	q := ""
	for k, v := range params {
		q += fmt.Sprintf("%s=%s&", k, v)
	}
	return q[:len(q)-1] // Remove trailing '&'
}
