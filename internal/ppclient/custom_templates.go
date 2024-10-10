package ppclient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetCustomTemplates() ([]CustomTemplate, error) {
	var allItems []CustomTemplate
	reqLimit := 120 // 1..120
	nextPage := ""
	hasMore := true

	for hasMore {
		// Build the request URL, append nextPage if it's not empty
		url := fmt.Sprintf("%s/custom-templates", c.HostURL)
		if nextPage != "" {
			url = fmt.Sprintf("%s?limit=%d&after=%s", url, reqLimit, nextPage) // possible options: name, machineId, limit
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		// Execute the request and get the response body
		body, err := c.doRequest(req)
		if err != nil {
			return nil, err
		}

		// Unmarshal the response into the struct
		res := &CustomTemplatesResponse{}
		err = json.Unmarshal(body, res)
		if err != nil {
			return nil, err
		}

		allItems = append(allItems, res.Items...)

		nextPage = res.NextPage
		hasMore = res.HasMore
	}

	// API returns duplicates. Workaround: get only unique items based on ID
	uniqueItems := make(map[string]CustomTemplate)
	for _, item := range allItems {
		uniqueItems[item.ID] = item
	}
	allItems = []CustomTemplate{}
	for _, item := range uniqueItems {
		allItems = append(allItems, item)
	}

	return allItems, nil
}
