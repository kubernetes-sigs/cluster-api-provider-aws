// +build e2e

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

package shared

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ssm"
	expect "github.com/google/goexpect"
)

type instance struct {
	name       string
	instanceID string
}

// allMachines gets all EC2 instances at once, to save on DescribeInstances
// calls
func allMachines(ctx context.Context, e2eCtx *E2EContext) ([]instance, error) {
	ec2Svc := ec2.New(e2eCtx.AWSSession)
	resp, err := ec2Svc.DescribeInstancesWithContext(ctx, &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, err
	}
	if len(resp.Reservations) == 0 || len(resp.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("no machines found")
	}
	instances := []instance{}
	for _, r := range resp.Reservations {
		for _, i := range r.Instances {
			tags := i.Tags
			name := ""
			for _, t := range tags {
				if aws.StringValue(t.Key) == "Name" {
					name = aws.StringValue(t.Value)
				}
			}
			if name == "" {
				continue
			}
			instances = append(instances,
				instance{
					name:       name,
					instanceID: aws.StringValue(i.InstanceId),
				},
			)
		}
	}
	return instances, nil
}

type command struct {
	title string
	cmd   string
}

// commandsForMachine opens a terminal connection using AWS SSM Session Manager
// and executes the given commands, outputting the results to a file for each.
func commandsForMachine(ctx context.Context, e2eCtx *E2EContext, f *os.File, instanceID string, commands []command) {
	ssmSvc := ssm.New(e2eCtx.BootstrapUserAWSSession)
	sess, err := ssmSvc.StartSessionWithContext(ctx, &ssm.StartSessionInput{
		Target: aws.String(instanceID),
	})
	if err != nil {
		fmt.Fprintf(f, "unable to start session: err=%s", err)
		return
	}
	defer func() {
		if _, err := ssmSvc.TerminateSession(&ssm.TerminateSessionInput{SessionId: sess.SessionId}); err != nil {
			fmt.Fprintf(f, "unable to terminate session: err=%s", err)
		}
	}()
	sessionToken, err := json.Marshal(sess)
	if err != nil {
		fmt.Fprintf(f, "unable to marshal session: err=%s", err)
		return
	}
	cmdLine := fmt.Sprintf("session-manager-plugin %s %s StartSession %s", string(sessionToken), *ssmSvc.Client.Config.Region, ssmSvc.Client.Endpoint)
	e, _, err := expect.Spawn(cmdLine, -1)
	if err != nil {
		fmt.Fprintf(f, "unable to spawn AWS SSM Session Manager plugin: %s", err)
		return
	}
	defer e.Close()
	shellStart := regexp.MustCompile(`\n\$`)
	if result, _, err := e.Expect(shellStart, 10*time.Second); err != nil {
		fmt.Fprintf(f, "did not find shell: err=%s, output=%s", err, result)
		return
	}
	for _, c := range commands {
		logFile := path.Join(filepath.Dir(f.Name()), c.title+".log")
		if err := e.Send("sudo " + c.cmd + "\n"); err != nil {
			fmt.Fprintf(f, "unable to send command: err=%s", err)
			return
		}
		result, _, err := e.Expect(shellStart, 20*time.Second)
		if err := ioutil.WriteFile(logFile, []byte(result), os.ModePerm); err != nil {
			fmt.Fprintf(f, "error writing log file: err=%s", err)
			return
		}
		if err != nil {
			fmt.Fprintf(f, "error awaiting command output: err=%s", err)
			return
		}
	}
}
