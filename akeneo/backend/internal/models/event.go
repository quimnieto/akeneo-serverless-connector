package models

type Subscriber struct {
	URL    string `json:"url" binding:"required"`
	Active bool   `json:"active"`
}

type Subscription struct {
	ConnectionCode string   `json:"connection_code" binding:"required"`
	Events         []string `json:"events" binding:"required"`
	Active         bool     `json:"active"`
}

type SubscriberResponse struct {
	URL    string `json:"url"`
	Active bool   `json:"active"`
}

type SubscriptionResponse struct {
	ConnectionCode string   `json:"connection_code"`
	Events         []string `json:"events"`
	Active         bool     `json:"active"`
}
