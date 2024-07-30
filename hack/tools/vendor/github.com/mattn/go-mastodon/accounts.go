package mastodon

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Account holds information for a mastodon account.
type Account struct {
	ID             ID             `json:"id"`
	Username       string         `json:"username"`
	Acct           string         `json:"acct"`
	DisplayName    string         `json:"display_name"`
	Locked         bool           `json:"locked"`
	CreatedAt      time.Time      `json:"created_at"`
	FollowersCount int64          `json:"followers_count"`
	FollowingCount int64          `json:"following_count"`
	StatusesCount  int64          `json:"statuses_count"`
	Note           string         `json:"note"`
	URL            string         `json:"url"`
	Avatar         string         `json:"avatar"`
	AvatarStatic   string         `json:"avatar_static"`
	Header         string         `json:"header"`
	HeaderStatic   string         `json:"header_static"`
	Emojis         []Emoji        `json:"emojis"`
	Moved          *Account       `json:"moved"`
	Fields         []Field        `json:"fields"`
	Bot            bool           `json:"bot"`
	Discoverable   bool           `json:"discoverable"`
	Source         *AccountSource `json:"source"`
}

// Field is a Mastodon account profile field.
type Field struct {
	Name       string    `json:"name"`
	Value      string    `json:"value"`
	VerifiedAt time.Time `json:"verified_at"`
}

// AccountSource is a Mastodon account profile field.
type AccountSource struct {
	Privacy   *string  `json:"privacy"`
	Sensitive *bool    `json:"sensitive"`
	Language  *string  `json:"language"`
	Note      *string  `json:"note"`
	Fields    *[]Field `json:"fields"`
}

