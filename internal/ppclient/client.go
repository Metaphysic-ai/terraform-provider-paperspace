package ppclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"terraform-provider-paperspace/internal/ppclient/models"
)

const HostURL string = "https://api.paperspace.com/v1"

type Client struct {
	HostURL     string
	HTTPClient  *http.Client
	Token       string
	AuthSession *models.AuthSession
}

func (c *Client) GetAuthSession() (*models.AuthSession, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/auth/session", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	authSession := &models.AuthSession{}
	err = json.Unmarshal(body, authSession)

	if err != nil {
		return nil, err
	}

	if authSession.User.ID == "" {
		return nil, errors.New("Unable to get auth session, possibly invalid token")
	}

	return authSession, nil
}

func NewClient(host, authToken *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    HostURL,
	}

	// If authToken not provided, return empty client
	if authToken == nil {
		return &c, nil
	}

	if host != nil {
		c.HostURL = *host
	}

	c.Token = *authToken

	authSession, err := c.GetAuthSession()
	if err != nil {
		return nil, err
	}

	c.AuthSession = authSession

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", "Bearer "+c.Token)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
