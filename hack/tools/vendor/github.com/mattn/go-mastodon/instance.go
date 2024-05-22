package mastodon

import (
	"context"
	"net/http"
)

// Instance holds information for a mastodon instance.
type Instance struct {
	URI            string            `json:"uri"`
	Title          string            `json:"title"`
	Description    string            `json:"description"`
	EMail          string            `json:"email"`
	Version        string            `json:"version,omitempty"`
	Thumbnail      string            `json:"thumbnail,omitempty"`
	URLs           map[string]string `json:"urls,omitempty"`
	Stats          *InstanceStats    `json:"stats,omitempty"`
	Languages      []string          `json:"languages"`
	ContactAccount *Account          `json:"contact_account"`
}

// InstanceStats holds information for mastodon instance stats.
type InstanceStats struct {
	UserCount   int64 `json:"user_count"`
	StatusCount int64 `json:"status_count"`
	DomainCount int64 `json:"domain_count"`
}

// GetInstance returns Instance.
func (c *Client) GetInstance(ctx context.Context) (*Instance, error) {
	var instance Instance
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/instance", nil, &instance, nil)
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

// WeeklyActivity holds information for mastodon weekly activity.
type WeeklyActivity struct {
	Week          Unixtime `json:"week"`
	Statuses      int64    `json:"statuses,string"`
	Logins        int64    `json:"logins,string"`
	Registrations int64    `json:"registrations,string"`
}

// GetInstanceActivity returns instance activity.
func (c *Client) GetInstanceActivity(ctx context.Context) ([]*WeeklyActivity, error) {
	var activity []*WeeklyActivity
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/instance/activity", nil, &activity, nil)
	if err != nil {
		return nil, err
	}
	return activity, nil
}

// GetInstancePeers returns instance peers.
func (c *Client) GetInstancePeers(ctx context.Context) ([]string, error) {
	var peers []string
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/instance/peers", nil, &peers, nil)
	if err != nil {
		return nil, err
	}
	return peers, nil
}
