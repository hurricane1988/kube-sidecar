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
	"kube-sidecar/pkg/clientset/kubernetes"
	"kube-sidecar/pkg/clientset/logging"

	ot "kube-sidecar/utils/opentelemetry"
	"kube-sidecar/utils/tools"
)

// StartKubeSidecar 启动kube-sideacar服务
var StartKubeSidecar = &cobra.Command{
	Use:              "start",
	Example:          "kube-sidecar start",
	Short:            "Start the kube-sidecar service",
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			param       string
			tracerName  string = "kube-sidecar"
			spanName    string = "sidecar"
			service     string = "kube-sidecar"
			environment string = "qkp"
			id          int64  = tools.RandomInt64()
		)
		switch len(args) {
		case 1:
			param = args[0]
		default:
			param = "start"
		}
		if param != "start" {
			logging.Logger.Error("输入参数错误")
		}
		// 打印终端提示
		logging.Logger.Info("成功启动kube-sidecar监听服务")
		tools.TerminalColor()
		cfg, _ := config.LoadConfigFromFile()
		// 注册全局tracer
		options := kubernetes.NewKubernetesOptions()
		client, _ := kubernetes.NewKubernetesClient(options)
		ot.NewOpenTelemetry(client, *cfg.FluentBitConfig, *cfg.Sidecar, *cfg.JaegerConfig, *cfg.WhiteList).TracerProvider(service, environment)
	},
}

// 注册到rootCmd
func init() {
	rootCmd.AddCommand(StartKubeSidecar)
}
