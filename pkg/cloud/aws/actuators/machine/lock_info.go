/*
Copyright 2019 The Kubernetes Authors.

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

package machine

import (
	"encoding/json"

	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
)

type LockInformation struct {
	MachineName      string `json:"machineName"`
	IdempotencyToken string `json:"idempotencyToken"`
}

const LockInformationConfigMapKey = "lock-information"

func ReadLockInfo(cm *corev1.ConfigMap) (*LockInformation, error) {
	if cm == nil || cm.Data == nil {
		return nil, errors.New("recieved nil config map")
	}

	liString := cm.Data[LockInformationConfigMapKey]

	var lockInfo LockInformation
	if err := json.Unmarshal([]byte(liString), &lockInfo); err != nil {
		return nil, errors.Wrapf(err, "failed to decode lock information %q", liString)
	}

	return &lockInfo, nil

}

func WriteLockInfo(cm *corev1.ConfigMap, li *LockInformation) error {
	bytes, err := json.Marshal(li)
	if err != nil {
		return errors.WithStack(err)
	}

	if cm.Data == nil {
		cm.Data = make(map[string]string)
	}

	cm.Data[LockInformationConfigMapKey] = string(bytes)

	return nil
}
