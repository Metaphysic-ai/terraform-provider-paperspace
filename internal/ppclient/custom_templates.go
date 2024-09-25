package ppclient

import (
	"encoding/json"
	"fmt"
	"net/http"

	"terraform-provider-paperspace/internal/ppclient/models"
)

func (c *Client) GetCustomTemplates() ([]models.CustomTemplate, error) {
	var allItems []models.CustomTemplate
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
		res := &models.CustomTemplatesResponse{}
		err = json.Unmarshal(body, res)
		if err != nil {
			return nil, err
		}

		allItems = append(allItems, res.Items...)

		nextPage = res.NextPage
		hasMore = res.HasMore
	}

	return allItems, nil
}
