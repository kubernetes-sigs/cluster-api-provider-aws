package internal

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"

	configv1 "github.com/openshift/api/config/v1"
	configclientv1 "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"

	"github.com/openshift/cluster-version-operator/lib"
	"github.com/openshift/cluster-version-operator/lib/resourcebuilder"
	"github.com/openshift/cluster-version-operator/lib/resourcemerge"
	"github.com/openshift/cluster-version-operator/pkg/payload"
)

var (
	osScheme = runtime.NewScheme()
	osCodecs = serializer.NewCodecFactory(osScheme)

	osMapper = resourcebuilder.NewResourceMapper()
)

func init() {
	if err := configv1.AddToScheme(osScheme); err != nil {
		panic(err)
	}

	osMapper.RegisterGVK(configv1.SchemeGroupVersion.WithKind("ClusterOperator"), newClusterOperatorBuilder)
	osMapper.AddToMap(resourcebuilder.Mapper)
}

// readClusterOperatorV1OrDie reads clusteroperator object from bytes. Panics on error.
func readClusterOperatorV1OrDie(objBytes []byte) *configv1.ClusterOperator {
	requiredObj, err := runtime.Decode(osCodecs.UniversalDecoder(configv1.SchemeGroupVersion), objBytes)
	if err != nil {
		panic(err)
	}
	return requiredObj.(*configv1.ClusterOperator)
}

type clusterOperatorBuilder struct {
	client   configclientv1.ConfigV1Interface
	raw      []byte
	modifier resourcebuilder.MetaV1ObjectModifierFunc
}

func newClusterOperatorBuilder(config *rest.Config, m lib.Manifest) resourcebuilder.Interface {
	return &clusterOperatorBuilder{
		client: configclientv1.NewForConfigOrDie(config),
		raw:    m.Raw,
	}
}

func (b *clusterOperatorBuilder) WithModifier(f resourcebuilder.MetaV1ObjectModifierFunc) resourcebuilder.Interface {
	b.modifier = f
	return b
}

func (b *clusterOperatorBuilder) Do() error {
	os := readClusterOperatorV1OrDie(b.raw)
	if b.modifier != nil {
		b.modifier(os)
	}

	return waitForOperatorStatusToBeDone(osPollInternal, osPollTimeout, b.client, os)
}

const (
	osPollInternal = 1 * time.Second
	osPollTimeout  = 1 * time.Minute
)

func waitForOperatorStatusToBeDone(interval, timeout time.Duration, client configclientv1.ClusterOperatorsGetter, expected *configv1.ClusterOperator) error {
	var lastErr error
	err := wait.Poll(interval, timeout, func() (bool, error) {
		actual, err := client.ClusterOperators().Get(expected.Name, metav1.GetOptions{})
		if err != nil {
			lastErr = &payload.UpdateError{
				Nested:  err,
				Reason:  "ClusterOperatorNotAvailable",
				Message: fmt.Sprintf("Cluster operator %s has not yet reported success", expected.Name),
				Name:    expected.Name,
			}
			return false, nil
		}

		// undone is map of operand to tuple of (expected version, actual version)
		// for incomplete operands.
		undone := map[string][]string{}
		for _, expOp := range expected.Status.Versions {
			undone[expOp.Name] = []string{expOp.Version}
			for _, actOp := range actual.Status.Versions {
				if actOp.Name == expOp.Name {
					undone[expOp.Name] = append(undone[expOp.Name], actOp.Version)
					if actOp.Version == expOp.Version {
						delete(undone, expOp.Name)
					}
					break
				}
			}
		}
		if len(undone) > 0 {
			var keys []string
			for k := range undone {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			var reports []string
			for _, op := range keys {
				// we do not need to report `operator` version.
				if op == "operator" {
					continue
				}
				ver := undone[op]
				if len(ver) == 1 {
					reports = append(reports, fmt.Sprintf("missing version information for %s", op))
					continue
				}
				reports = append(reports, fmt.Sprintf("upgrading %s from %s to %s", op, ver[1], ver[0]))
			}
			message := fmt.Sprintf("Cluster operator %s is still updating", actual.Name)
			if len(reports) > 0 {
				message = fmt.Sprintf("Cluster operator %s is still updating: %s", actual.Name, strings.Join(reports, ", "))
			}
			lastErr = &payload.UpdateError{
				Nested:  errors.New(lowerFirst(message)),
				Reason:  "ClusterOperatorNotAvailable",
				Message: message,
				Name:    actual.Name,
			}
			return false, nil
		}

		available := false
		progressing := true
		failing := true
		for _, condition := range actual.Status.Conditions {
			switch {
			case condition.Type == configv1.OperatorAvailable && condition.Status == configv1.ConditionTrue:
				available = true
			case condition.Type == configv1.OperatorProgressing && condition.Status == configv1.ConditionFalse:
				progressing = false
			case condition.Type == configv1.OperatorFailing && condition.Status == configv1.ConditionFalse:
				failing = false
			}
		}
		// if we're at the correct version, and available, not progressing, and not failing, we are done
		if available && !progressing && !failing {
			return true, nil
		}

		if c := resourcemerge.FindOperatorStatusCondition(actual.Status.Conditions, configv1.OperatorFailing); c != nil && c.Status == configv1.ConditionTrue {
			message := fmt.Sprintf("Cluster operator %s is reporting a failure", actual.Name)
			if len(c.Message) > 0 {
				message = fmt.Sprintf("Cluster operator %s is reporting a failure: %s", actual.Name, c.Message)
			}
			lastErr = &payload.UpdateError{
				Nested:  errors.New(lowerFirst(message)),
				Reason:  "ClusterOperatorFailing",
				Message: message,
				Name:    actual.Name,
			}
			return false, nil
		}

		lastErr = &payload.UpdateError{
			Nested: fmt.Errorf("cluster operator %s is not done; it is available=%v, progressing=%v, failing=%v",
				actual.Name, available, progressing, failing,
			),
			Reason:  "ClusterOperatorNotAvailable",
			Message: fmt.Sprintf("Cluster operator %s has not yet reported success", actual.Name),
			Name:    actual.Name,
		}
		return false, nil
	})
	if err != nil {
		if err == wait.ErrWaitTimeout && lastErr != nil {
			return lastErr
		}
		return err
	}
	return nil
}

func lowerFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}
