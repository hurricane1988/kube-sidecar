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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"kube-sidecar/pkg/clientset/fluent"
	"kube-sidecar/pkg/clientset/jaeger"
	"kube-sidecar/pkg/clientset/kubernetes"
	lg "kube-sidecar/pkg/clientset/logging"
	"kube-sidecar/pkg/clientset/sidecar"
	"kube-sidecar/pkg/clientset/workload"
	"kube-sidecar/pkg/model/deploy"
	"kube-sidecar/utils/tools"
)

// 定义全局annotation key值
const (
	annotationKey = "deployment.kubernetes.io/sidecar"
)

type deployment struct {
	K8sClient kubernetes.Client
	FluentBit fluent.Options
	Sidecar   sidecar.Options
	Jeager    jaeger.Options
	WhiteList workload.Options
}

type Deployment interface {
	Watch(ctx context.Context, tracerName, spanName string)
}

func NewDeployment(k8sClient kubernetes.Client, fluentBit fluent.Options, sidecar sidecar.Options, jeager jaeger.Options, whiteList workload.Options) Deployment {
	return &deployment{
		K8sClient: k8sClient,
		FluentBit: fluentBit,
		Sidecar:   sidecar,
		Jeager:    jeager,
		WhiteList: whiteList,
	}
}

// Watch watching kubernetes deployment changes
func (d *deployment) Watch(ctx context.Context, tracerName, spanName string) {
	// 增加链路跟踪
	if d.Jeager.Enable {
		tr := otel.Tracer(tracerName)
		_, span := tr.Start(ctx, spanName)
		span.SetAttributes(attribute.Key("update").String("deployment"))
		lg.Logger.Info("添加trace链路跟踪成功")
		defer span.End()
	}
	// 创建watchInterface接口
	watchInterface, err := d.K8sClient.Kubernetes().AppsV1().Deployments("").Watch(context.TODO(), metav1.ListOptions{})
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
				tools.WhetherExists(dp.Namespace, d.WhiteList.Namespaces) == false &&
				// 判断该工作负载是否已经在deployment白名单中
				tools.WhetherExists(dp.Name, d.WhiteList.Deployments) == false &&
				// 判断该deployment是否已经注入sidecar容器
				tools.WhetherExists(d.Sidecar.Name, tools.WorkloadContainerNames("Deployment", event.Object)) == false {
				// 执行自动添加sidecar容器
				err = deploy.NewDeploy(d.K8sClient, d.FluentBit).AddSidecar(dp)
				if err != nil {
					lg.Logger.Error(dp.Name + " 自动添加sidecar容器镜像失败,错误信息," + err.Error())
				}
			}
			continue
		}
	}
}
