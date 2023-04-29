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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"strconv"
	"strings"

	"kube-sidecar/config"
	"kube-sidecar/pkg/model/container"
	s "kube-sidecar/pkg/model/secret"
	"kube-sidecar/utils/tools"

	lg "kube-sidecar/utils/logging"

	"kube-sidecar/utils/clients/k8s"
)

// 定义全局annotation key值
const (
	annotationKey = "deployment.kubernetes.io/sidecar"
)

// AddDeploymentSidecar 为deployment添加sidecar容器方法
func AddDeploymentSidecar(deployment *appsv1.Deployment, client k8s.Client) error {
	// 定义全局错误信息
	var (
		errMsg error
		c      = config.Config
	)
	// 创建sidecar容器对象
	sidecar := container.CreateContainer()
	// 增加sidecar容器到deployment
	deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, *sidecar)
	// 获取sidecar 后端存储类型
	interval, _ := strconv.Atoi(tools.SetDefaultValueNotExist(
		deployment.Annotations["deployment.kubernetes.io/sidecar.inputRefreshInterval"],
		string(rune(c.FluentBitConfig.InputRefreshInterval))))

	fluent := s.FluentBitConf{
		ServiceLogLevel: tools.SetDefaultValueNotExist(
			deployment.Annotations["deployment.kubernetes.io/sidecar.serviceLogLevel"],
			c.FluentBitConfig.ServiceLogLevel),
		InputAppName: deployment.Name,
		InputLogPath: deployment.Annotations["deployment.kubernetes.io/sidecar.inputLogPath"],
		// InputAppTag:  deployment.Name,
		InputMemBufLimit: tools.SetDefaultValueNotExist(
			deployment.Annotations["deployment.kubernetes.io/sidecar.inputMemBufLimit"],
			c.FluentBitConfig.InputMemBufLimit), InputRefreshInterval: interval,
		OutputEsHost:        deployment.Annotations["deployment.kubernetes.io/sidecar.outputEsHost"],
		OutputEsPort:        deployment.Annotations["deployment.kubernetes.io/sidecar.outputEsHost"],
		OutputEsIndex:       deployment.Annotations["deployment.kubernetes.io/sidecar.outputEsIndex"],
		OutputEsUser:        deployment.Annotations["deployment.kubernetes.io/sidecar.outputEsUser"],
		OutputEsPassword:    deployment.Annotations["deployment.kubernetes.io/sidecar.outputEsPassword"],
		OutputKafkaHost:     deployment.Annotations["deployment.kubernetes.io/sidecar.outputKafkaHost"],
		OutputKafkaPort:     deployment.Annotations["deployment.kubernetes.io/sidecar.outputKafkaPort"],
		OutputKafkaTopic:    deployment.Annotations["deployment.kubernetes.io/sidecar.outputKafkaTopic"],
		OutputKafkaUser:     deployment.Annotations["deployment.kubernetes.io/sidecar.outputKafkaUser"],
		OutputKafkaPassword: deployment.Annotations["deployment.kubernetes.io/sidecar.outputKafkaPassword"],
	}
	// 获取fluentBit output类型
	backendType := deployment.Annotations["deployment.kubernetes.io/sidecar.backend"]
	// 基于backendType创建不同的secret配置
	err := s.CreateFluentBitSecret(backendType, deployment.Name, deployment.Namespace, client, fluent)
	if err != nil {
		lg.Logger.Error(err.Error())
		errMsg = err
	}
	// 创建secret volume对象
	secretVolume := corev1.Volume{
		Name: c.Sidecar.VolumeName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: strings.Join([]string{deployment.Name, "sidecar"}, "-"),
			},
		},
	}
	// 添加卷至Deployment
	deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, secretVolume)
	// 更新Deployment object添加新的sidecar容器和卷
	_, err = client.Kubernetes().AppsV1().Deployments(deployment.Namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		lg.Logger.Error("更新deployment " + deployment.Name + "失败,错误信息," + err.Error())
		errMsg = err
	}
	lg.Logger.Info("更新deployment " + deployment.Name + "成功!")
	return errMsg
}

// WatchDeployment watching kubernetes deployment changes
func WatchDeployment(client k8s.Client) {
	// 定义全局Deployment对象
	var (
		c = config.Config
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
			if annotations[annotationKey] == "true" && tools.WhetherExists(
				dp.Namespace, c.NamespacesWhiteList.Names) == false && tools.WhetherExists(
				dp.Name, c.DeploymentWhiteList.Names) == false {
				// 执行自动添加sidecar容器
				err = AddDeploymentSidecar(dp, client)
				if err != nil {
					lg.Logger.Error(dp.Name + " 自动添加sidecar容器镜像失败,错误信息," + err.Error())
				}
			}
			continue
		}
	}
}
