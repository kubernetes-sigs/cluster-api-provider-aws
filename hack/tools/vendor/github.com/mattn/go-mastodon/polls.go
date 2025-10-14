package mastodon

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Poll holds information for mastodon polls.
type Poll struct {
	ID          ID           `json:"id"`
	ExpiresAt   time.Time    `json:"expires_at"`
	Expired     bool         `json:"expired"`
	Multiple    bool         `json:"multiple"`
	VotesCount  int64        `json:"votes_count"`
	VotersCount int64        `json:"voters_count"`
	Options     []PollOption `json:"options"`
	Voted       bool         `json:"voted"`
	OwnVotes    []int        `json:"own_votes"`
	Emojis      []Emoji      `json:"emojis"`
}

// Poll holds information for a mastodon poll option.
type PollOption struct {
	Title      string `json:"title"`
	VotesCount int64  `json:"votes_count"`
}

// GetPoll returns poll specified by id.
func (c *Client) GetPoll(ctx context.Context, id ID) (*Poll, error) {
	var poll Poll
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/polls/%s", id), nil, &poll, nil)
	if err != nil {
		return nil, err
	}
	return &poll, nil
}

// PollVote votes on a poll specified by id, choices is the Poll.Options index to vote on
func (c *Client) PollVote(ctx context.Context, id ID, choices ...int) (*Poll, error) {
	params := url.Values{}
	for _, c := range choices {
		params.Add("choices[]", fmt.Sprintf("%d", c))
	}

	var poll Poll
	err := c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/polls/%s/votes", url.PathEscape(string(id))), params, &poll, nil)
	if err != nil {
		return nil, err
	}
	return &poll, nil
}
