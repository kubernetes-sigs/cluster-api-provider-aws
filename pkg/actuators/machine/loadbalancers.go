package machine

import (
	"fmt"

	errorutil "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog/v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"

	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/client"
)

func registerWithClassicLoadBalancers(client awsclient.Client, names []string, instance *ec2.Instance) error {
	klog.V(4).Infof("Updating classic load balancer registration for %q", *instance.InstanceId)
	elbInstance := &elb.Instance{InstanceId: instance.InstanceId}
	var errs []error
	for _, elbName := range names {
		req := &elb.RegisterInstancesWithLoadBalancerInput{
			Instances:        []*elb.Instance{elbInstance},
			LoadBalancerName: aws.String(elbName),
		}
		_, err := client.RegisterInstancesWithLoadBalancer(req)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %v", elbName, err))
		}
	}

	if len(errs) > 0 {
		return errorutil.NewAggregate(errs)
	}
	return nil
}

func registerWithNetworkLoadBalancers(client awsclient.Client, names []string, instance *ec2.Instance) error {
	klog.V(4).Infof("Updating network load balancer registration for %q", *instance.InstanceId)
	targetGroups, err := gatherLoadBalancerTargetGroups(client, names)
	if err != nil {
		return err
	}

	errs := []error{}
	for _, targetGroup := range targetGroups {

		var target *elbv2.TargetDescription
		switch *targetGroup.TargetType {
		case elbv2.TargetTypeEnumInstance:
			target = &elbv2.TargetDescription{
				Id: instance.InstanceId,
			}
			klog.V(4).Infof("Registering instance %q by instance ID to target group: %v", *instance.InstanceId, *targetGroup.TargetGroupArn)
		case elbv2.TargetTypeEnumIp:
			target = &elbv2.TargetDescription{
				Id: instance.PrivateIpAddress,
			}
			klog.V(4).Infof("Registering instance %q by IP to target group: %v", *instance.InstanceId, *targetGroup.TargetGroupArn)
		}

		registeredTargets, err := gatherLoadBalancerTargetGroupRegisteredTargets(client, targetGroup.TargetGroupArn)
		if err != nil {
			klog.Errorf("Failed to gather registered targets for target group %q: %v", *targetGroup.TargetGroupArn, err)
			errs = append(errs, fmt.Errorf("%s: %v", *targetGroup.TargetGroupArn, err))
		}
		if registeredTargets != nil {
			if _, ok := registeredTargets[*target.Id]; ok {
				klog.V(4).Infof("Skipping registration for instance %q to target group %q: Instance already registered", *instance.InstanceId, *targetGroup.TargetGroupArn)
				continue
			}
		}

		registerTargetsInput := &elbv2.RegisterTargetsInput{
			TargetGroupArn: targetGroup.TargetGroupArn,
			Targets:        []*elbv2.TargetDescription{target},
		}
		if _, err := client.ELBv2RegisterTargets(registerTargetsInput); err != nil {
			klog.Errorf("Failed to register instance %q with target group %q: %v", *instance.InstanceId, *targetGroup.TargetGroupArn, err)
			errs = append(errs, fmt.Errorf("%s: %v", *targetGroup.TargetGroupArn, err))
		}
	}
	if len(errs) > 0 {
		return errorutil.NewAggregate(errs)
	}
	return nil
}

