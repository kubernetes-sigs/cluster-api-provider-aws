package notary

import (
	"context"
	"fmt"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/anchore/quill/internal/log"
)

type SubmissionStatus string

const (
	AcceptedStatus = "Accepted"
	PendingStatus  = "Pending"
	InvalidStatus  = "Invalid"
	RejectedStatus = "Rejected"
	TimeoutStatus  = "Timeout"
)

func (s SubmissionStatus) isCompleted() bool {
	switch s {
	case AcceptedStatus, RejectedStatus, InvalidStatus, TimeoutStatus:
		return true
	default:
		return false
	}
}

func (s SubmissionStatus) isSuccessful() bool {
	return s == AcceptedStatus
}

type Submission struct {
	api    api
	binary *Payload
	name   string
	id     string
}

type SubmissionList struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	CreatedDate string `json:"createdDate"`
}

func NewSubmission(a api, bin *Payload) *Submission {
	return &Submission{
		name:   filepath.Base(bin.Path + "-" + bin.Digest + "-" + randomString(8)),
		binary: bin,
		api:    a,
	}
}

func ExistingSubmission(a api, id string) *Submission {
	return &Submission{
		id:  id,
		api: a,
	}
}

func (s Submission) ID() string {
	return s.id
}

func (s *Submission) Start(ctx context.Context) error {
	if s.id != "" {
		return fmt.Errorf("submission already started")
	}

	log.WithFields("name", s.name).Debug("starting submission")

	if s.binary == nil {
		return fmt.Errorf("unable to start Submission without a binary")
	}

	response, err := s.api.submissionRequest(
		ctx,
		submissionRequest{
			Sha256:         s.binary.Digest,
			SubmissionName: s.name,
		},
	)

	if err != nil {
		return err
	}

	s.id = response.Data.ID

	log.WithFields("id", s.id, "name", s.name).Trace("received submission id")

	return s.api.uploadBinary(ctx, *response, *s.binary)
}

func (s Submission) Status(ctx context.Context) (SubmissionStatus, error) {
	log.WithFields("id", s.id).Trace("checking submission status")

	response, err := s.api.submissionStatusRequest(ctx, s.id)
	if err != nil {
		return "", err
	}

	log.WithFields("status", fmt.Sprintf("%q", response.Data.Attributes.Status), "id", s.id).Debug("submission status")

	switch response.Data.Attributes.Status {
	case "In Progress":
		return PendingStatus, nil
	case "Accepted":
		return AcceptedStatus, nil
	case "Invalid":
		return InvalidStatus, nil
	case "Rejected":
		return RejectedStatus, nil
	default:
		return "", fmt.Errorf("unexpected status: %s", response.Data.Attributes.Status)
	}
}

func (s Submission) Logs(ctx context.Context) (string, error) {
	return s.api.submissionLogs(ctx, s.id)
}

func (s Submission) List(ctx context.Context) ([]SubmissionList, error) {
	resp, err := s.api.submissionList(ctx)
	if err != nil {
		return nil, err
	}

	var results []SubmissionList
	for _, item := range resp.Data {
		results = append(results, SubmissionList{
			ID:          item.ID,
			Name:        item.Attributes.Name,
			Status:      item.Attributes.Status,
			CreatedDate: item.Attributes.CreatedDate,
		})
	}
	return results, nil
}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b) //nolint:gosec
	return fmt.Sprintf("%x", b)[:length]
}
