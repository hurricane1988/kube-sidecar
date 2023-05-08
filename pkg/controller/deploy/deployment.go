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

package deploy

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	"kube-sidecar/config"
	"kube-sidecar/pkg/model/deploy"
	"kube-sidecar/utils/clients/k8s"
	lg "kube-sidecar/utils/logging"
	"kube-sidecar/utils/tools"
)

// 定义全局annotation key值
const (
	annotationKey = "deployment.kubernetes.io/sidecar"
)

// WatchDeployment watching kubernetes deployment changes
func WatchDeployment(client k8s.Client) {
	// 定义全局Deployment对象
	var (
		cfg = config.Config
	)
	// 创建watchInterface接口
	watchInterface, err := client.Kubernetes().AppsV1().Deployments("").Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		lg.Logger.Error("创建Deployment的watch失败,错误信息" + err.Error())
	}
	// 开始执行watching deployment
	for event := range watchInterface.ResultChan() {
		// 检查操作事件是否为修改或者新增
		if event.Type == watch.Modified || event.Type == watch.Added {
			// 获取更新的Deployment 对象
			dp := event.Object.(*appsv1.Deployment)
			// 检查deployment是否有required annotation
			annotations := dp.GetAnnotations()
			if annotations[annotationKey] == "true" &&
				// 判断该工作负载是否已经在namespace白名单中
				tools.WhetherExists(dp.Namespace, cfg.NamespacesWhiteList.Names) == false &&
				// 判断该工作负载是否已经在deployment白名单中
				tools.WhetherExists(dp.Name, cfg.DeploymentWhiteList.Names) == false &&
				// 判断该deployment是否已经注入sidecar容器
				tools.WhetherExists(cfg.Sidecar.Name, tools.WorkloadContainerNames("Deployment", event.Object)) == false {
				// 执行自动添加sidecar容器
				err = deploy.AddDeploymentSidecar(dp, client)
				if err != nil {
					lg.Logger.Error(dp.Name + " 自动添加sidecar容器镜像失败,错误信息," + err.Error())
				}
			}
			continue
		}
	}
}
