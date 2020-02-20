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
	"k8s.io/apimachinery/pkg/util/sets"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/wait"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/internal/hash"
)

// ResourceGroups are filtered by ARN identifier: https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html#arns-syntax
// this is the identifier for classic ELBs: https://docs.aws.amazon.com/IAM/latest/UserGuide/list_elasticloadbalancing.html#elasticloadbalancing-resources-for-iam-policies
const elbResourceType = "elasticloadbalancing:loadbalancer"

// ReconcileLoadbalancers reconciles the load balancers for the given cluster.
func (s *Service) ReconcileLoadbalancers() error {
	s.scope.V(2).Info("Reconciling load balancers")

	// Get default api server spec.
	spec, err := s.getAPIServerClassicELBSpec()

	if err != nil {
		return err
	}

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

	if err := s.reconcileELBTags(apiELB.Name, spec.Tags); err != nil {
		return errors.Wrapf(err, "failed to reconcile tags for apiserver load balancer %q", apiELB.Name)
	}

	// Reconcile the subnets and availability zones from the spec
	// and the ones currently attached to the load balancer.
	if len(apiELB.SubnetIDs) != len(spec.SubnetIDs) {
		_, err := s.scope.ELB.AttachLoadBalancerToSubnets(&elb.AttachLoadBalancerToSubnetsInput{
			LoadBalancerName: &apiELB.Name,
			Subnets:          aws.StringSlice(spec.SubnetIDs),
		})
		if err != nil {
			return errors.Wrapf(err, "failed to attach apiserver load balancer %q to subnets", apiELB.Name)
		}
	}
	if len(apiELB.AvailabilityZones) != len(spec.AvailabilityZones) {
		apiELB.AvailabilityZones = spec.AvailabilityZones
	}

	// Reconcile the security groups from the spec and the ones currently attached to the load balancer
	if !sets.NewString(apiELB.SecurityGroupIDs...).Equal(sets.NewString(spec.SecurityGroupIDs...)) {
		_, err := s.scope.ELB.ApplySecurityGroupsToLoadBalancer(&elb.ApplySecurityGroupsToLoadBalancerInput{
			LoadBalancerName: &apiELB.Name,
			SecurityGroups:   aws.StringSlice(spec.SecurityGroupIDs),
		})
		if err != nil {
			return errors.Wrapf(err, "failed to apply security groups to load balancer %q", apiELB.Name)
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
	elbName, err := GenerateELBName(s.scope.Name())
	if err != nil {
		return "", err
	}
	apiELB, err := s.describeClassicELB(elbName)
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
		if err := s.deleteClassicELB(elb); err != nil {
			return err
		}
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
	name, err := GenerateELBName(s.scope.Name())
	if err != nil {
		return err
	}
	out, err := s.describeClassicELB(name)
	if err != nil {
		return err
	}

	// Validate that the subnets associated with the load balancer has the instance AZ.
	subnet := s.scope.Subnets().FindByID(i.SubnetID)
	if subnet == nil {
		return errors.Errorf("failed to attach load balancer subnets, could not find subnet %q description in AWSCluster", i.SubnetID)
	}
	instanceAZ := subnet.AvailabilityZone
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

// GenerateELBName generates a formatted ELB name via either
// concatenating the cluster name to the "-apiserver" suffix
// or computing a hash for clusters with names above 32 characters.
func GenerateELBName(clusterName string) (string, error) {
	standardELBName := generateStandardELBName(clusterName)
	if len(standardELBName) <= 32 {
		return standardELBName, nil
	}

	elbName, err := generateHashedELBName(clusterName)
	if err != nil {
		return "", err
	}

	return elbName, nil
}

// generateStandardELBName generates a formatted ELB name based on cluster
// and ELB name
func generateStandardELBName(clusterName string) string {
	elbCompatibleClusterName := strings.Replace(clusterName, ".", "-", -1)
	return fmt.Sprintf("%s-%s", elbCompatibleClusterName, infrav1.APIServerRoleTagValue)
}

// generateHashedELBName generates a 32-character hashed name based on cluster
// and ELB name
func generateHashedELBName(clusterName string) (string, error) {
	// hashSize = 32 - length of "k8s" - length of "-" = 28
	shortName, err := hash.Base36TruncatedHash(clusterName, 28)
	if err != nil {
		return "", errors.Wrap(err, "unable to create ELB name")
	}

	return fmt.Sprintf("%s-%s", shortName, "k8s"), nil
}

func (s *Service) getAPIServerClassicELBSpec() (*infrav1.ClassicELB, error) {
	elbName, err := GenerateELBName(s.scope.Name())
	if err != nil {
		return nil, err
	}
	res := &infrav1.ClassicELB{
		Name:   elbName,
		Scheme: s.scope.ControlPlaneLoadBalancerScheme(),
		Listeners: []*infrav1.ClassicELBListener{
			{
				Protocol:         infrav1.ClassicELBProtocolTCP,
				Port:             int64(s.scope.APIServerPort()),
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
		SecurityGroupIDs: []string{s.scope.SecurityGroups()[infrav1.SecurityGroupAPIServerLB].ID},
		Attributes: infrav1.ClassicELBAttributes{
			IdleTimeout: 10 * time.Minute,
		},
	}

	if s.scope.AWSCluster.Spec.ControlPlaneLoadBalancer != nil {
		res.Attributes.CrossZoneLoadBalancing = s.scope.AWSCluster.Spec.ControlPlaneLoadBalancer.CrossZoneLoadBalancing
	}

	res.Tags = infrav1.Build(infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Role:        aws.String(infrav1.APIServerRoleTagValue),
		Additional:  s.scope.AdditionalTags(),
	})

	// The load balancer APIs require us to only attach one subnet for each AZ.
	subnets := s.scope.Subnets().FilterPrivate()

	if s.scope.ControlPlaneLoadBalancerScheme() == infrav1.ClassicELBSchemeInternetFacing {
		subnets = s.scope.Subnets().FilterPublic()
	}

	for _, sn := range subnets {
		for _, az := range res.AvailabilityZones {
			if sn.AvailabilityZone == az {
				// If we already attached another subnet in the same AZ, there is no need to
				// add this subnet to the list of the ELB's subnets.
				continue
			}
		}
		res.AvailabilityZones = append(res.AvailabilityZones, sn.AvailabilityZone)
		res.SubnetIDs = append(res.SubnetIDs, sn.ID)
	}

	return res, nil
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
		LoadBalancerName: aws.String(name),
		LoadBalancerAttributes: &elb.LoadBalancerAttributes{
			CrossZoneLoadBalancing: &elb.CrossZoneLoadBalancing{
				Enabled: aws.Bool(attributes.CrossZoneLoadBalancing),
			},
		},
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

func (s *Service) reconcileELBTags(name string, desiredTags map[string]string) error {
	tags, err := s.scope.ELB.DescribeTags(&elb.DescribeTagsInput{
		LoadBalancerNames: []*string{aws.String(name)},
	})
	if err != nil {
		return err
	}

	if len(tags.TagDescriptions) == 0 {
		return errors.Errorf("no tag information returned for load balancer %q", name)
	}

	currentTags := converters.ELBTagsToMap(tags.TagDescriptions[0].Tags)

	addTagsInput := &elb.AddTagsInput{
		LoadBalancerNames: []*string{aws.String(name)},
	}

	removeTagsInput := &elb.RemoveTagsInput{
		LoadBalancerNames: []*string{aws.String(name)},
	}

	for k, v := range desiredTags {
		if val, ok := currentTags[k]; !ok || val != v {
			s.scope.V(4).Info("adding tag to load balancer", "elb-name", name, "key", k, "value", v)
			addTagsInput.Tags = append(addTagsInput.Tags, &elb.Tag{Key: aws.String(k), Value: aws.String(v)})
		}
	}

	for k := range currentTags {
		if _, ok := desiredTags[k]; !ok {
			s.scope.V(4).Info("removing tag from load balancer", "elb-name", name, "key", k)
			removeTagsInput.Tags = append(removeTagsInput.Tags, &elb.TagKeyOnly{Key: aws.String(k)})
		}
	}

	if len(addTagsInput.Tags) > 0 {
		if _, err := s.scope.ELB.AddTags(addTagsInput); err != nil {
			return err
		}
	}

	if len(removeTagsInput.Tags) > 0 {
		if _, err := s.scope.ELB.RemoveTags(removeTagsInput); err != nil {
			return err
		}
	}

	return nil
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

	res.Attributes.CrossZoneLoadBalancing = aws.BoolValue(attrs.CrossZoneLoadBalancing.Enabled)

	return res
}
