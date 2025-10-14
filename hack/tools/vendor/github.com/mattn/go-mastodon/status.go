package mastodon

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Status is struct to hold status.
type Status struct {
	ID                 ID              `json:"id"`
	URI                string          `json:"uri"`
	URL                string          `json:"url"`
	Account            Account         `json:"account"`
	InReplyToID        interface{}     `json:"in_reply_to_id"`
	InReplyToAccountID interface{}     `json:"in_reply_to_account_id"`
	Reblog             *Status         `json:"reblog"`
	Content            string          `json:"content"`
	CreatedAt          time.Time       `json:"created_at"`
	EditedAt           time.Time       `json:"edited_at"`
	Emojis             []Emoji         `json:"emojis"`
	RepliesCount       int64           `json:"replies_count"`
	ReblogsCount       int64           `json:"reblogs_count"`
	FavouritesCount    int64           `json:"favourites_count"`
	Reblogged          interface{}     `json:"reblogged"`
	Favourited         interface{}     `json:"favourited"`
	Bookmarked         interface{}     `json:"bookmarked"`
	Muted              interface{}     `json:"muted"`
	Sensitive          bool            `json:"sensitive"`
	SpoilerText        string          `json:"spoiler_text"`
	Visibility         string          `json:"visibility"`
	MediaAttachments   []Attachment    `json:"media_attachments"`
	Mentions           []Mention       `json:"mentions"`
	Tags               []Tag           `json:"tags"`
	Card               *Card           `json:"card"`
	Poll               *Poll           `json:"poll"`
	Application        Application     `json:"application"`
	Language           string          `json:"language"`
	Pinned             interface{}     `json:"pinned"`
	ScheduledParams    ScheduledParams `json:"params"`
}

// StatusHistory is a struct to hold status history data.
type StatusHistory struct {
	Content          string       `json:"content"`
	SpoilerText      string       `json:"spoiler_text"`
	Account          Account      `json:"account"`
	Sensitive        bool         `json:"sensitive"`
	CreatedAt        time.Time    `json:"created_at"`
	Emojis           []Emoji      `json:"emojis"`
	MediaAttachments []Attachment `json:"media_attachments"`
}

// ScheduledStatus holds information returned when ScheduledAt is set on a status
type ScheduledParams struct {
	ApplicationID ID          `json:"application_id"`
	Idempotency   string      `json:"idempotency"`
	InReplyToID   interface{} `json:"in_reply_to_id"`
	MediaIDs      []ID        `json:"media_ids"`
	Poll          *Poll       `json:"poll"`
	ScheduledAt   *time.Time  `json:"scheduled_at,omitempty"`
	Sensitive     bool        `json:"sensitive"`
	SpoilerText   string      `json:"spoiler_text"`
	Text          string      `json:"text"`
	Visibility    string      `json:"visibility"`
}

// Context holds information for a mastodon context.
type Context struct {
	Ancestors   []*Status `json:"ancestors"`
	Descendants []*Status `json:"descendants"`
}

