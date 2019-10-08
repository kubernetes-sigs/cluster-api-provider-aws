/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package elb

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elb"
	rgapi "github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/pkg/errors"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/wait"
)

// ResourceGroups are filtered by ARN identifier: https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html#arns-syntax
// this is the identifier for classic ELBs: https://docs.aws.amazon.com/IAM/latest/UserGuide/list_elasticloadbalancing.html#elasticloadbalancing-resources-for-iam-policies
const elbResourceType = "elasticloadbalancing:loadbalancer"

// ReconcileLoadbalancers reconciles the load balancers for the given cluster.
func (s *Service) ReconcileLoadbalancers() error {
	s.scope.V(2).Info("Reconciling load balancers")

	// Get default api server spec.
	spec := s.getAPIServerClassicELBSpec()

	// Describe or create.
	apiELB, err := s.describeClassicELB(spec.Name)
	if IsNotFound(err) {
		apiELB, err = s.createClassicELB(spec)
		if err != nil {
			return err
		}

		s.scope.V(2).Info("Created new classic load balancer for apiserver", "api-server-elb-name", apiELB.Name)
	} else if err != nil {
		return err
	}

	if !reflect.DeepEqual(spec.Attributes, apiELB.Attributes) {
		err := s.configureAttributes(apiELB.Name, spec.Attributes)
		if err != nil {
			return err
		}
	}

	// Reconciliate the subnets from the spec and the ones currently attached to the load balancer.
	if len(apiELB.SubnetIDs) != len(spec.SubnetIDs) {
		_, err := s.scope.ELB.AttachLoadBalancerToSubnets(&elb.AttachLoadBalancerToSubnetsInput{
			Subnets: aws.StringSlice(spec.SubnetIDs),
		})
		if err != nil {
			return errors.Wrapf(err, "failed to attach apiserver load balancer %q to subnets", apiELB.Name)
		}
	}

	// TODO(vincepri): check if anything has changed and reconcile as necessary.
	apiELB.DeepCopyInto(&s.scope.Network().APIServerELB)
	s.scope.V(4).Info("Control plane load balancer", "api-server-elb", apiELB)

	s.scope.V(2).Info("Reconcile load balancers completed successfully")
	return nil
}

// GetAPIServerDNSName returns the DNS name endpoint for the API server
func (s *Service) GetAPIServerDNSName() (string, error) {
	apiELB, err := s.describeClassicELB(GenerateELBName(s.scope.Name(), infrav1.APIServerRoleTagValue))
	if err != nil {
		return "", err
	}
	return apiELB.DNSName, nil
}

// DeleteLoadbalancers deletes the load balancers for the given cluster.
func (s *Service) DeleteLoadbalancers() error {
	s.scope.V(2).Info("Deleting load balancers")

	elbs, err := s.listOwnedELBs()
	if err != nil {
		return err
	}

	for _, elb := range elbs {
		s.scope.V(3).Info("deleting load balancer", "arn", elb)
		s.deleteClassicELB(elb)
	}

	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (done bool, err error) {
		elbs, err := s.listOwnedELBs()
		if err != nil {
			return false, err
		}

		return len(elbs) == 0, nil
	}); err != nil {
		return errors.Wrapf(err, "failed to wait for %q ELB deletions", s.scope.Name())
	}

	s.scope.V(2).Info("Deleting load balancers completed successfully")
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
func (s *Service) RegisterInstanceWithAPIServerELB(i *infrav1.Instance) error {
	name := GenerateELBName(s.scope.Name(), infrav1.APIServerRoleTagValue)
	out, err := s.describeClassicELB(name)
	if err != nil {
		return err
	}

	// Validate that the subnets associated with the load balancer has the instance AZ.
	instanceAZ := s.scope.Subnets().FindByID(i.SubnetID).AvailabilityZone
	found := false
	for _, subnetID := range out.SubnetIDs {
		if subnet := s.scope.Subnets().FindByID(subnetID); subnet != nil && instanceAZ == subnet.AvailabilityZone {
			found = true
			break
		}
	}
	if !found {
		return errors.Errorf("failed to register instance with APIServer ELB %q: instance is in availability zone %q, no public subnets attached to the ELB in the same zone", name, instanceAZ)
	}

	input := &elb.RegisterInstancesWithLoadBalancerInput{
		Instances:        []*elb.Instance{{InstanceId: aws.String(i.ID)}},
		LoadBalancerName: aws.String(name),
	}

	_, err = s.scope.ELB.RegisterInstancesWithLoadBalancer(input)
	if err != nil {
		return err
	}

	return nil
}

