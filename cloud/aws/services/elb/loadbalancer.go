// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package elb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/instrumentation"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/wait"
	"strings"
	"time"
)

// ReconcileLoadbalancers reconciles the load balancers for the given cluster.
func (s *Service) ReconcileLoadbalancers(ctx context.Context, clusterName string, network *v1alpha1.Network) error {
	ctx, span := trace.StartSpan(
		ctx, instrumentation.MethodName("services", "ec2", "ReconcileLoadbalancers"),
	)
	defer span.End()
	glog.V(2).Info("Reconciling load balancers")

	// Get default api server spec.
	spec, err := s.getAPIServerClassicELBSpec(clusterName, network)

	if err != nil {
		return err
	}

	// Describe or create.
	apiELB, err := s.describeClassicELB(ctx, spec.Name)
	if IsNotFound(err) {
		apiELB, err = s.createClassicELB(ctx, spec)
		if err != nil {
			return err
		}

		glog.V(2).Infof("Created new classic load balancer for apiserver: %v", apiELB)
	} else if err != nil {
		return err
	}

	// TODO(vincepri): check if anything has changed and reconcile as necessary.

	apiELB.DeepCopyInto(&network.APIServerELB)

	glog.V(2).Info("Reconcile load balancers completed successfully")
	return nil
}

// DeleteLoadbalancers deletes the load balancers for the given cluster.
func (s *Service) DeleteLoadbalancers(ctx context.Context, clusterName string, network *v1alpha1.Network) error {
	ctx, span := trace.StartSpan(
		ctx, instrumentation.MethodName("services", "ec2", "DeleteLoadbalancers"),
	)
	defer span.End()
	glog.V(2).Info("Delete load balancers")

	// Get default api server spec.
	spec, err := s.getAPIServerClassicELBSpec(clusterName, network)

	if err != nil {
		// The ELB never existed
		return nil
	}

	if spec.Name == "" {
		// The ELB never existed
		return nil
	}

	// Describe or create.
	apiELB, err := s.describeClassicELB(ctx, spec.Name)
	if IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}

	if err := s.deleteClassicELBAndWait(ctx, apiELB.Name); err != nil {
		return err
	}

	glog.V(2).Info("Deleting load balancers completed successfully")
	return nil
}

// RegisterInstanceWithClassicELB registers an instance with a classic ELB
func (s *Service) RegisterInstanceWithClassicELB(ctx context.Context, instanceID string, loadBalancer string) error {
	ctx, span := trace.StartSpan(
		ctx, instrumentation.MethodName("services", "ec2", "DeleteLoadbalancers"),
	)
	defer span.End()
	input := &elb.RegisterInstancesWithLoadBalancerInput{
		Instances:        []*elb.Instance{{InstanceId: aws.String(instanceID)}},
		LoadBalancerName: aws.String(loadBalancer),
	}

	check := func() (done bool, err error) {
		_, err = s.ELB.RegisterInstancesWithLoadBalancer(input)

		if code, _ := awserrors.Code(err); code == "LoadBalancerNotFound" {
			return false, nil
		}

		if err != nil {
			return false, err
		}

		return true, nil

	}

	if err := wait.WaitForWithRetryable(ctx, wait.NewBackoff(), check, []string{}); err != nil {
		return errors.Wrapf(err, "failed to wait register instance %q with classic ELB %q", instanceID, loadBalancer)
	}

	return nil
}

func (s *Service) getAPIServerClassicELBSpec(clusterName string, network *v1alpha1.Network) (*v1alpha1.ClassicELB, error) {
	if network.SecurityGroups[v1alpha1.SecurityGroupControlPlane] == nil {
		return nil, NewFailedDependency(errors.New("security groups not ready"))
	}

	res := &v1alpha1.ClassicELB{
		Name:   fmt.Sprintf("%s-apiserver", clusterName),
		Scheme: v1alpha1.ClassicELBSchemeInternetFacing,
		Listeners: []*v1alpha1.ClassicELBListener{
			{
				Protocol:         v1alpha1.ClassicELBProtocolTCP,
				Port:             6443,
				InstanceProtocol: v1alpha1.ClassicELBProtocolTCP,
				InstancePort:     6443,
			},
		},
		HealthCheck: &v1alpha1.ClassicELBHealthCheck{
			Target:             fmt.Sprintf("%v:%d", v1alpha1.ClassicELBProtocolTCP, 6443),
			Interval:           10 * time.Second,
			Timeout:            5 * time.Second,
			HealthyThreshold:   5,
			UnhealthyThreshold: 3,
		},
		SecurityGroupIDs: []string{network.SecurityGroups[v1alpha1.SecurityGroupControlPlane].ID},
		Tags:             s.buildTags(clusterName, ResourceLifecycleOwned, "", TagValueAPIServerRole, nil),
	}

	for _, sn := range network.Subnets.FilterPrivate() {
		res.SubnetIDs = append(res.SubnetIDs, sn.ID)
	}

	return res, nil
}

