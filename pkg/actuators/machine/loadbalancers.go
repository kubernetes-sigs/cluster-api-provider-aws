package machine

import (
	"fmt"
	"strings"

	errorutil "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog"

	"github.com/aws/aws-sdk-go/aws"
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
		return err
	}
	// Use a map for target groups to get unique target group entries across load balancers
	targetGroups := map[string]*elbv2.TargetGroup{}
	for _, loadBalancer := range lbsResponse.LoadBalancers {
		klog.V(4).Infof("Retrieving target groups for load balancer %q", *loadBalancer.LoadBalancerName)
		targetGroupsInput := &elbv2.DescribeTargetGroupsInput{
			LoadBalancerArn: loadBalancer.LoadBalancerArn,
		}
		targetGroupsOutput, err := client.ELBv2DescribeTargetGroups(targetGroupsInput)
		if err != nil {
			klog.Errorf("Failed to retrieve load balancer target groups for %q: %v", *loadBalancer.LoadBalancerName, err)
			return err
		}
		for _, targetGroup := range targetGroupsOutput.TargetGroups {
			targetGroups[*targetGroup.TargetGroupArn] = targetGroup
		}
	}
	if klog.V(4) {
		targetGroupArns := make([]string, 0, len(targetGroups))
		for arn := range targetGroups {
			targetGroupArns = append(targetGroupArns, fmt.Sprintf("%q", arn))
		}
		klog.Infof("Registering instance %q with target groups: %v", *instance.InstanceId, strings.Join(targetGroupArns, ","))
	}
	errs := []error{}
	for _, targetGroup := range targetGroups {
		var target *elbv2.TargetDescription
		switch *targetGroup.TargetType {
		case elbv2.TargetTypeEnumInstance:
			target = &elbv2.TargetDescription{
				Id: instance.InstanceId,
			}
		case elbv2.TargetTypeEnumIp:
			target = &elbv2.TargetDescription{
				Id: instance.PrivateIpAddress,
			}
		}
		registerTargetsInput := &elbv2.RegisterTargetsInput{
			TargetGroupArn: targetGroup.TargetGroupArn,
			Targets:        []*elbv2.TargetDescription{target},
		}
		_, err := client.ELBv2RegisterTargets(registerTargetsInput)
		if err != nil {
			klog.Errorf("Failed to register instance %q with target group %q: %v", *instance.InstanceId, *targetGroup.TargetGroupArn, err)
			errs = append(errs, fmt.Errorf("%s: %v", *targetGroup.TargetGroupArn, err))
		}
	}
	if len(errs) > 0 {
		return errorutil.NewAggregate(errs)
	}
	return nil
}