// Card holds information for a mastodon card.
type Card struct {
	URL          string `json:"url"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Image        string `json:"image"`
	Type         string `json:"type"`
	AuthorName   string `json:"author_name"`
	AuthorURL    string `json:"author_url"`
	ProviderName string `json:"provider_name"`
	ProviderURL  string `json:"provider_url"`
	HTML         string `json:"html"`
	Width        int64  `json:"width"`
	Height       int64  `json:"height"`
}

// Source holds source properties so a status can be edited.
type Source struct {
	ID          ID     `json:"id"`
	Text        string `json:"text"`
	SpoilerText string `json:"spoiler_text"`
}

// Conversation holds information for a mastodon conversation.
type Conversation struct {
	ID         ID         `json:"id"`
	Accounts   []*Account `json:"accounts"`
	Unread     bool       `json:"unread"`
	LastStatus *Status    `json:"last_status"`
}

// Media is struct to hold media.
type Media struct {
	File        io.Reader
	Thumbnail   io.Reader
	Description string
	Focus       string
}

func (m *Media) bodyAndContentType() (io.Reader, string, error) {
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)

	fileName := "upload"
	if f, ok := m.File.(*os.File); ok {
		fileName = f.Name()
	}
	file, err := mw.CreateFormFile("file", fileName)
	if err != nil {
		return nil, "", err
	}
	if _, err := io.Copy(file, m.File); err != nil {
		return nil, "", err
	}

	if m.Thumbnail != nil {
		thumbName := "upload"
		if f, ok := m.Thumbnail.(*os.File); ok {
			thumbName = f.Name()
		}
		thumb, err := mw.CreateFormFile("thumbnail", thumbName)
		if err != nil {
			return nil, "", err
		}
		if _, err := io.Copy(thumb, m.Thumbnail); err != nil {
			return nil, "", err
		}
	}

	if m.Description != "" {
		desc, err := mw.CreateFormField("description")
		if err != nil {
			return nil, "", err
		}
		if _, err := io.Copy(desc, strings.NewReader(m.Description)); err != nil {
			return nil, "", err
		}
	}

	if m.Focus != "" {
		focus, err := mw.CreateFormField("focus")
		if err != nil {
			return nil, "", err
		}
		if _, err := io.Copy(focus, strings.NewReader(m.Focus)); err != nil {
			return nil, "", err
		}
	}

	if err := mw.Close(); err != nil {
		return nil, "", err
	}

	return &buf, mw.FormDataContentType(), nil
}

// GetFavourites returns the favorite list of the current user.
func (c *Client) GetFavourites(ctx context.Context, pg *Pagination) ([]*Status, error) {
	var statuses []*Status
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/favourites", nil, &statuses, pg)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// GetBookmarks returns the bookmark list of the current user.
func (c *Client) GetBookmarks(ctx context.Context, pg *Pagination) ([]*Status, error) {
	var statuses []*Status
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/bookmarks", nil, &statuses, pg)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// GetStatus returns status specified by id.
func (c *Client) GetStatus(ctx context.Context, id ID) (*Status, error) {
	var status Status
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/statuses/%s", id), nil, &status, nil)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// GetStatusContext returns status specified by id.
func (c *Client) GetStatusContext(ctx context.Context, id ID) (*Context, error) {
	var context Context
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/statuses/%s/context", id), nil, &context, nil)
	if err != nil {
		return nil, err
	}
	return &context, nil
}

// GetStatusCard returns status specified by id.
func (c *Client) GetStatusCard(ctx context.Context, id ID) (*Card, error) {
	var card Card
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/statuses/%s/card", id), nil, &card, nil)
	if err != nil {
		return nil, err
	}
	return &card, nil
}

// GetStatusSource returns source data specified by id.
func (c *Client) GetStatusSource(ctx context.Context, id ID) (*Source, error) {
	var source Source
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/statuses/%s/source", id), nil, &source, nil)
	if err != nil {
		return nil, err
	}
	return &source, nil
}

// GetStatusHistory returns the status history specified by id.
func (c *Client) GetStatusHistory(ctx context.Context, id ID) ([]*StatusHistory, error) {
	var statuses []*StatusHistory
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/statuses/%s/history", id), nil, &statuses, nil)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// GetRebloggedBy returns the account list of the user who reblogged the toot of id.
func (c *Client) GetRebloggedBy(ctx context.Context, id ID, pg *Pagination) ([]*Account, error) {
	var accounts []*Account
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/statuses/%s/reblogged_by", id), nil, &accounts, pg)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// GetFavouritedBy returns the account list of the user who liked the toot of id.
func (c *Client) GetFavouritedBy(ctx context.Context, id ID, pg *Pagination) ([]*Account, error) {
	var accounts []*Account
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/statuses/%s/favourited_by", id), nil, &accounts, pg)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// Reblog reblogs the toot of id and returns status of reblog.
func (c *Client) Reblog(ctx context.Context, id ID) (*Status, error) {
	var status Status
	err := c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/statuses/%s/reblog", id), nil, &status, nil)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// Unreblog unreblogs the toot of id and returns status of the original toot.
func (c *Client) Unreblog(ctx context.Context, id ID) (*Status, error) {
	var status Status
	err := c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/statuses/%s/unreblog", id), nil, &status, nil)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// Favourite favourites the toot of id and returns status of the favourite toot.
func (c *Client) Favourite(ctx context.Context, id ID) (*Status, error) {
	var status Status
	err := c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/statuses/%s/favourite", id), nil, &status, nil)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// Unfavourite unfavourites the toot of id and returns status of the unfavourite toot.
func (c *Client) Unfavourite(ctx context.Context, id ID) (*Status, error) {
	var status Status
	err := c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/statuses/%s/unfavourite", id), nil, &status, nil)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// Bookmark bookmarks the toot of id and returns status of the bookmark toot.
func (c *Client) Bookmark(ctx context.Context, id ID) (*Status, error) {
	var status Status
	err := c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/statuses/%s/bookmark", id), nil, &status, nil)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// Unbookmark is unbookmark the toot of id and return status of the unbookmark toot.
func (c *Client) Unbookmark(ctx context.Context, id ID) (*Status, error) {
	var status Status
	err := c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/statuses/%s/unbookmark", id), nil, &status, nil)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// GetTimelineHome return statuses from home timeline.
func (c *Client) GetTimelineHome(ctx context.Context, pg *Pagination) ([]*Status, error) {
	var statuses []*Status
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/timelines/home", nil, &statuses, pg)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// GetTimelinePublic return statuses from public timeline.
func (c *Client) GetTimelinePublic(ctx context.Context, isLocal bool, pg *Pagination) ([]*Status, error) {
	params := url.Values{}
	if isLocal {
		params.Set("local", "t")
	}

	var statuses []*Status
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/timelines/public", params, &statuses, pg)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// GetTimelineHashtag return statuses from tagged timeline.
func (c *Client) GetTimelineHashtag(ctx context.Context, tag string, isLocal bool, pg *Pagination) ([]*Status, error) {
	params := url.Values{}
	if isLocal {
		params.Set("local", "t")
	}

	var statuses []*Status
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/timelines/tag/%s", url.PathEscape(tag)), params, &statuses, pg)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// GetTimelineList return statuses from a list timeline.
func (c *Client) GetTimelineList(ctx context.Context, id ID, pg *Pagination) ([]*Status, error) {
	var statuses []*Status
	err := c.doAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/timelines/list/%s", url.PathEscape(string(id))), nil, &statuses, pg)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// GetTimelineMedia return statuses from media timeline.
// NOTE: This is an experimental feature of pawoo.net.
func (c *Client) GetTimelineMedia(ctx context.Context, isLocal bool, pg *Pagination) ([]*Status, error) {
	params := url.Values{}
	params.Set("media", "t")
	if isLocal {
		params.Set("local", "t")
	}

	var statuses []*Status
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/timelines/public", params, &statuses, pg)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// PostStatus post the toot.
func (c *Client) PostStatus(ctx context.Context, toot *Toot) (*Status, error) {
	return c.postStatus(ctx, toot, false, ID("none"))
}

// UpdateStatus updates the toot.
func (c *Client) UpdateStatus(ctx context.Context, toot *Toot, id ID) (*Status, error) {
	return c.postStatus(ctx, toot, true, id)
}

func (c *Client) postStatus(ctx context.Context, toot *Toot, update bool, updateID ID) (*Status, error) {
	params := url.Values{}
	params.Set("status", toot.Status)
	if toot.InReplyToID != "" {
		params.Set("in_reply_to_id", string(toot.InReplyToID))
	}
	if toot.MediaIDs != nil {
		for _, media := range toot.MediaIDs {
			params.Add("media_ids[]", string(media))
		}
	}
	// Can't use Media and Poll at the same time.
	if toot.Poll != nil && toot.Poll.Options != nil && toot.MediaIDs == nil {
		for _, opt := range toot.Poll.Options {
			params.Add("poll[options][]", string(opt))
		}
		params.Add("poll[expires_in]", fmt.Sprintf("%d", toot.Poll.ExpiresInSeconds))
		if toot.Poll.Multiple {
			params.Add("poll[multiple]", "true")
		}
		if toot.Poll.HideTotals {
			params.Add("poll[hide_totals]", "true")
		}
	}
	if toot.Visibility != "" {
		params.Set("visibility", fmt.Sprint(toot.Visibility))
	}
	if toot.Language != "" {
		params.Set("language", fmt.Sprint(toot.Language))
	}
	if toot.Sensitive {
		params.Set("sensitive", "true")
	}
	if toot.SpoilerText != "" {
		params.Set("spoiler_text", toot.SpoilerText)
	}
	if toot.ScheduledAt != nil {
		params.Set("scheduled_at", toot.ScheduledAt.Format(time.RFC3339))
	}

	var status Status
	var err error
	if !update {
		err = c.doAPI(ctx, http.MethodPost, "/api/v1/statuses", params, &status, nil)
	} else {
		err = c.doAPI(ctx, http.MethodPut, fmt.Sprintf("/api/v1/statuses/%s", updateID), params, &status, nil)
	}
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// DeleteStatus delete the toot.
func (c *Client) DeleteStatus(ctx context.Context, id ID) error {
	return c.doAPI(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/statuses/%s", id), nil, nil, nil)
}

// Search search content with query.
func (c *Client) Search(ctx context.Context, q string, resolve bool) (*Results, error) {
	params := url.Values{}
	params.Set("q", q)
	params.Set("resolve", fmt.Sprint(resolve))
	var results Results
	err := c.doAPI(ctx, http.MethodGet, "/api/v2/search", params, &results, nil)
	if err != nil {
		return nil, err
	}
	return &results, nil
}

// UploadMedia upload a media attachment from a file.
func (c *Client) UploadMedia(ctx context.Context, file string) (*Attachment, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return c.UploadMediaFromMedia(ctx, &Media{File: f})
}

// UploadMediaFromBytes uploads a media attachment from a byte slice.
func (c *Client) UploadMediaFromBytes(ctx context.Context, b []byte) (*Attachment, error) {
	return c.UploadMediaFromReader(ctx, bytes.NewReader(b))
}

// UploadMediaFromReader uploads a media attachment from an io.Reader.
func (c *Client) UploadMediaFromReader(ctx context.Context, reader io.Reader) (*Attachment, error) {
	return c.UploadMediaFromMedia(ctx, &Media{File: reader})
}

// UploadMediaFromMedia uploads a media attachment from a Media struct.
func (c *Client) UploadMediaFromMedia(ctx context.Context, media *Media) (*Attachment, error) {
	var attachment Attachment
	if err := c.doAPI(ctx, http.MethodPost, "/api/v1/media", media, &attachment, nil); err != nil {
		return nil, err
	}
	return &attachment, nil
}

// GetTimelineDirect return statuses from direct timeline.
func (c *Client) GetTimelineDirect(ctx context.Context, pg *Pagination) ([]*Status, error) {
	params := url.Values{}

	var conversations []*Conversation
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/conversations", params, &conversations, pg)
	if err != nil {
		return nil, err
	}

	var statuses = []*Status{}

	for _, c := range conversations {
		s := c.LastStatus
		statuses = append(statuses, s)
	}

	return statuses, nil
}

// GetConversations return direct conversations.
func (c *Client) GetConversations(ctx context.Context, pg *Pagination) ([]*Conversation, error) {
	params := url.Values{}

	var conversations []*Conversation
	err := c.doAPI(ctx, http.MethodGet, "/api/v1/conversations", params, &conversations, pg)
	if err != nil {
		return nil, err
	}
	return conversations, nil
}

// DeleteConversation delete the conversation specified by id.
func (c *Client) DeleteConversation(ctx context.Context, id ID) error {
	return c.doAPI(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/conversations/%s", id), nil, nil, nil)
}

// MarkConversationAsRead mark the conversation as read.
func (c *Client) MarkConversationAsRead(ctx context.Context, id ID) error {
	return c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/conversations/%s/read", id), nil, nil, nil)
}
