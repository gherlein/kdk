// Copyright © 2018 Cisco Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/manifoldco/promptui"
)

func Destroy(ctx context.Context, dockerClient *client.Client, logger logrus.Entry) error {

	var containerIds []string

	containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{})

	if err != nil {
		logger.WithField("error", err).Fatal("Failed to list docker containers")
	}
	for _, container := range containers {
		for _, name := range container.Names {
			if strings.HasPrefix(name, "/kdk") {
				containerIds = append(containerIds, container.ID)
				break
			}
		}
	}
	if len(containerIds) > 0 {
		logger.Info("Destroying KDK container(s)...")
		for _, containerId := range containerIds {
			fmt.Printf("Delete KDK container [%v]\n", containerId[:8])
			prompt := promptui.Prompt{
				Label:     "Continue",
				IsConfirm: true,
			}
			if _, err := prompt.Run(); err != nil {
				logger.Error("KDK container deletion canceled or invalid input.")
				return nil
			}
			if err := dockerClient.ContainerRemove(ctx, containerId, types.ContainerRemoveOptions{Force: true}); err != nil {
				logger.WithField("error", err).Fatal("Failed to remove KDK container")
			}
		}
		logger.Info("KDK destroy complete.")
	} else {
		logger.Info("No KDK containers found. Nothing to destroy...")
	}
	return nil
}