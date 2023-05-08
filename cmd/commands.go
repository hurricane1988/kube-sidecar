/*
Copyright 2023 QKP Authors

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

package cmd

import (
	"github.com/spf13/cobra"

	"kube-sidecar/config"
	"kube-sidecar/pkg/model/deploy"
	"kube-sidecar/utils/clients/k8s"
	lg "kube-sidecar/utils/logging"
)

// StartKubeSidecar 启动kube-sideacar服务
var StartKubeSidecar = &cobra.Command{
	Use:              "start",
	Version:          config.Config.Version,
	Example:          "kube-sidecar start",
	Short:            "Start the kube-sidecar service",
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {
		var param string
		switch len(args) {
		case 1:
			param = args[0]
		default:
			param = "start"
		}
		if param != "start" {
			lg.Logger.Error("输入参数错误")
		}
		options := k8s.NewKubernetesOptions()
		client, _ := k8s.NewKubernetesClient(options)
		deploy.WatchDeployment(client)
	},
}

// 注册到rootCmd
func init() {
	rootCmd.AddCommand(StartKubeSidecar)
}
