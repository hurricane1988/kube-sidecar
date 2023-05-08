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
	"strconv"
	"strings"
	"time"

	"kube-sidecar/config"
	"kube-sidecar/pkg/model/container"
	s "kube-sidecar/pkg/model/secret"
	"kube-sidecar/utils/tools"

	lg "kube-sidecar/utils/logging"

	"kube-sidecar/utils/clients/k8s"
)

// AddDeploymentSidecar 为deployment添加sidecar容器方法
func AddDeploymentSidecar(deployment *appsv1.Deployment, client k8s.Client) error {
	// 定义全局错误信息
	var (
		errMsg error
		cfg    = config.Config
	)
	// 创建sidecar容器对象
	sidecar := container.CreateContainer()
	// 增加sidecar容器到deployment
	deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, *sidecar)
	// 获取sidecar 后端存储类型
	interval, _ := strconv.Atoi(tools.SetDefaultValueNotExist(
		deployment.Annotations["deployment.kubernetes.io/sidecar.inputRefreshInterval"],
		string(rune(cfg.FluentBitConfig.InputRefreshInterval))))

	fluent := s.FluentBitConf{
		ServiceLogLevel: tools.SetDefaultValueNotExist(deployment.Annotations["deployment.kubernetes.io/sidecar.serviceLogLevel"], cfg.FluentBitConfig.ServiceLogLevel),
		InputAppName:    deployment.Name,
		InputLogPath:    tools.SetDefaultValueNotExist(deployment.Annotations["deployment.kubernetes.io/sidecar.inputLogPath"], "/tmp"),
		// InputAppTag:  deployment.Name,
		InputMemBufLimit:     tools.SetDefaultValueNotExist(deployment.Annotations["deployment.kubernetes.io/sidecar.inputMemBufLimit"], cfg.FluentBitConfig.InputMemBufLimit),
		InputRefreshInterval: interval,
		OutputEsHost:         deployment.Annotations["deployment.kubernetes.io/sidecar.outputEsHost"],
		OutputEsPort:         tools.SetDefaultValueNotExist(deployment.Annotations["deployment.kubernetes.io/sidecar.outputEsPort"], "9200"),
		OutputEsIndex:        tools.SetDefaultValueNotExist(deployment.Annotations["deployment.kubernetes.io/sidecar.outputEsIndex"], deployment.Name+time.Now().Format("2006-01-02")),
		OutputEsUser:         deployment.Annotations["deployment.kubernetes.io/sidecar.outputEsUser"],
		OutputEsPassword:     deployment.Annotations["deployment.kubernetes.io/sidecar.outputEsPassword"],
		OutputKafkaHost:      deployment.Annotations["deployment.kubernetes.io/sidecar.outputKafkaHost"],
		OutputKafkaPort:      tools.SetDefaultValueNotExist(deployment.Annotations["deployment.kubernetes.io/sidecar.outputKafkaPort"], "9092"),
		OutputKafkaTopic:     deployment.Annotations["deployment.kubernetes.io/sidecar.outputKafkaTopic"],
		OutputKafkaUser:      deployment.Annotations["deployment.kubernetes.io/sidecar.outputKafkaUser"],
		OutputKafkaPassword:  deployment.Annotations["deployment.kubernetes.io/sidecar.outputKafkaPassword"],
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
		Name: cfg.Sidecar.VolumeName,
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
