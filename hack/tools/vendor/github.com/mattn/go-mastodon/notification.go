package mastodon

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Notification holds information for a mastodon notification.
type Notification struct {
	ID        ID        `json:"id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	Account   Account   `json:"account"`
	Status    *Status   `json:"status"`
}

type PushSubscription struct {
	ID        ID          `json:"id"`
	Endpoint  string      `json:"endpoint"`
	ServerKey string      `json:"server_key"`
	Alerts    *PushAlerts `json:"alerts"`
}

type PushAlerts struct {
	Follow    *Sbool `json:"follow"`
	Favourite *Sbool `json:"favourite"`
	Reblog    *Sbool `json:"reblog"`
	Mention   *Sbool `json:"mention"`
}

// GetNotifications returns notifications.
func (c *Client) GetNotifications(ctx context.Context, pg *Pagination) ([]*Notification, error) {
	var notifications []*Notification
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/notifications", nil, &notifications, pg)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

// GetNotification returns notification.
func (c *Client) GetNotification(ctx context.Context, id ID) (*Notification, error) {
	var notification Notification
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/notifications/%v", id), nil, &notification, nil)
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

// DismissNotification deletes a single notification.
func (c *Client) DismissNotification(ctx context.Context, id ID) error {
	return c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/notifications/%v/dismiss", id), nil, nil, nil)
}

// ClearNotifications clears notifications.
func (c *Client) ClearNotifications(ctx context.Context) error {
	return c.doAPI(ctx, http.MethodPost, "/api/v1/notifications/clear", nil, nil, nil)
}

// AddPushSubscription adds a new push subscription.
func (c *Client) AddPushSubscription(ctx context.Context, endpoint string, public ecdsa.PublicKey, shared []byte, alerts PushAlerts) (*PushSubscription, error) {
	var subscription PushSubscription
	pk := elliptic.Marshal(public.Curve, public.X, public.Y)
	params := url.Values{}
	params.Add("subscription[endpoint]", endpoint)
	params.Add("subscription[keys][p256dh]", base64.RawURLEncoding.EncodeToString(pk))
	params.Add("subscription[keys][auth]", base64.RawURLEncoding.EncodeToString(shared))
	if alerts.Follow != nil {
		params.Add("data[alerts][follow]", strconv.FormatBool(bool(*alerts.Follow)))
	}
	if alerts.Favourite != nil {
		params.Add("data[alerts][favourite]", strconv.FormatBool(bool(*alerts.Favourite)))
	}
	if alerts.Reblog != nil {
		params.Add("data[alerts][reblog]", strconv.FormatBool(bool(*alerts.Reblog)))
	}
	if alerts.Mention != nil {
		params.Add("data[alerts][mention]", strconv.FormatBool(bool(*alerts.Mention)))
	}
	err := c.doAPI(ctx, http.MethodPost, "/api/v1/push/subscription", params, &subscription, nil)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// UpdatePushSubscription updates which type of notifications are sent for the active push subscription.
func (c *Client) UpdatePushSubscription(ctx context.Context, alerts *PushAlerts) (*PushSubscription, error) {
	var subscription PushSubscription
	params := url.Values{}
	if alerts.Follow != nil {
		params.Add("data[alerts][follow]", strconv.FormatBool(bool(*alerts.Follow)))
	}
	if alerts.Mention != nil {
		params.Add("data[alerts][favourite]", strconv.FormatBool(bool(*alerts.Favourite)))
	}
	if alerts.Reblog != nil {
		params.Add("data[alerts][reblog]", strconv.FormatBool(bool(*alerts.Reblog)))
	}
	if alerts.Mention != nil {
		params.Add("data[alerts][mention]", strconv.FormatBool(bool(*alerts.Mention)))
	}
	err := c.doAPI(ctx, http.MethodPut, "/api/v1/push/subscription", params, &subscription, nil)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// RemovePushSubscription deletes the active push subscription.
func (c *Client) RemovePushSubscription(ctx context.Context) error {
	return c.doAPI(ctx, http.MethodDelete, "/api/v1/push/subscription", nil, nil, nil)
}

// GetPushSubscription retrieves information about the active push subscription.
func (c *Client) GetPushSubscription(ctx context.Context) (*PushSubscription, error) {
	var subscription PushSubscription
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/push/subscription", nil, &subscription, nil)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}
