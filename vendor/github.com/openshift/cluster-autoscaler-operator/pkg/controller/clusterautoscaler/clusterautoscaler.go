package clusterautoscaler

import (
	"fmt"

	v1 "github.com/openshift/cluster-autoscaler-operator/pkg/apis/autoscaling/v1"
)

// AutoscalerArg represents a command line argument to the cluster-autoscaler
// that may be combined with a value or numerical range.
type AutoscalerArg string

// String returns the argument as a plain string.
func (a AutoscalerArg) String() string {
	return string(a)
}

// Value returns the argument with the given value set.
func (a AutoscalerArg) Value(v interface{}) string {
	return fmt.Sprintf("%s=%v", a.String(), v)
}

// Range returns the argument with the given numerical range set.
func (a AutoscalerArg) Range(min, max int) string {
	return fmt.Sprintf("%s=%d:%d", a.String(), min, max)
}

// TypeRange returns the argument with the given type and numerical range set.
func (a AutoscalerArg) TypeRange(t string, min, max int) string {
	return fmt.Sprintf("%s=%s:%d:%d", a.String(), t, min, max)
}

// These constants represent the cluster-autoscaler arguments used by the
// operator when processing ClusterAutoscaler resources.
const (
	LogToStderrArg                  AutoscalerArg = "--logtostderr"
	NamespaceArg                    AutoscalerArg = "--namespace"
	CloudProviderArg                AutoscalerArg = "--cloud-provider"
	MaxGracefulTerminationSecArg    AutoscalerArg = "--max-graceful-termination-sec"
	ExpendablePodsPriorityCutoffArg AutoscalerArg = "--expendable-pods-priority-cutoff"
	ScaleDownEnabledArg             AutoscalerArg = "--scale-down-enabled"
	ScaleDownDelayAfterAddArg       AutoscalerArg = "--scale-down-delay-after-add"
	ScaleDownDelayAfterDeleteArg    AutoscalerArg = "--scale-down-delay-after-delete"
	ScaleDownDelayAfterFailureArg   AutoscalerArg = "--scale-down-delay-after-failure"
	ScaleDownUnneededTimeArg        AutoscalerArg = "--scale-down-unneeded-time"
	MaxNodesTotalArg                AutoscalerArg = "--max-nodes-total"
	CoresTotalArg                   AutoscalerArg = "--cores-total"
	MemoryTotalArg                  AutoscalerArg = "--memory-total"
	GPUTotalArg                     AutoscalerArg = "--gpu-total"
	VerbosityArg                    AutoscalerArg = "--v"
)

// AutoscalerArgs returns a slice of strings representing command line arguments
// to the cluster-autoscaler corresponding to the values in the given
// ClusterAutoscaler resource.
func AutoscalerArgs(ca *v1.ClusterAutoscaler, cfg *Config) []string {
	s := &ca.Spec

	args := []string{
		LogToStderrArg.String(),
		VerbosityArg.Value(cfg.Verbosity),
		CloudProviderArg.Value(cfg.CloudProvider),
		NamespaceArg.Value(cfg.Namespace),
	}

	if ca.Spec.MaxPodGracePeriod != nil {
		v := MaxGracefulTerminationSecArg.Value(*s.MaxPodGracePeriod)
		args = append(args, v)
	}

	if ca.Spec.PodPriorityThreshold != nil {
		v := ExpendablePodsPriorityCutoffArg.Value(*s.PodPriorityThreshold)
		args = append(args, v)
	}

	if ca.Spec.ResourceLimits != nil {
		args = append(args, ResourceArgs(s.ResourceLimits)...)
	}

	if ca.Spec.ScaleDown != nil {
		args = append(args, ScaleDownArgs(s.ScaleDown)...)
	}

	return args
}

// ScaleDownArgs returns a slice of strings representing command line arguments
// to the cluster-autoscaler corresponding to the values in the given
// ScaleDownConfig object.
func ScaleDownArgs(sd *v1.ScaleDownConfig) []string {
	if !sd.Enabled {
		return []string{ScaleDownEnabledArg.Value(false)}
	}

	args := []string{
		ScaleDownEnabledArg.Value(true),
	}

	if sd.DelayAfterAdd != nil {
		args = append(args, ScaleDownDelayAfterAddArg.Value(*sd.DelayAfterAdd))
	}

	if sd.DelayAfterDelete != nil {
		args = append(args, ScaleDownDelayAfterDeleteArg.Value(*sd.DelayAfterDelete))
	}

	if sd.DelayAfterFailure != nil {
		args = append(args, ScaleDownDelayAfterFailureArg.Value(*sd.DelayAfterFailure))
	}

	if sd.UnneededTime != nil {
		args = append(args, ScaleDownUnneededTimeArg.Value(*sd.UnneededTime))
	}

	return args
}

// ResourceArgs returns a slice of strings representing command line arguments
// to the cluster-autoscaler corresponding to the values in the given
// ResourceLimits object.
func ResourceArgs(rl *v1.ResourceLimits) []string {
	args := []string{}

	if rl.MaxNodesTotal != nil {
		args = append(args, MaxNodesTotalArg.Value(*rl.MaxNodesTotal))
	}

	if rl.Cores != nil {
		min, max := int(rl.Cores.Min), int(rl.Cores.Max)
		args = append(args, CoresTotalArg.Range(min, max))
	}

	if rl.Memory != nil {
		min, max := int(rl.Memory.Min), int(rl.Memory.Max)
		args = append(args, MemoryTotalArg.Range(min, max))
	}

	for _, g := range rl.GPUS {
		min, max := int(g.Min), int(g.Max)
		args = append(args, GPUTotalArg.TypeRange(g.Type, min, max))
	}

	return args
}