// GenerateELBName generates a formatted ELB name
func GenerateELBName(clusterName string, elbName string) string {
	return fmt.Sprintf("%s-%s", clusterName, elbName)
}

func (s *Service) getAPIServerClassicELBSpec() *infrav1.ClassicELB {
	res := &infrav1.ClassicELB{
		Name:   GenerateELBName(s.scope.Name(), infrav1.APIServerRoleTagValue),
		Scheme: s.scope.ControlPlaneLoadBalancerScheme(),
		Listeners: []*infrav1.ClassicELBListener{
			{
				Protocol:         infrav1.ClassicELBProtocolTCP,
				Port:             s.scope.APIServerPort(),
				InstanceProtocol: infrav1.ClassicELBProtocolTCP,
				InstancePort:     6443,
			},
		},
		HealthCheck: &infrav1.ClassicELBHealthCheck{
			Target:             fmt.Sprintf("%v:%d", infrav1.ClassicELBProtocolSSL, 6443),
			Interval:           10 * time.Second,
			Timeout:            5 * time.Second,
			HealthyThreshold:   5,
			UnhealthyThreshold: 3,
		},
		SecurityGroupIDs: []string{s.scope.SecurityGroups()[infrav1.SecurityGroupControlPlane].ID},
		Attributes: infrav1.ClassicELBAttributes{
			IdleTimeout: 10 * time.Minute,
		},
	}

	res.Tags = infrav1.Build(infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Role:        aws.String(infrav1.APIServerRoleTagValue),
		Additional:  s.scope.AdditionalTags(),
	})

	// The load balancer APIs require us to only attach one subnet for each AZ.
	zones := map[string]struct{}{}
	for _, sn := range s.scope.Subnets().FilterPublic() {
		if _, ok := zones[sn.AvailabilityZone]; !ok {
			zones[sn.AvailabilityZone] = struct{}{}
			res.SubnetIDs = append(res.SubnetIDs, sn.ID)
		}
	}

	return res
}

func (s *Service) createClassicELB(spec *infrav1.ClassicELB) (*infrav1.ClassicELB, error) {
	input := &elb.CreateLoadBalancerInput{
		LoadBalancerName: aws.String(spec.Name),
		Subnets:          aws.StringSlice(spec.SubnetIDs),
		SecurityGroups:   aws.StringSlice(spec.SecurityGroupIDs),
		Scheme:           aws.String(string(spec.Scheme)),
		Tags:             converters.MapToELBTags(spec.Tags),
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
		if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
			if _, err := s.scope.ELB.ConfigureHealthCheck(&elb.ConfigureHealthCheckInput{
				LoadBalancerName: aws.String(spec.Name),
				HealthCheck: &elb.HealthCheck{
					Target:             aws.String(spec.HealthCheck.Target),
					Interval:           aws.Int64(int64(spec.HealthCheck.Interval.Seconds())),
					Timeout:            aws.Int64(int64(spec.HealthCheck.Timeout.Seconds())),
					HealthyThreshold:   aws.Int64(spec.HealthCheck.HealthyThreshold),
					UnhealthyThreshold: aws.Int64(spec.HealthCheck.UnhealthyThreshold),
				},
			}); err != nil {
				return false, err
			}
			return true, nil
		}, awserrors.LoadBalancerNotFound); err != nil {
			return nil, errors.Wrapf(err, "failed to configure health check for classic load balancer: %v", spec)
		}
	}

	s.scope.V(2).Info("Created classic load balancer", "dns-name", *out.DNSName)

	res := spec.DeepCopy()
	res.DNSName = *out.DNSName
	return res, nil
}

