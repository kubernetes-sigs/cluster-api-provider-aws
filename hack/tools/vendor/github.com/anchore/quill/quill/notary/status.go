package notary

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/wagoodman/go-progress"
)

type StatusConfig struct {
	Timeout time.Duration
	Poll    time.Duration
	Wait    bool
	stage   *progress.Stage
}

func (c *StatusConfig) WithProgress(stage *progress.Stage) *StatusConfig {
	c.stage = stage
	return c
}

func PollStatus(ctx context.Context, sub *Submission, cfg StatusConfig) (SubmissionStatus, error) {
	var err error

	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	var status SubmissionStatus = PendingStatus

	var count int
	for !status.isCompleted() {
		select {
		case <-ctx.Done():
			return TimeoutStatus, errors.New("timeout waiting for notarize submission response")

		default:
			count++
			status, err = sub.Status(ctx)
			if err != nil {
				return "", err
			}

			if cfg.stage != nil {
				cfg.stage.Current = fmt.Sprintf("status %q, poll %d", strings.ToLower(string(status)), count)
			}

			if !status.isCompleted() {
				time.Sleep(cfg.Poll)
			}
		}
	}

	if !status.isSuccessful() {
		logs, err := sub.Logs(ctx)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("submission result is %+v:\n%+v", status, logs)
	}

	return status, nil
}
