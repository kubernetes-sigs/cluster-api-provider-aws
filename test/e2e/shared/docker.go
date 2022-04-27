//go:build e2e
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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type AuthsObj struct {
	Auths map[string]types.AuthConfig `json:"auths"`
}

func DockerTag(source, target string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	if err := cli.ImageTag(ctx, source, target); err != nil {
		return err
	}
	return nil
}

func DockerImageList() ([]types.ImageSummary, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	result, err := cli.ImageList(ctx, types.ImageListOptions{All: true})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func DockerLogin(server, username, password string) (bool, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return false, err
	}

	auth := types.AuthConfig{
		Username:      username,
		Password:      password,
		ServerAddress: server,
	}

	result, err := cli.RegistryLogin(ctx, auth)
	if err != nil {
		return false, err
	}
	fmt.Printf("%+v\n", result)
	return true, nil
}

func DockerPush(image, username, password string) (bool, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return false, err
	}

	authObj := types.AuthConfig{
		Username: username,
		Password: password,
	}

	data, err := json.Marshal(authObj)
	if err != nil {
		return false, err
	}
	authStr := base64.URLEncoding.EncodeToString(data)

	opt := types.ImagePushOptions{
		RegistryAuth: authStr,
	}

	push, err := cli.ImagePush(ctx, image, opt)
	if err != nil {
		return false, err
	}

	defer push.Close()
	if _, err := io.Copy(os.Stdout, push); err != nil {
		return false, err
	}
	return true, nil
}

func DockerPull(image string) (bool, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return false, err
	}

	opt := types.ImagePullOptions{}

	pull, err := cli.ImagePull(ctx, image, opt)
	if err != nil {
		return false, err
	}
	defer pull.Close()
	if _, err := io.Copy(os.Stdout, pull); err != nil {
		return false, err
	}
	return true, nil
}
