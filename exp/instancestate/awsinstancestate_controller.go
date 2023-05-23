/*
Copyright 2020 The Kubernetes Authors.

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

package instancestate

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/controllers"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/instancestate"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
)

// Ec2InstanceStateLabelKey defines an ec2 instance state label.
const Ec2InstanceStateLabelKey = "ec2-instance-state"

// ASGInstanceStateLabelKey defines an ASG EC2 instance state label.
const ASGInstanceStateLabelKey = "asg-instance-state"

// EC2InstanceHealthStateLabelKey defines an EC2 instance health state label.
const EC2InstanceHealthStateLabelKey = "ec2-instance-health-state"

// AwsInstanceStateReconciler reconciles a AwsInstanceState object.
type AwsInstanceStateReconciler struct {
	client.Client
	Log               logr.Logger
	sqsServiceFactory func() sqsiface.SQSAPI
	queueURLs         sync.Map
	Endpoints         []scope.ServiceEndpoint
	WatchFilterValue  string
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsclusters,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsmachines,verbs=get;list;watch

func (r *AwsInstanceStateReconciler) getSQSService(region string) (sqsiface.SQSAPI, error) {
	if r.sqsServiceFactory != nil {
		return r.sqsServiceFactory(), nil
	}

	globalScope, err := scope.NewGlobalScope(scope.GlobalScopeParams{
		ControllerName: "awsinstancestate",
		Region:         region,
		Endpoints:      r.Endpoints,
	})

	if err != nil {
		return nil, err
	}
	return scope.NewGlobalSQSClient(globalScope, globalScope), nil
}

func (r *AwsInstanceStateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Fetch the AWSCluster instance
	awsCluster := &infrav1.AWSCluster{}
	err := r.Get(ctx, req.NamespacedName, awsCluster)
	if err != nil {
		if apierrors.IsNotFound(err) {
			r.Log.Info("cluster not found, removing queue URL", "cluster", klog.KRef(req.Namespace, req.Name))
			r.queueURLs.Delete(req.Name)
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Handle deleted clusters
	if !awsCluster.DeletionTimestamp.IsZero() {
		r.queueURLs.Delete(req.Name)
		return reconcile.Result{}, nil
	}

	// retrieve queue URL if it isn't already tracked
	if _, ok := r.queueURLs.Load(awsCluster.Name); !ok {
		URL, err := r.getQueueURL(awsCluster)
		if err != nil {
			if queueNotFoundError(err) {
				return reconcile.Result{}, nil
			}
			return reconcile.Result{}, err
		}
		r.queueURLs.Store(awsCluster.Name, queueParams{region: awsCluster.Spec.Region, URL: URL})
	}

	return ctrl.Result{}, nil
}

func (r *AwsInstanceStateReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	go func() {
		r.watchQueuesForInstanceEvents()
	}()
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.AWSCluster{}).
		WithOptions(options).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(logger.FromContext(ctx).GetLogger(), r.WatchFilterValue)).
		Complete(r)
}

func (r *AwsInstanceStateReconciler) watchQueuesForInstanceEvents() {
	ctx := context.TODO()
	awsClusterList := &infrav1.AWSClusterList{}
	if err := r.Client.List(ctx, awsClusterList); err == nil {
		for i, cluster := range awsClusterList.Items {
			if URL, err := r.getQueueURL(&awsClusterList.Items[i]); err == nil {
				r.queueURLs.Store(cluster.Name, queueParams{region: cluster.Spec.Region, URL: URL})
			}
		}
	}
	for range time.Tick(1 * time.Second) {
		// go through each cluster and check for messages on its queue
		r.queueURLs.Range(func(key, val interface{}) bool {
			qp := val.(queueParams)
			go r.processQueue(ctx, qp)
			return true
		})
	}
}

func (r *AwsInstanceStateReconciler) processQueue(ctx context.Context, qp queueParams) {
	sqsSvs, err := r.getSQSService(qp.region)
	if err != nil {
		r.Log.Error(err, "unable to create SQS client")
		return
	}
	resp, err := receiveMessageFromQueue(sqsSvs, qp)
	if err != nil {
		r.Log.Error(err, "unable to receive message from queue")
		return
	}

	queueEvents := make([]instancestate.EventDetails, 10)
	for _, msg := range resp.Messages {
		event := instancestate.EventBridgeEvent{}
		err := json.Unmarshal([]byte(*msg.Body), &event)

		if err != nil {
			r.Log.Error(err, "unable to marshall")
			return
		}

		err = r.processMessage(ctx, event)
		if err != nil {
			r.Log.Error(err, "unable to process the message from queue")
			return
		}
		message := msg
		queueEvents = append(queueEvents, instancestate.EventDetails{
			Message:          message,
			EventBridgeEvent: event,
		})
	}
	instancestate.QueueEventMapping.Store(qp.URL, queueEvents)
}

func receiveMessageFromQueue(sqsSvs sqsiface.SQSAPI, qp queueParams) (*sqs.ReceiveMessageOutput, error) {
	resp, err := sqsSvs.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(qp.URL),
		MaxNumberOfMessages: aws.Int64(10),
		VisibilityTimeout:   aws.Int64(20),
		WaitTimeSeconds:     aws.Int64(10),
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// processMessage triggers reconcile on an AWSMachine if its EC2 instance state changed.
func (r *AwsInstanceStateReconciler) processMessage(ctx context.Context, event instancestate.EventBridgeEvent) error {
	instanceID := event.Detail.InstanceID
	if event.Source == "aws.health" {
		instanceID = event.Detail.AffectedEntities[0].EntityValue
	}
	// Fetch the awsMachine instance by InstanceID
	awsMachines := &infrav1.AWSMachineList{}
	err := r.List(ctx, awsMachines, client.MatchingFields{controllers.InstanceIDIndex: instanceID})

	if err != nil {
		r.Log.Error(err, "unable to list machines by instance ID", "instanceID", instanceID)
	}

	if len(awsMachines.Items) > 0 {
		machine := awsMachines.Items[0]
		if !machine.ObjectMeta.DeletionTimestamp.IsZero() {
			return nil
		}
		patchHelper, err := patch.NewHelper(&machine, r.Client)
		if err != nil {
			r.Log.Error(err, "unable to create patch helper")
		}

		labels, err := r.getMachineLabelOnEventType(machine, event)
		if err != nil {
			return err
		}

		// Trigger an update on the machine
		machine.SetLabels(labels)
		err = patchHelper.Patch(ctx, &machine)
		if err != nil {
			r.Log.Error(err, "unable to patch AWS machine")
		}
	}
	return nil
}

// getQueueURL retrieves the SQS queue URL for a given cluster.
func (r *AwsInstanceStateReconciler) getQueueURL(cluster *infrav1.AWSCluster) (string, error) {
	sqsSvs, err := r.getSQSService(cluster.Spec.Region)
	if err != nil {
		return "", err
	}
	queueName := instancestate.GenerateQueueName(cluster.Name)
	resp, err := sqsSvs.GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: aws.String(queueName)})

	if err != nil {
		return "", err
	}

	return *resp.QueueUrl, nil
}

func (r *AwsInstanceStateReconciler) getMachineLabelOnEventType(machine infrav1.AWSMachine, event instancestate.EventBridgeEvent) (map[string]string, error) {
	labels := machine.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}

	switch event.Source {
	case "aws.autoscaling":
		labels[ASGInstanceStateLabelKey] = string(event.Detail.State)
	case "aws.ec2":
		labels[Ec2InstanceStateLabelKey] = string(event.Detail.State)
	case "aws.health":
		if event.Detail.Service != "EC2" {
			return nil, fmt.Errorf("events from Amazon EventBridge for service (%s) are not supported", event.Detail.Service)
		}
		if event.Detail.EventTypeCategory != "scheduledChange" {
			return nil, fmt.Errorf("events from Amazon EventBridge with EventTypeCategory (%s) are not supported", event.Detail.EventTypeCategory)
		}
		labels[EC2InstanceHealthStateLabelKey] = event.Detail.EventTypeCategory
	default:
		r.Log.V(4).Info("event type from Amazon EventBridge is not supported", "event", event.Source)
	}
	return labels, nil
}

func queueNotFoundError(err error) bool {
	if aerr, ok := err.(awserr.Error); ok {
		if aerr.Code() == sqs.ErrCodeQueueDoesNotExist {
			return true
		}
	}
	return false
}

type queueParams struct {
	region string
	URL    string
}
