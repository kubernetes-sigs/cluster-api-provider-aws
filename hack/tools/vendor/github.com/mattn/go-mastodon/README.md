# go-mastodon

[![Build Status](https://github.com/mattn/go-mastodon/workflows/test/badge.svg?branch=master)](https://github.com/mattn/go-mastodon/actions?query=workflow%3Atest)
[![Codecov](https://codecov.io/gh/mattn/go-mastodon/branch/master/graph/badge.svg)](https://codecov.io/gh/mattn/go-mastodon)
[![Go Reference](https://pkg.go.dev/badge/github.com/mattn/go-mastodon.svg)](https://pkg.go.dev/github.com/mattn/go-mastodon)
[![Go Report Card](https://goreportcard.com/badge/github.com/mattn/go-mastodon)](https://goreportcard.com/report/github.com/mattn/go-mastodon)


## Usage

### Application

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mattn/go-mastodon"
)

func main() {
	app, err := mastodon.RegisterApp(context.Background(), &mastodon.AppConfig{
		Server:     "https://mstdn.jp",
		ClientName: "client-name",
		Scopes:     "read write follow",
		Website:    "https://github.com/mattn/go-mastodon",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("client-id    : %s\n", app.ClientID)
	fmt.Printf("client-secret: %s\n", app.ClientSecret)
}
```

### Client

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mattn/go-mastodon"
)

func main() {
	c := mastodon.NewClient(&mastodon.Config{
		Server:       "https://mstdn.jp",
		ClientID:     "client-id",
		ClientSecret: "client-secret",
	})
	err := c.Authenticate(context.Background(), "your-email", "your-password")
	if err != nil {
		log.Fatal(err)
	}
	timeline, err := c.GetTimelineHome(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	for i := len(timeline) - 1; i >= 0; i-- {
		fmt.Println(timeline[i])
	}
}
```

## Status of implementations

* [x] GET /api/v1/accounts/:id
* [x] GET /api/v1/accounts/verify_credentials
* [x] PATCH /api/v1/accounts/update_credentials
* [x] GET /api/v1/accounts/:id/followers
* [x] GET /api/v1/accounts/:id/following
* [x] GET /api/v1/accounts/:id/statuses
* [x] POST /api/v1/accounts/:id/follow
* [x] POST /api/v1/accounts/:id/unfollow
* [x] GET /api/v1/accounts/:id/block
* [x] GET /api/v1/accounts/:id/unblock
* [x] GET /api/v1/accounts/:id/mute
* [x] GET /api/v1/accounts/:id/unmute
* [x] GET /api/v1/accounts/:id/lists
* [x] GET /api/v1/accounts/relationships
* [x] GET /api/v1/accounts/search
* [x] GET /api/v1/apps/verify_credentials
* [x] GET /api/v1/bookmarks
* [x] POST /api/v1/apps
* [x] GET /api/v1/blocks
* [x] GET /api/v1/conversations
* [x] DELETE /api/v1/conversations/:id
* [x] POST /api/v1/conversations/:id/read
* [x] GET /api/v1/favourites
* [x] GET /api/v1/filters
* [x] POST /api/v1/filters
* [x] GET /api/v1/filters/:id
* [x] PUT /api/v1/filters/:id
* [x] DELETE /api/v1/filters/:id
* [x] GET /api/v1/follow_requests
* [x] POST /api/v1/follow_requests/:id/authorize
* [x] POST /api/v1/follow_requests/:id/reject
* [x] POST /api/v1/follows
* [x] GET /api/v1/instance
* [x] GET /api/v1/instance/activity
* [x] GET /api/v1/instance/peers
* [x] GET /api/v1/lists
* [x] GET /api/v1/lists/:id/accounts
* [x] GET /api/v1/lists/:id
* [x] POST /api/v1/lists
* [x] PUT /api/v1/lists/:id
* [x] DELETE /api/v1/lists/:id
* [x] POST /api/v1/lists/:id/accounts
* [x] DELETE /api/v1/lists/:id/accounts
* [x] POST /api/v1/media
* [x] GET /api/v1/mutes
* [x] GET /api/v1/notifications
* [x] GET /api/v1/notifications/:id
* [x] POST /api/v1/notifications/:id/dismiss
* [x] POST /api/v1/notifications/clear
* [x] POST /api/v1/push/subscription
* [x] GET /api/v1/push/subscription
* [x] PUT /api/v1/push/subscription
* [x] DELETE /api/v1/push/subscription
* [x] GET /api/v1/reports
* [x] POST /api/v1/reports
* [x] GET /api/v2/search
* [x] GET /api/v1/statuses/:id
* [x] GET /api/v1/statuses/:id/context
* [x] GET /api/v1/statuses/:id/card
* [x] GET /api/v1/statuses/:id/reblogged_by
* [x] GET /api/v1/statuses/:id/favourited_by
* [x] POST /api/v1/statuses
* [x] DELETE /api/v1/statuses/:id
* [x] POST /api/v1/statuses/:id/reblog
* [x] POST /api/v1/statuses/:id/unreblog
* [x] POST /api/v1/statuses/:id/favourite
* [x] POST /api/v1/statuses/:id/unfavourite
* [x] POST /api/v1/statuses/:id/bookmark
* [x] POST /api/v1/statuses/:id/unbookmark
* [x] GET /api/v1/timelines/home
* [x] GET /api/v1/timelines/public
* [x] GET /api/v1/timelines/tag/:hashtag
* [x] GET /api/v1/timelines/list/:id
* [x] GET /api/v1/streaming/user
* [x] GET /api/v1/streaming/public
* [x] GET /api/v1/streaming/hashtag?tag=:hashtag
* [x] GET /api/v1/streaming/hashtag/local?tag=:hashtag
* [x] GET /api/v1/streaming/list?list=:list_id
* [x] GET /api/v1/streaming/direct

## Installation

```shell
go install github.com/mattn/go-mastodon@latest
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)