func (s *Service) configureAttributes(name string, attributes infrav1.ClassicELBAttributes) error {
	attrs := &elb.ModifyLoadBalancerAttributesInput{
		LoadBalancerName:       aws.String(name),
		LoadBalancerAttributes: &elb.LoadBalancerAttributes{},
	}

	if attributes.IdleTimeout > 0 {
		attrs.LoadBalancerAttributes.ConnectionSettings = &elb.ConnectionSettings{
			IdleTimeout: aws.Int64(int64(attributes.IdleTimeout.Seconds())),
		}
	}

	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		if _, err := s.scope.ELB.ModifyLoadBalancerAttributes(attrs); err != nil {
			return false, err
		}
		return true, nil
	}, awserrors.LoadBalancerNotFound); err != nil {
		return errors.Wrapf(err, "failed to configure attributes for classic load balancer: %v", name)
	}

	return nil
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

func (s *Service) listByTag(tag string) ([]string, error) {
	input := rgapi.GetResourcesInput{
		ResourceTypeFilters: aws.StringSlice([]string{elbResourceType}),
		TagFilters: []*rgapi.TagFilter{
			{
				Key:    aws.String(tag),
				Values: aws.StringSlice([]string{string(infrav1.ResourceLifecycleOwned)}),
			},
		},
	}

	names := []string{}

	err := s.scope.ResourceTagging.GetResourcesPages(&input, func(r *rgapi.GetResourcesOutput, last bool) bool {
		for _, tagmapping := range r.ResourceTagMappingList {
			if tagmapping.ResourceARN != nil {
				// We can't use arn.Parse because the "Resource" is loadbalancer/<name>
				parts := strings.Split(*tagmapping.ResourceARN, "/")
				name := parts[len(parts)-1]
				if name == "" {
					s.scope.Info("failed to parse ARN", "arn", *tagmapping.ResourceARN, "tag", tag)
					continue
				}
				names = append(names, name)
			}
		}
		return true
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed to list %s ELBs by tag group", s.scope.Name())
	}

	return names, nil
}

func (s *Service) listOwnedELBs() ([]string, error) {
	// k8s.io/cluster/<name>, created by k/k cloud provider
	serviceTag := infrav1.ClusterAWSCloudProviderTagKey(s.scope.Name())
	arns, err := s.listByTag(serviceTag)
	if err != nil {
		return nil, err
	}

	// sigs.k8s.io/cluster-api-provider-aws/cluster/<name>, created by CAPA
	capaTag := infrav1.ClusterTagKey(s.scope.Name())
	clusterArns, err := s.listByTag(capaTag)
	if err != nil {
		return nil, err
	}

	return append(arns, clusterArns...), nil
}

func (s *Service) describeClassicELB(name string) (*infrav1.ClassicELB, error) {
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

	if out != nil && len(out.LoadBalancerDescriptions) == 0 {
		return nil, NewNotFound(fmt.Errorf("no classic load balancer found with name %q", name))
	}

	if s.scope.VPC().ID != "" && s.scope.VPC().ID != *out.LoadBalancerDescriptions[0].VPCId {
		return nil, errors.Errorf(
			"ELB names must be unique within a region: %q ELB already exists in this region in VPC %q",
			name, *out.LoadBalancerDescriptions[0].VPCId)
	}

	outAtt, err := s.scope.ELB.DescribeLoadBalancerAttributes(&elb.DescribeLoadBalancerAttributesInput{
		LoadBalancerName: aws.String(name),
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe classic load balancer %q attributes", name)
	}

	return fromSDKTypeToClassicELB(out.LoadBalancerDescriptions[0], outAtt.LoadBalancerAttributes), nil
}

func fromSDKTypeToClassicELB(v *elb.LoadBalancerDescription, attrs *elb.LoadBalancerAttributes) *infrav1.ClassicELB {
	res := &infrav1.ClassicELB{
		Name:             aws.StringValue(v.LoadBalancerName),
		Scheme:           infrav1.ClassicELBScheme(*v.Scheme),
		SubnetIDs:        aws.StringValueSlice(v.Subnets),
		SecurityGroupIDs: aws.StringValueSlice(v.SecurityGroups),
		DNSName:          aws.StringValue(v.DNSName),
	}

	if attrs.ConnectionSettings != nil && attrs.ConnectionSettings.IdleTimeout != nil {
		res.Attributes.IdleTimeout = time.Duration(*attrs.ConnectionSettings.IdleTimeout) * time.Second
	}

	return res
}
