package akeneo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	BaseURL          string // PIM URL for authentication
	EventPlatformURL string // Event Platform URL
	ClientID         string
	ClientSecret     string
	Username         string
	Password         string
	accessToken      string
	subscriberID     string
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

func NewClient(baseURL, eventPlatformURL, clientID, clientSecret, username, password string) *Client {
	return &Client{
		BaseURL:          strings.TrimSuffix(baseURL, "/"),
		EventPlatformURL: strings.TrimSuffix(eventPlatformURL, "/"),
		ClientID:         clientID,
		ClientSecret:     clientSecret,
		Username:         username,
		Password:         password,
	}
}

func (c *Client) authenticate() error {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", c.Username)
	data.Set("password", c.Password)

	req, err := http.NewRequest("POST", c.BaseURL+"/api/oauth/v1/token", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.ClientID, c.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("authentication failed: %s", string(body))
	}

	var token tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return err
	}

	c.accessToken = token.AccessToken
	return nil
}

func (c *Client) request(method, path string, body interface{}) ([]byte, error) {
	if c.accessToken == "" {
		if err := c.authenticate(); err != nil {
			return nil, err
		}
	}

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (c *Client) eventPlatformRequest(method, path string, body interface{}) ([]byte, error) {
	if c.accessToken == "" {
		if err := c.authenticate(); err != nil {
			return nil, err
		}
	}

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.EventPlatformURL+path, reqBody)
	if err != nil {
		return nil, err
	}

	// Event Platform requires special headers as per documentation
	req.Header.Set("X-PIM-URL", c.BaseURL)
	req.Header.Set("X-PIM-TOKEN", c.accessToken)
	req.Header.Set("X-PIM-CLIENT-ID", c.ClientID)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Event Platform API error %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (c *Client) GetSubscriber() ([]map[string]interface{}, error) {
	data, err := c.eventPlatformRequest("GET", "/api/v1/subscribers", nil)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	// Store first subscriber ID for later use if available
	if len(result) > 0 {
		if id, ok := result[0]["id"].(string); ok {
			c.subscriberID = id
		}
	}

	return result, nil
}

func (c *Client) CreateSubscriber(subscriber map[string]interface{}) error {
	data, err := c.eventPlatformRequest("POST", "/api/v1/subscribers", subscriber)
	if err != nil {
		return err
	}

	// Extract and store subscriber ID from response
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err == nil {
		if id, ok := result["id"].(string); ok {
			c.subscriberID = id
		}
	}

	return nil
}

func (c *Client) UpdateSubscriber(subscriber map[string]interface{}) error {
	_, err := c.eventPlatformRequest("PATCH", "/api/v1/subscribers", subscriber)
	return err
}

func (c *Client) GetSubscriptions() ([]map[string]interface{}, error) {
	data, err := c.eventPlatformRequest("GET", "/api/v1/subscriptions", nil)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) CreateSubscription(subscription map[string]interface{}) error {
	// Ensure we have a subscriber ID
	if c.subscriberID == "" {
		// Try to get subscriber first
		_, err := c.GetSubscriber()
		if err != nil {
			return fmt.Errorf("failed to get subscriber ID: %w", err)
		}
	}

	// Add subscriber_id to the subscription payload
	subscription["subscriber_id"] = c.subscriberID

	_, err := c.eventPlatformRequest("POST", "/api/v1/subscriptions", subscription)
	return err
}

func (c *Client) UpdateSubscription(connectionCode string, subscription map[string]interface{}) error {
	_, err := c.eventPlatformRequest("PATCH", fmt.Sprintf("/api/v1/subscriptions/%s", connectionCode), subscription)
	return err
}

func (c *Client) DeleteSubscription(connectionCode string) error {
	_, err := c.eventPlatformRequest("DELETE", fmt.Sprintf("/api/v1/subscriptions/%s", connectionCode), nil)
	return err
}

func (c *Client) GetEventTypes() []string {
	return GetEventTypes()
}
