package cincinnati

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/blang/semver"
	"github.com/google/uuid"
)

const (
	// GraphMediaType is the media-type specified in the HTTP Accept header
	// of requests sent to the Cincinnati-v1 Graph API.
	GraphMediaType = "application/json"
)

// Client is a Cincinnati client which can be used to fetch update graphs from
// an upstream Cincinnati stack.
type Client struct {
	id uuid.UUID
}

// NewClient creates a new Cincinnati client with the given client identifier.
func NewClient(id uuid.UUID) Client {
	return Client{id: id}
}

// Update is a single node from the update graph.
type Update node

// GetUpdates fetches the next-applicable update payloads from the specified
// upstream Cincinnati stack given the current version and channel. The next-
// applicable updates are determined by downloading the update graph, finding
// the current version within that graph (typically the root node), and then
// finding all of the children. These children are the available updates for
// the current version and their payloads indicate from where the actual update
// image can be downloaded.
func (c Client) GetUpdates(upstream string, channel string, version semver.Version) ([]Update, error) {
	// Prepare parametrized cincinnati query.
	cincinnatiURL, err := url.Parse(upstream)
	if err != nil {
		return nil, fmt.Errorf("failed to parse upstream URL: %s", err)
	}
	queryParams := cincinnatiURL.Query()
	queryParams.Add("channel", channel)
	queryParams.Add("id", c.id.String())
	queryParams.Add("version", version.String())
	cincinnatiURL.RawQuery = queryParams.Encode()

	// Download the update graph.
	req, err := http.NewRequest("GET", cincinnatiURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", GraphMediaType)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	// Parse the graph.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var graph graph
	if err = json.Unmarshal(body, &graph); err != nil {
		return nil, err
	}

	// Find the current version within the graph.
	var currentIdx int
	found := false
	for i, node := range graph.Nodes {
		if version.EQ(node.Version) {
			currentIdx = i
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("unknown version %s", version)
	}

	// Find the children of the current version.
	var nextIdxs []int
	for _, edge := range graph.Edges {
		if edge.Origin == currentIdx {
			nextIdxs = append(nextIdxs, edge.Destination)
		}
	}

	var updates []Update
	for _, i := range nextIdxs {
		updates = append(updates, Update(graph.Nodes[i]))
	}

	return updates, nil
}

type graph struct {
	Nodes []node
	Edges []edge
}

type node struct {
	Version semver.Version `json:"version"`
	Image   string         `json:"payload"`
}

type edge struct {
	Origin      int
	Destination int
}

// UnmarshalJSON unmarshals an edge in the update graph. The edge's JSON
// representation is a two-element array of indices, but Go's representation is
// a struct with two elements so this custom unmarshal method is required.
func (e *edge) UnmarshalJSON(data []byte) error {
	var fields []int
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	if len(fields) != 2 {
		return fmt.Errorf("expected 2 fields, found %d", len(fields))
	}

	e.Origin = fields[0]
	e.Destination = fields[1]

	return nil
}
