package ppclient

func (c *Client) GetCustomTemplates() ([]CustomTemplate, error) {
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

	return allItems, nil
}
