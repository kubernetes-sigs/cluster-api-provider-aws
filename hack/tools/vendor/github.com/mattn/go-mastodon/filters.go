package mastodon

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Filter is metadata for a filter of users.
type Filter struct {
	ID           ID        `json:"id"`
	Phrase       string    `json:"phrase"`
	Context      []string  `json:"context"`
	WholeWord    bool      `json:"whole_word"`
	ExpiresAt    time.Time `json:"expires_at"`
	Irreversible bool      `json:"irreversible"`
}

// GetFilters returns all the filters on the current account.
func (c *Client) GetFilters(ctx context.Context) ([]*Filter, error) {
	var filters []*Filter
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/filters", nil, &filters, nil)
	if err != nil {
		return nil, err
	}
	return filters, nil
}

// GetFilter retrieves a filter by ID.
func (c *Client) GetFilter(ctx context.Context, id ID) (*Filter, error) {
	var filter Filter
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/filters/%s", url.PathEscape(string(id))), nil, &filter, nil)
	if err != nil {
		return nil, err
	}
	return &filter, nil
}

// CreateFilter creates a new filter.
func (c *Client) CreateFilter(ctx context.Context, filter *Filter) (*Filter, error) {
	if filter == nil {
		return nil, errors.New("filter can't be nil")
	}
	if filter.Phrase == "" {
		return nil, errors.New("phrase can't be empty")
	}
	if len(filter.Context) == 0 {
		return nil, errors.New("context can't be empty")
	}
	params := url.Values{}
	params.Set("phrase", filter.Phrase)
	for _, c := range filter.Context {
		params.Add("context[]", c)
	}
	if filter.WholeWord {
		params.Add("whole_word", "true")
	}
	if filter.Irreversible {
		params.Add("irreversible", "true")
	}
	if !filter.ExpiresAt.IsZero() {
		diff := time.Until(filter.ExpiresAt)
		params.Add("expires_in", fmt.Sprintf("%.0f", diff.Seconds()))
	}

	var f Filter
	err := c.doAPI(ctx, http.MethodPost, "/api/v1/filters", params, &f, nil)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// UpdateFilter updates a filter.
func (c *Client) UpdateFilter(ctx context.Context, id ID, filter *Filter) (*Filter, error) {
	if filter == nil {
		return nil, errors.New("filter can't be nil")
	}
	if id == ID("") {
		return nil, errors.New("ID can't be empty")
	}
	if filter.Phrase == "" {
		return nil, errors.New("phrase can't be empty")
	}
	if len(filter.Context) == 0 {
		return nil, errors.New("context can't be empty")
	}
	params := url.Values{}
	params.Set("phrase", filter.Phrase)
	for _, c := range filter.Context {
		params.Add("context[]", c)
	}
	if filter.WholeWord {
		params.Add("whole_word", "true")
	} else {
		params.Add("whole_word", "false")
	}
	if filter.Irreversible {
		params.Add("irreversible", "true")
	} else {
		params.Add("irreversible", "false")
	}
	if !filter.ExpiresAt.IsZero() {
		diff := time.Until(filter.ExpiresAt)
		params.Add("expires_in", fmt.Sprintf("%.0f", diff.Seconds()))
	} else {
		params.Add("expires_in", "")
	}

	var f Filter
	err := c.doAPI(ctx, http.MethodPut, fmt.Sprintf("/api/v1/filters/%s", url.PathEscape(string(id))), params, &f, nil)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// DeleteFilter removes a filter.
func (c *Client) DeleteFilter(ctx context.Context, id ID) error {
	return c.doAPI(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/filters/%s", url.PathEscape(string(id))), nil, nil, nil)
}
