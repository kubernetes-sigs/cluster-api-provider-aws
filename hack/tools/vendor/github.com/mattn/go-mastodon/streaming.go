package mastodon

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// UpdateEvent is a struct for passing status event to app.
type UpdateEvent struct {
	Status *Status `json:"status"`
}

func (e *UpdateEvent) event() {}

// UpdateEditEvent is a struct for passing status edit event to app.
type UpdateEditEvent struct {
	Status *Status `json:"status"`
}

func (e *UpdateEditEvent) event() {}

// NotificationEvent is a struct for passing notification event to app.
type NotificationEvent struct {
	Notification *Notification `json:"notification"`
}

func (e *NotificationEvent) event() {}

// DeleteEvent is a struct for passing deletion event to app.
type DeleteEvent struct{ ID ID }

func (e *DeleteEvent) event() {}

// ErrorEvent is a struct for passing errors to app.
type ErrorEvent struct{ Err error }

func (e *ErrorEvent) event()        {}
func (e *ErrorEvent) Error() string { return e.Err.Error() }

// Event is an interface passing events to app.
type Event interface {
	event()
}

func handleReader(q chan Event, r io.Reader) error {
	var name string
	var lineBuf bytes.Buffer
	br := bufio.NewReader(r)
	for {
		line, isPrefix, err := br.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		if isPrefix {
			lineBuf.Write(line)
			continue
		}
		if lineBuf.Len() > 0 {
			lineBuf.Write(line)
			line = lineBuf.Bytes()
			lineBuf.Reset()
		}

		token := strings.SplitN(string(line), ":", 2)
		if len(token) != 2 {
			continue
		}
		switch strings.TrimSpace(token[0]) {
		case "event":
			name = strings.TrimSpace(token[1])
		case "data":
			var err error
			switch name {
			case "update":
				var status Status
				err = json.Unmarshal([]byte(token[1]), &status)
				if err == nil {
					q <- &UpdateEvent{&status}
				}
			case "status.update":
				var status Status
				err = json.Unmarshal([]byte(token[1]), &status)
				if err == nil {
					q <- &UpdateEditEvent{&status}
				}
			case "notification":
				var notification Notification
				err = json.Unmarshal([]byte(token[1]), &notification)
				if err == nil {
					q <- &NotificationEvent{&notification}
				}
			case "delete":
				q <- &DeleteEvent{ID: ID(strings.TrimSpace(token[1]))}
			}
			if err != nil {
				q <- &ErrorEvent{err}
			}
		}
	}
}

func (c *Client) streaming(ctx context.Context, p string, params url.Values) (chan Event, error) {
	u, err := url.Parse(c.Config.Server)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, "/api/v1/streaming", p)
	u.RawQuery = params.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	if c.Config.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.Config.AccessToken)
	}

	q := make(chan Event)
	go func() {
		defer close(q)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			c.doStreaming(req, q)
		}
	}()
	return q, nil
}

func (c *Client) doStreaming(req *http.Request, q chan Event) {
	resp, err := c.Do(req)
	if err != nil {
		q <- &ErrorEvent{err}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		q <- &ErrorEvent{parseAPIError("bad request", resp)}
		return
	}

	err = handleReader(q, resp.Body)
	if err != nil {
		q <- &ErrorEvent{err}
	}
}

// StreamingUser returns a channel to read events on home.
func (c *Client) StreamingUser(ctx context.Context) (chan Event, error) {
	return c.streaming(ctx, "user", nil)
}

// StreamingPublic returns a channel to read events on public.
func (c *Client) StreamingPublic(ctx context.Context, isLocal bool) (chan Event, error) {
	p := "public"
	if isLocal {
		p = path.Join(p, "local")
	}

	return c.streaming(ctx, p, nil)
}

// StreamingHashtag returns a channel to read events on tagged timeline.
func (c *Client) StreamingHashtag(ctx context.Context, tag string, isLocal bool) (chan Event, error) {
	params := url.Values{}
	params.Set("tag", tag)

	p := "hashtag"
	if isLocal {
		p = path.Join(p, "local")
	}

	return c.streaming(ctx, p, params)
}

// StreamingList returns a channel to read events on a list.
func (c *Client) StreamingList(ctx context.Context, id ID) (chan Event, error) {
	params := url.Values{}
	params.Set("list", string(id))

	return c.streaming(ctx, "list", params)
}

// StreamingDirect returns a channel to read events on a direct messages.
func (c *Client) StreamingDirect(ctx context.Context) (chan Event, error) {
	return c.streaming(ctx, "direct", nil)
}
