package mastodon

import (
	"context"
	"fmt"
	"net/http"
)

// TagUnfollow unfollows a hashtag.
func (c *Client) TagUnfollow(ctx context.Context, ID string) (*FollowedTag, error) {
	var tag FollowedTag
	err := c.doAPI(ctx, http.MethodPost, fmt.Sprintf("/api/v1/tags/%s/unfollow", ID), nil, &tag, nil)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}