func (s *Service) createClassicELB(ctx context.Context, spec *v1alpha1.ClassicELB) (*v1alpha1.ClassicELB, error) {
	ctx, span := trace.StartSpan(
		ctx, instrumentation.MethodName("services", "elb", "createClassicELB"),
	)
	defer span.End()
	input := &elb.CreateLoadBalancerInput{
		LoadBalancerName: aws.String(spec.Name),
		Subnets:          aws.StringSlice(spec.SubnetIDs),
		SecurityGroups:   aws.StringSlice(spec.SecurityGroupIDs),
		Scheme:           aws.String(string(spec.Scheme)),
		Tags:             mapToTags(spec.Tags),
	}

	for _, ln := range spec.Listeners {
		input.Listeners = append(input.Listeners, &elb.Listener{
			Protocol:         aws.String(string(ln.Protocol)),
			LoadBalancerPort: aws.Int64(ln.Port),
			InstanceProtocol: aws.String(string(ln.InstanceProtocol)),
			InstancePort:     aws.Int64(ln.InstancePort),
		})
	}

	out, err := s.ELB.CreateLoadBalancer(input)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create classic load balancer: %v", spec)
	}

	if spec.HealthCheck != nil {
		hc := &elb.ConfigureHealthCheckInput{
			LoadBalancerName: aws.String(spec.Name),
			HealthCheck: &elb.HealthCheck{
				Target:             aws.String(spec.HealthCheck.Target),
				Interval:           aws.Int64(int64(spec.HealthCheck.Interval.Seconds())),
				Timeout:            aws.Int64(int64(spec.HealthCheck.Timeout.Seconds())),
				HealthyThreshold:   aws.Int64(spec.HealthCheck.HealthyThreshold),
				UnhealthyThreshold: aws.Int64(spec.HealthCheck.UnhealthyThreshold),
			},
		}

		if _, err := s.ELB.ConfigureHealthCheck(hc); err != nil {
			return nil, errors.Wrapf(err, "failed to configure health check for classic load balancer: %v", spec)
		}
	}

	glog.V(2).Infof("Created load balancer with dns name: %q", *out.DNSName)

	res := spec.DeepCopy()
	res.DNSName = *out.DNSName
	return res, nil
}

func (s *Service) deleteClassicELB(ctx context.Context, name string) error {
	ctx, span := trace.StartSpan(
		ctx, instrumentation.MethodName("services", "elb", "deleteClassicELB"),
	)
	defer span.End()
	input := &elb.DeleteLoadBalancerInput{
		LoadBalancerName: aws.String(name),
	}

	if _, err := s.ELB.DeleteLoadBalancer(input); err != nil {
		return err
	}
	return nil
}

func (s *Service) deleteClassicELBAndWait(ctx context.Context, name string) error {
	ctx, span := trace.StartSpan(
		ctx, instrumentation.MethodName("services", "elb", "deleteClassicELBAndWait"),
	)
	defer span.End()
	if err := s.deleteClassicELB(ctx, name); err != nil {
		return err
	}

	input := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: aws.StringSlice([]string{name}),
	}

	check := func() (done bool, err error) {
		out, err := s.ELB.DescribeLoadBalancers(input)

		// ELB already deleted.
		if len(out.LoadBalancerDescriptions) == 0 {
			return true, nil
		}

		if code, _ := awserrors.Code(err); code == "LoadBalancerNotFound" {
			return true, nil
		}

		if err != nil {
			return false, err
		}

		return false, nil

	}

	if err := wait.WaitForWithRetryable(ctx, wait.NewBackoff(), check, []string{}); err != nil {
		return errors.Wrapf(err, "failed to wait for ELB deletion %q", name)
	}

	return nil
}

func (s *Service) describeClassicELB(ctx context.Context, name string) (*v1alpha1.ClassicELB, error) {
	ctx, span := trace.StartSpan(
		ctx, instrumentation.MethodName("services", "elb", "describeClassicELB"),
	)
	defer span.End()
	input := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: aws.StringSlice([]string{name}),
	}

	out, err := s.ELB.DescribeLoadBalancers(input)
	if (err != nil && strings.Contains(err.Error(), "There is no ACTIVE Load Balancer")) || (out != nil && len(out.LoadBalancerDescriptions) == 0) {
		return nil, NewNotFound(errors.Errorf("no classic load balancer found with name %q", name))
	} else if err != nil {
		return nil, errors.Wrapf(err, "failed to describe classic load balancer: %s", name)
	}

	return fromSDKTypeToClassicELB(out.LoadBalancerDescriptions[0]), nil
}

func fromSDKTypeToClassicELB(v *elb.LoadBalancerDescription) *v1alpha1.ClassicELB {
	return &v1alpha1.ClassicELB{
		Name:             aws.StringValue(v.LoadBalancerName),
		Scheme:           v1alpha1.ClassicELBScheme(*v.Scheme),
		SubnetIDs:        aws.StringValueSlice(v.Subnets),
		SecurityGroupIDs: aws.StringValueSlice(v.SecurityGroups),
		DNSName:          aws.StringValue(v.DNSName),
	}
}

func fromSDKTypeToClassicListener(v *elb.Listener) *v1alpha1.ClassicELBListener {
	return &v1alpha1.ClassicELBListener{
		Protocol:         v1alpha1.ClassicELBProtocol(*v.Protocol),
		Port:             *v.LoadBalancerPort,
		InstanceProtocol: v1alpha1.ClassicELBProtocol(*v.InstanceProtocol),
		InstancePort:     *v.InstancePort,
	}
}

func fromSDKTypeToClassicHealthCheck(v *elb.HealthCheck) *v1alpha1.ClassicELBHealthCheck {
	return &v1alpha1.ClassicELBHealthCheck{
		Target:             *v.Target,
		Interval:           time.Duration(*v.Interval) * time.Second,
		Timeout:            time.Duration(*v.Timeout) * time.Second,
		HealthyThreshold:   *v.HealthyThreshold,
		UnhealthyThreshold: *v.UnhealthyThreshold,
	}
}