// deregisterNetworkLoadBalancers serves manual instance removal from Network LoadBalancer TargetGroup list
// for the instances attached by IP. Unlike instance reference, IP attachment should be cleaned manually.
func deregisterNetworkLoadBalancers(client awsclient.Client, names []string, instance *ec2.Instance) error {
	if instance.PrivateIpAddress == nil {
		klog.V(4).Infof("Instance %q does not have private ip, skipping...", *instance.InstanceId)
		return nil
	}

	klog.V(4).Infof("Removing network load balancer registration for %q", *instance.InstanceId)
	targetGroupsOutput, err := gatherLoadBalancerTargetGroups(client, names)
	if err != nil {
		return err
	}

	filteredGroupsByIP := []*elbv2.TargetGroup{}
	for _, targetGroup := range targetGroupsOutput {
		if *targetGroup.TargetType == elbv2.TargetTypeEnumIp {
			filteredGroupsByIP = append(filteredGroupsByIP, targetGroup)
		}
	}

	errs := []error{}
	for _, targetGroup := range filteredGroupsByIP {
		klog.V(4).Infof("Unregistering instance %q registered by ip from target group: %v", *instance.InstanceId, *targetGroup.TargetGroupArn)

		deregisterTargetsInput := &elbv2.DeregisterTargetsInput{
			TargetGroupArn: targetGroup.TargetGroupArn,
			Targets: []*elbv2.TargetDescription{{
				Id: instance.PrivateIpAddress,
			}},
		}
		_, err := client.ELBv2DeregisterTargets(deregisterTargetsInput)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case elbv2.ErrCodeInvalidTargetException, elbv2.ErrCodeTargetGroupNotFoundException:
					// Ignoring error when LB target group was already removed
					continue
				}
			}
			klog.Errorf("Failed to unregister instance %q from target group %q: %v", *instance.InstanceId, *targetGroup.TargetGroupArn, err)
			errs = append(errs, fmt.Errorf("%s: %v", *targetGroup.TargetGroupArn, err))
		}
	}
	if len(errs) > 0 {
		return errorutil.NewAggregate(errs)
	}
	return nil
}

func gatherLoadBalancerTargetGroups(client awsclient.Client, names []string) ([]*elbv2.TargetGroup, error) {
	lbNames := make([]*string, len(names))
	for i, name := range names {
		lbNames[i] = aws.String(name)
	}
	lbsRequest := &elbv2.DescribeLoadBalancersInput{
		Names: lbNames,
	}
	lbsResponse, err := client.ELBv2DescribeLoadBalancers(lbsRequest)
	if err != nil {
		klog.Errorf("Failed to describe load balancers %v: %v", names, err)
		return nil, err
	}
	// Use a map for target groups to get unique target group entries across load balancers
	targetGroups := []*elbv2.TargetGroup{}
	for _, loadBalancer := range lbsResponse.LoadBalancers {
		klog.V(4).Infof("Retrieving target groups for load balancer %s", *loadBalancer.LoadBalancerName)
		targetGroupsInput := &elbv2.DescribeTargetGroupsInput{
			LoadBalancerArn: loadBalancer.LoadBalancerArn,
		}
		targetGroupsOutput, err := client.ELBv2DescribeTargetGroups(targetGroupsInput)
		if err != nil {
			klog.Errorf("Failed to retrieve load balancer target groups for %q: %v", *loadBalancer.LoadBalancerName, err)
			return nil, err
		}
		targetGroups = append(targetGroups, targetGroupsOutput.TargetGroups...)
	}

	return targetGroups, nil
}

// gatherLoadBalancerTargetGroupRegisteredTargets looks for all targets that are registered to a particular targetGroup.
// Within the AWS API, the only way to find the targets that are registered is to look at the target health for the group.
// The target health response contains all of the targets and importantly, their IDs which we need later to compare with
// the target ID we are wanting to register.
func gatherLoadBalancerTargetGroupRegisteredTargets(client awsclient.Client, targetGroupArn *string) (map[string]struct{}, error) {
	targetHealthRequest := &elbv2.DescribeTargetHealthInput{
		TargetGroupArn: targetGroupArn,
	}
	targetHealthResponse, err := client.ELBv2DescribeTargetHealth(targetHealthRequest)
	if err != nil {
		klog.Errorf("Failed to describe target health: %v", err)
		return nil, err
	}

	targetIDs := make(map[string]struct{})
	for _, targetHealth := range targetHealthResponse.TargetHealthDescriptions {
		targetIDs[*targetHealth.Target.Id] = struct{}{}
	}
	return targetIDs, nil
}
