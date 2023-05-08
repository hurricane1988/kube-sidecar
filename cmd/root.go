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
	"log"
	"os"
)

// 定义rootCmd
var rootCmd = &cobra.Command{
	Use:   "kube-sidecar",
	Short: "kube-sidecar is a service that inject sidecar container into deployment.",
	Long: `============================================================================================================
        # Use watch method to watch the deployment when Annotations "deployment.kubernetes.io/sidecar": 'true'         #
        # Create secret of FluentBit config and mount on sidecar container  automatically                              #
        # Copyright 2023 QiMing Kubernetes Platform Authors                                                            #
        ============================================================================================================`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute 定义cobra执行器
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