// GetAccount return Account.
func (c *Client) GetAccount(ctx context.Context, id ID) (*Account, error) {
	var account Account
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/accounts/%s", url.PathEscape(string(id))), nil, &account, nil)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetAccountCurrentUser returns the Account of current user.
func (c *Client) GetAccountCurrentUser(ctx context.Context) (*Account, error) {
	var account Account
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/accounts/verify_credentials", nil, &account, nil)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// Profile is a struct for updating profiles.
type Profile struct {
	// If it is nil it will not be updated.
	// If it is empty, update it with empty.
	DisplayName *string
	Note        *string
	Locked      *bool
	Fields      *[]Field
	Source      *AccountSource

	// Set the base64 encoded character string of the image.
	Avatar string
	Header string
}

// AccountUpdate updates the information of the current user.
func (c *Client) AccountUpdate(ctx context.Context, profile *Profile) (*Account, error) {
	params := url.Values{}
	if profile.DisplayName != nil {
		params.Set("display_name", *profile.DisplayName)
	}
	if profile.Note != nil {
		params.Set("note", *profile.Note)
	}
	if profile.Locked != nil {
		params.Set("locked", strconv.FormatBool(*profile.Locked))
	}
	if profile.Fields != nil {
		for idx, field := range *profile.Fields {
			params.Set(fmt.Sprintf("fields_attributes[%d][name]", idx), field.Name)
			params.Set(fmt.Sprintf("fields_attributes[%d][value]", idx), field.Value)
		}
	}
	if profile.Source != nil {
		if profile.Source.Privacy != nil {
			params.Set("source[privacy]", *profile.Source.Privacy)
		}
		if profile.Source.Sensitive != nil {
			params.Set("source[sensitive]", strconv.FormatBool(*profile.Source.Sensitive))
		}
		if profile.Source.Language != nil {
			params.Set("source[language]", *profile.Source.Language)
		}
	}
	if profile.Avatar != "" {
		params.Set("avatar", profile.Avatar)
	}
	if profile.Header != "" {
		params.Set("header", profile.Header)
	}

	var account Account
	err := c.doAPI(ctx, http.MethodPatch, "/api/v1/accounts/update_credentials", params, &account, nil)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetAccountStatuses return statuses by specified account.
func (c *Client) GetAccountStatuses(ctx context.Context, id ID, pg *Pagination) ([]*Status, error) {
	var statuses []*Status
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/accounts/%s/statuses", url.PathEscape(string(id))), nil, &statuses, pg)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// GetAccountPinnedStatuses returns statuses pinned by specified accuont.
func (c *Client) GetAccountPinnedStatuses(ctx context.Context, id ID) ([]*Status, error) {
	var statuses []*Status
	params := url.Values{}
	params.Set("pinned", "true")
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/accounts/%s/statuses", url.PathEscape(string(id))), params, &statuses, nil)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// GetAccountFollowers returns followers list.
func (c *Client) GetAccountFollowers(ctx context.Context, id ID, pg *Pagination) ([]*Account, error) {
	var accounts []*Account
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/accounts/%s/followers", url.PathEscape(string(id))), nil, &accounts, pg)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// GetAccountFollowing returns following list.
func (c *Client) GetAccountFollowing(ctx context.Context, id ID, pg *Pagination) ([]*Account, error) {
	var accounts []*Account
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/accounts/%s/following", url.PathEscape(string(id))), nil, &accounts, pg)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// GetBlocks returns block list.
func (c *Client) GetBlocks(ctx context.Context, pg *Pagination) ([]*Account, error) {
	var accounts []*Account
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/blocks", nil, &accounts, pg)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// Relationship holds information for relationship to the account.
type Relationship struct {
	ID                  ID   `json:"id"`
	Following           bool `json:"following"`
	FollowedBy          bool `json:"followed_by"`
	Blocking            bool `json:"blocking"`
	Muting              bool `json:"muting"`
	MutingNotifications bool `json:"muting_notifications"`
	Requested           bool `json:"requested"`
	DomainBlocking      bool `json:"domain_blocking"`
	ShowingReblogs      bool `json:"showing_reblogs"`
	Endorsed            bool `json:"endorsed"`
}

// AccountFollow follows the account.
func (c *Client) AccountFollow(ctx context.Context, id ID) (*Relationship, error) {
	var relationship Relationship
	err := c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/accounts/%s/follow", url.PathEscape(string(id))), nil, &relationship, nil)
	if err != nil {
		return nil, err
	}
	return &relationship, nil
}

// AccountUnfollow unfollows the account.
func (c *Client) AccountUnfollow(ctx context.Context, id ID) (*Relationship, error) {
	var relationship Relationship
	err := c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/accounts/%s/unfollow", url.PathEscape(string(id))), nil, &relationship, nil)
	if err != nil {
		return nil, err
	}
	return &relationship, nil
}

// AccountBlock blocks the account.
func (c *Client) AccountBlock(ctx context.Context, id ID) (*Relationship, error) {
	var relationship Relationship
	err := c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/accounts/%s/block", url.PathEscape(string(id))), nil, &relationship, nil)
	if err != nil {
		return nil, err
	}
	return &relationship, nil
}

// AccountUnblock unblocks the account.
func (c *Client) AccountUnblock(ctx context.Context, id ID) (*Relationship, error) {
	var relationship Relationship
	err := c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/accounts/%s/unblock", url.PathEscape(string(id))), nil, &relationship, nil)
	if err != nil {
		return nil, err
	}
	return &relationship, nil
}

// AccountMute mutes the account.
func (c *Client) AccountMute(ctx context.Context, id ID) (*Relationship, error) {
	var relationship Relationship
	err := c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/accounts/%s/mute", url.PathEscape(string(id))), nil, &relationship, nil)
	if err != nil {
		return nil, err
	}
	return &relationship, nil
}

// AccountUnmute unmutes the account.
func (c *Client) AccountUnmute(ctx context.Context, id ID) (*Relationship, error) {
	var relationship Relationship
	err := c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/accounts/%s/unmute", url.PathEscape(string(id))), nil, &relationship, nil)
	if err != nil {
		return nil, err
	}
	return &relationship, nil
}

// GetAccountRelationships returns relationship for the account.
func (c *Client) GetAccountRelationships(ctx context.Context, ids []string) ([]*Relationship, error) {
	params := url.Values{}
	for _, id := range ids {
		params.Add("id[]", id)
	}

	var relationships []*Relationship
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/accounts/relationships", params, &relationships, nil)
	if err != nil {
		return nil, err
	}
	return relationships, nil
}

// AccountsSearch searches accounts by query.
func (c *Client) AccountsSearch(ctx context.Context, q string, limit int64) ([]*Account, error) {
	params := url.Values{}
	params.Set("q", q)
	params.Set("limit", fmt.Sprint(limit))

	var accounts []*Account
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/accounts/search", params, &accounts, nil)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// FollowRemoteUser sends follow-request.
func (c *Client) FollowRemoteUser(ctx context.Context, uri string) (*Account, error) {
	params := url.Values{}
	params.Set("uri", uri)

	var account Account
	err := c.doAPI(ctx, http.MethodPost, "/api/v1/follows", params, &account, nil)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetFollowRequests returns follow requests.
func (c *Client) GetFollowRequests(ctx context.Context, pg *Pagination) ([]*Account, error) {
	var accounts []*Account
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/follow_requests", nil, &accounts, pg)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// FollowRequestAuthorize authorizes the follow request of user with id.
func (c *Client) FollowRequestAuthorize(ctx context.Context, id ID) error {
	return c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/follow_requests/%s/authorize", url.PathEscape(string(id))), nil, nil, nil)
}

// FollowRequestReject rejects the follow request of user with id.
func (c *Client) FollowRequestReject(ctx context.Context, id ID) error {
	return c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/follow_requests/%s/reject", url.PathEscape(string(id))), nil, nil, nil)
}

// GetMutes returns the list of users muted by the current user.
func (c *Client) GetMutes(ctx context.Context, pg *Pagination) ([]*Account, error) {
	var accounts []*Account
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/mutes", nil, &accounts, pg)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
