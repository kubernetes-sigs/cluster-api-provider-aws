/*
Copyright 2022 The Kubernetes Authors.

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
	"reflect"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

// QueueEventMapping stores the mapping of queueURL with all the events and its related details
// necessary to act upon based on the event type.
var QueueEventMapping sync.Map

// EventDetails represents all the necessary details present in the SQS message.
type EventDetails struct {
	EventBridgeEvent EventBridgeEvent
	Message          *sqs.Message
}

// EventBridgeEvent is a structure to hold generic event details from Amazon EventBridge.
type EventBridgeEvent struct {
	Version    string         `json:"version"`
	ID         string         `json:"id"`
	DetailType string         `json:"detail-type"`
	Source     string         `json:"source"`
	Account    string         `json:"account"`
	Time       string         `json:"time"`
	Region     string         `json:"region"`
	Resources  []string       `json:"resources"`
	Detail     *messageDetail `json:"detail"`
}

// messageDetail holds information on the affected instance/entity.
type messageDetail struct {
	InstanceID        string                `json:"instance-id,omitempty"`
	State             infrav1.InstanceState `json:"state,omitempty"`
	EventTypeCategory string                `json:"eventTypeCategory,omitempty"`
	Service           string                `json:"service,omitempty"`
	AffectedEntities  []affectedEntity      `json:"affectedEntities,omitempty"`
}

// affectedEntity holds information about an entity that is affected by a Health event.
type affectedEntity struct {
	EntityValue string `json:"entityValue"`
}

// GetEventsFromSQS is a utility used across CAPA controllers in tandem with a watcher
// to get the events in SQS queue to be processed further by the controller.
func GetEventsFromSQS(queueURL string) []EventDetails {
	var eventBridgeValues reflect.Value
	if events, ok := QueueEventMapping.Load(queueURL); ok {
		eventBridgeValues = reflect.ValueOf(events)
	}

	return eventBridgeValues.Interface().([]EventDetails)
}

// DeleteProcessedEvent deletes the event from queue if it's already processed by the awsinstancestate controller.
func (s Service) DeleteProcessedEvent(msg sqs.Message, queueURL string) error {
	if _, ok := QueueEventMapping.Load(queueURL); ok {
		QueueEventMapping.Delete(queueURL)
		_, err := s.SQSClient.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      aws.String(queueURL),
			ReceiptHandle: msg.ReceiptHandle,
		})

		if err != nil {
			return errors.Wrapf(err, "error deleting message, queueURL %v, messageReceiptHandle %v", queueURL, msg.ReceiptHandle)
		}
	}
	return nil
}
