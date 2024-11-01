package ppclient

import (
	"fmt"
	"sort"
	"time"
)

func (c *Client) GetCustomTemplates() (*[]CustomTemplate, error) {
	allItems := []CustomTemplate{}
	params := map[string]string{}

	err := fetchAllItems(c, &allItems, "custom-templates", params)
	if err != nil {
		return nil, err
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

	err = sortCustomTemplates(allItems, "DtCreated")
	if err != nil {
		return nil, err
	}

	return &allItems, nil
}

func sortCustomTemplates(templates []CustomTemplate, sortBy string) error {
	switch sortBy {
	case "ID":
		sort.Slice(templates, func(i, j int) bool {
			return templates[i].ID < templates[j].ID
		})
		return nil
	case "DtCreated":
		sort.Slice(templates, func(i, j int) bool {
			time1, _ := time.Parse(time.RFC3339, templates[i].DtCreated)
			time2, _ := time.Parse(time.RFC3339, templates[j].DtCreated)
			return time1.Before(time2)
		})
		return nil
	default:
		return fmt.Errorf("Invalid sort option: %s", sortBy)
	}
}
