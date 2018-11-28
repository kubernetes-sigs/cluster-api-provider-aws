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
	"fmt"
	"time"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/wait"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/pkg/errors"
	"k8s.io/klog"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
)

// ReconcileLoadbalancers reconciles the load balancers for the given cluster.
func (s *Service) ReconcileLoadbalancers() error {
	klog.V(2).Info("Reconciling load balancers")

	// Get default api server spec.
	spec := s.getAPIServerClassicELBSpec()

	// Describe or create.
	apiELB, err := s.describeClassicELB(spec.Name)
	if IsNotFound(err) {
		apiELB, err = s.createClassicELB(spec)
		if err != nil {
			return err
		}

		klog.V(2).Infof("Created new classic load balancer for apiserver: %v", apiELB)
	} else if err != nil {
		return err
	}

	// TODO(vincepri): check if anything has changed and reconcile as necessary.
	apiELB.DeepCopyInto(&s.scope.Network().APIServerELB)
	klog.V(2).Info("Reconcile load balancers completed successfully")
	return nil
}

// GetAPIServerDNSName returns the DNS name endpoint for the API server
func (s *Service) GetAPIServerDNSName() (string, error) {
	apiELB, err := s.describeClassicELB(GenerateELBName(s.scope.Name(), TagValueAPIServerRole))

	if err != nil {
		return "", err
	}

	return apiELB.DNSName, nil
}

// DeleteLoadbalancers deletes the load balancers for the given cluster.
func (s *Service) DeleteLoadbalancers() error {
	klog.V(2).Info("Deleting load balancers")

	// Get default api server spec.
	spec := s.getAPIServerClassicELBSpec()

	// Describe or create.
	apiELB, err := s.describeClassicELB(spec.Name)
	if IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}

	if err := s.deleteClassicELBAndWait(apiELB.Name); err != nil {
		return err
	}

	klog.V(2).Info("Deleting load balancers completed successfully")
	return nil
}

// RegisterInstanceWithClassicELB registers an instance with a classic ELB
func (s *Service) RegisterInstanceWithClassicELB(instanceID string, loadBalancer string) error {
	input := &elb.RegisterInstancesWithLoadBalancerInput{
		Instances:        []*elb.Instance{{InstanceId: aws.String(instanceID)}},
		LoadBalancerName: aws.String(loadBalancer),
	}

	_, err := s.scope.ELB.RegisterInstancesWithLoadBalancer(input)
	if err != nil {
		return err
	}

	return nil
}

// RegisterInstanceWithAPIServerELB registers an instance with a classic ELB
func (s *Service) RegisterInstanceWithAPIServerELB(instanceID string) error {
	input := &elb.RegisterInstancesWithLoadBalancerInput{
		Instances:        []*elb.Instance{{InstanceId: aws.String(instanceID)}},
		LoadBalancerName: aws.String(GenerateELBName(s.scope.Name(), TagValueAPIServerRole)),
	}

	_, err := s.scope.ELB.RegisterInstancesWithLoadBalancer(input)
	if err != nil {
		return err
	}

	return nil
}

// GenerateELBName generates a formatted ELB name
func GenerateELBName(clusterName string, elbName string) string {
	return fmt.Sprintf("%s-%s", clusterName, elbName)
}

func (s *Service) getAPIServerClassicELBSpec() *v1alpha1.ClassicELB {
	res := &v1alpha1.ClassicELB{
		Name:   GenerateELBName(s.scope.Name(), TagValueAPIServerRole),
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
		SecurityGroupIDs: []string{s.scope.SecurityGroups()[v1alpha1.SecurityGroupControlPlane].ID},
		Tags:             s.buildTags(s.scope.Name(), ResourceLifecycleOwned, "", TagValueAPIServerRole, nil),
	}

	for _, sn := range s.scope.Subnets().FilterPublic() {
		res.SubnetIDs = append(res.SubnetIDs, sn.ID)
	}

	return res
}

func (s *Service) createClassicELB(spec *v1alpha1.ClassicELB) (*v1alpha1.ClassicELB, error) {
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

	out, err := s.scope.ELB.CreateLoadBalancer(input)
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

		if _, err := s.scope.ELB.ConfigureHealthCheck(hc); err != nil {
			return nil, errors.Wrapf(err, "failed to configure health check for classic load balancer: %v", spec)
		}
	}

	klog.V(2).Infof("Created load balancer with dns name: %q", *out.DNSName)

	res := spec.DeepCopy()
	res.DNSName = *out.DNSName
	return res, nil
}

func (s *Service) deleteClassicELB(name string) error {
	input := &elb.DeleteLoadBalancerInput{
		LoadBalancerName: aws.String(name),
	}

	if _, err := s.scope.ELB.DeleteLoadBalancer(input); err != nil {
		return err
	}
	return nil
}

func (s *Service) deleteClassicELBAndWait(name string) error {
	if err := s.deleteClassicELB(name); err != nil {
		return err
	}

	input := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: aws.StringSlice([]string{name}),
	}

	checkForELBDeletion := func() (done bool, err error) {
		out, err := s.scope.ELB.DescribeLoadBalancers(input)

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

	if err := wait.WaitForWithRetryable(wait.NewBackoff(), checkForELBDeletion, []string{}); err != nil {
		return errors.Wrapf(err, "failed to wait for ELB deletion %q", name)
	}

	return nil
}

func (s *Service) describeClassicELB(name string) (*v1alpha1.ClassicELB, error) {
	input := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: aws.StringSlice([]string{name}),
	}

	out, err := s.scope.ELB.DescribeLoadBalancers(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case elb.ErrCodeAccessPointNotFoundException:
				return nil, errors.Wrapf(err, "no classic load balancer found with name: %q", name)
			case elb.ErrCodeDependencyThrottleException:
				return nil, errors.Wrap(err, "too many requests made to the ELB service")
			default:
				return nil, errors.Wrap(err, "unexpected aws error")
			}
		} else {
			return nil, errors.Wrapf(err, "failed to describe classic load balancer: %s", name)
		}
	}

	if out == nil && len(out.LoadBalancerDescriptions) == 0 {
		return nil, NewNotFound(fmt.Errorf("no classic load balancer found with name %q", name))
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
