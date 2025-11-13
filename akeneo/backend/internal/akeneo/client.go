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

func (c *Client) UpdateSubscriber(subscriberID string, subscriber map[string]interface{}) error {
	_, err := c.eventPlatformRequest("PATCH", fmt.Sprintf("/api/v1/subscribers/%s", subscriberID), subscriber)
	return err
}

func (c *Client) DeleteSubscriber(subscriberID string) error {
	_, err := c.eventPlatformRequest("DELETE", fmt.Sprintf("/api/v1/subscribers/%s", subscriberID), nil)
	return err
}

func (c *Client) GetSubscriptions() ([]map[string]interface{}, error) {
	// Get all subscribers first
	subscribers, err := c.GetSubscriber()
	if err != nil {
		return nil, fmt.Errorf("failed to get subscribers: %w", err)
	}

	// Collect all subscriptions from all subscribers
	var allSubscriptions []map[string]interface{}
	for _, subscriber := range subscribers {
		subscriberID, ok := subscriber["id"].(string)
		if !ok {
			continue
		}

		data, err := c.eventPlatformRequest("GET", fmt.Sprintf("/api/v1/subscribers/%s/subscriptions", subscriberID), nil)
		if err != nil {
			// Skip if this subscriber has no subscriptions or error
			continue
		}

		var subscriptions []map[string]interface{}
		if err := json.Unmarshal(data, &subscriptions); err != nil {
			continue
		}

		// Add subscriber info to each subscription
		for i := range subscriptions {
			subscriptions[i]["subscriber_id"] = subscriberID
			subscriptions[i]["subscriber_name"] = subscriber["name"]
		}

		allSubscriptions = append(allSubscriptions, subscriptions...)
	}

	return allSubscriptions, nil
}

func (c *Client) CreateSubscription(subscription map[string]interface{}) error {
	// Get subscriber ID from the subscription payload
	subscriberID, ok := subscription["subscriber_id"].(string)
	if !ok || subscriberID == "" {
		return fmt.Errorf("subscriber_id is required")
	}
	
	// Remove subscriber_id from payload as it's in the URL
	delete(subscription, "subscriber_id")
	
	// Set subject to PIM URL if not provided
	if subscription["subject"] == nil || subscription["subject"] == "" {
		subscription["subject"] = c.BaseURL
	}
	
	_, err := c.eventPlatformRequest("POST", fmt.Sprintf("/api/v1/subscribers/%s/subscriptions", subscriberID), subscription)
	return err
}

func (c *Client) UpdateSubscription(subscriptionID string, subscription map[string]interface{}) error {
	// Get subscriber ID from the subscription payload
	subscriberID, ok := subscription["subscriber_id"].(string)
	if !ok || subscriberID == "" {
		return fmt.Errorf("subscriber_id is required in subscription payload")
	}

	// Remove subscriber_id from payload as it's in the URL
	delete(subscription, "subscriber_id")

	_, err := c.eventPlatformRequest("PATCH", fmt.Sprintf("/api/v1/subscribers/%s/subscriptions/%s", subscriberID, subscriptionID), subscription)
	return err
}

func (c *Client) DeleteSubscription(subscriptionID string) error {
	// We need to find which subscriber owns this subscription
	// For now, we'll try to get all subscribers and find the right one
	subscribers, err := c.GetSubscriber()
	if err != nil {
		return fmt.Errorf("failed to get subscribers: %w", err)
	}

	// Try to delete from each subscriber until we find the right one
	for _, subscriber := range subscribers {
		subscriberID, ok := subscriber["id"].(string)
		if !ok {
			continue
		}

		_, err := c.eventPlatformRequest("DELETE", fmt.Sprintf("/api/v1/subscribers/%s/subscriptions/%s", subscriberID, subscriptionID), nil)
		if err == nil {
			return nil // Successfully deleted
		}
	}

	return fmt.Errorf("subscription not found in any subscriber")
}

func (c *Client) GetEventTypes() []string {
	return GetEventTypes()
}
