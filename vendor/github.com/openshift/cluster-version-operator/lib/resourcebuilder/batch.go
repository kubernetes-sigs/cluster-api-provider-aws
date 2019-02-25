package resourcebuilder

import (
	"fmt"
	"time"

	"github.com/golang/glog"

	"github.com/openshift/cluster-version-operator/lib"
	"github.com/openshift/cluster-version-operator/lib/resourceapply"
	"github.com/openshift/cluster-version-operator/lib/resourceread"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	batchclientv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	"k8s.io/client-go/rest"
)

type jobBuilder struct {
	client   *batchclientv1.BatchV1Client
	raw      []byte
	modifier MetaV1ObjectModifierFunc
}

func newJobBuilder(config *rest.Config, m lib.Manifest) Interface {
	return &jobBuilder{
		client: batchclientv1.NewForConfigOrDie(withProtobuf(config)),
		raw:    m.Raw,
	}
}

func (b *jobBuilder) WithModifier(f MetaV1ObjectModifierFunc) Interface {
	b.modifier = f
	return b
}

func (b *jobBuilder) Do() error {
	job := resourceread.ReadJobV1OrDie(b.raw)
	if b.modifier != nil {
		b.modifier(job)
	}
	_, updated, err := resourceapply.ApplyJob(b.client, job)
	if err != nil {
		return err
	}
	if updated {
		return WaitForJobCompletion(b.client, job)
	}
	return nil
}

const (
	jobPollInterval = 1 * time.Second
	jobPollTimeout  = 5 * time.Minute
)

// WaitForJobCompletion waits for job to complete.
func WaitForJobCompletion(client batchclientv1.JobsGetter, job *batchv1.Job) error {
	return wait.Poll(jobPollInterval, jobPollTimeout, func() (bool, error) {
		j, err := client.Jobs(job.Namespace).Get(job.Name, metav1.GetOptions{})
		if err != nil {
			glog.Errorf("error getting Job %s: %v", job.Name, err)
			return false, nil
		}

		if j.Status.Succeeded > 0 {
			return true, nil
		}

		// Since we have filled in "activeDeadlineSeconds",
		// the Job will 'Active == 0' iff it exceeds the deadline.
		// Failed jobs will be recreated in the next run.
		if j.Status.Active == 0 && j.Status.Failed > 0 {
			reason := "DeadlineExceeded"
			message := "Job was active longer than specified deadline"
			if len(j.Status.Conditions) > 0 {
				reason, message = j.Status.Conditions[0].Reason, j.Status.Conditions[0].Message
			}
			return false, fmt.Errorf("deadline exceeded, reason: %q, message: %q", reason, message)
		}
		return false, nil
	})
}
