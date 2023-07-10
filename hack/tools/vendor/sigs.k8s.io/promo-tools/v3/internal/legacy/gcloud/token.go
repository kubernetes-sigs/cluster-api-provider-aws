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

package gcloud

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/sirupsen/logrus"

	"sigs.k8s.io/release-utils/command"
)

// Token is the oauth2 access token used for API calls over HTTP.
type Token string

// GetServiceAccountToken calls gcloud to get an access token for the specified
// service account.
func GetServiceAccountToken(
	serviceAccount string,
	useServiceAccount bool,
) (Token, error) {
	logrus.Infof("Obtaining access token for %s", serviceAccount)
	args := []string{
		"auth",
		"print-access-token",
	}
	args = MaybeUseServiceAccount(serviceAccount, useServiceAccount, args)

	cmd := command.New("gcloud", args...)

	// We use RunSilentSuccessOutput() to ensure the access token is captured,
	// but not displayed in logs.
	std, err := cmd.RunSilentSuccessOutput()
	// Do not log the token (stdout) on error, because it could
	// be that the token was valid, but that Run() failed for
	// other reasons. NEVER print the token as part of an error message!
	if err != nil {
		logrus.Errorf("could not execute cmd %v", cmd)
		return "", err
	}

	stdout := std.Output()
	token := Token(strings.TrimSpace(stdout))
	return token, nil
}

// MaybeUseServiceAccount injects a '--account=...' argument to the command with
// the given service account.
func MaybeUseServiceAccount(
	serviceAccount string,
	useServiceAccount bool,
	cmd []string,
) []string {
	if useServiceAccount && len(serviceAccount) > 0 {
		cmd = append(cmd, "")
		copy(cmd[2:], cmd[1:])
		cmd[1] = fmt.Sprintf("--account=%v", serviceAccount)
	}
	return cmd
}

// ActivateServiceAccounts uses the given CSV of JSON key filepaths to activate
// the associated service accounts.
func ActivateServiceAccounts(keyFilePaths string) error {
	r := csv.NewReader(strings.NewReader(keyFilePaths))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logrus.Fatal(err)
		}

		for _, keyFilePath := range record {
			if err := ActivateServiceAccount(keyFilePath); err != nil {
				logrus.Fatal(err)
			}
		}
	}
	return nil
}

// ActivateServiceAccount activates the service account with gcloud.
func ActivateServiceAccount(keyFilePath string) error {
	cmd := command.New(
		"gcloud",
		"auth",
		"activate-service-account",
		"--key-file="+keyFilePath,
	)

	return cmd.RunSuccess()
}
