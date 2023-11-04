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
	"kube-sidecar/pkg/clientset/fluent"
	"kube-sidecar/pkg/clientset/kubernetes"
	"kube-sidecar/pkg/clientset/logging"
	"kube-sidecar/pkg/clientset/sidecar"
	"kube-sidecar/pkg/model/secret"
	"strconv"
	"strings"
	"time"

	"kube-sidecar/pkg/model/container"
	"kube-sidecar/utils/tools"
)

type deploy struct {
	k8sClient kubernetes.Client
	fluentBit fluent.Options
	sidecar   sidecar.Options
}

type Deploy interface {
	AddSidecar(deployment *appsv1.Deployment) error
}

func NewDeploy(k8sClient kubernetes.Client, fluentBit fluent.Options) Deploy {
	return &deploy{
		k8sClient: k8sClient,
		fluentBit: fluentBit,
	}
}

// AddSidecar 为deployment添加sidecar容器方法
func (d *deploy) AddSidecar(deployment *appsv1.Deployment) error {
	// 定义全局错误信息
	var (
		errMsg error
	)
	// 创建sidecar容器对象
	s := container.NewContainer(d.sidecar).Create()
	// 增加sidecar容器到deployment
	deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, *s)
	// 获取sidecar 后端存储类型
	interval, _ := strconv.Atoi(tools.SetDefaultValueNotExist(
		deployment.Annotations["deployment.kubernetes.io/sidecar.inputRefreshInterval"],
		string(rune(d.fluentBit.InputRefreshInterval))))
	f := fluent.Options{
		ServiceLogLevel: tools.SetDefaultValueNotExist(deployment.Annotations["deployment.kubernetes.io/sidecar.serviceLogLevel"], d.fluentBit.ServiceLogLevel),
		InputAppName:    deployment.Name,
		InputLogPath:    tools.SetDefaultValueNotExist(deployment.Annotations["deployment.kubernetes.io/sidecar.inputLogPath"], "/tmp"),
		// InputAppTag:  deployment.Name,
		InputMemBufLimit:     tools.SetDefaultValueNotExist(deployment.Annotations["deployment.kubernetes.io/sidecar.inputMemBufLimit"], d.fluentBit.InputMemBufLimit),
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
	err := secret.NewSecret(d.k8sClient).FluentBit(backendType, deployment.Name, deployment.Namespace, f)
	if err != nil {
		logging.Logger.Error(err.Error())
		errMsg = err
	}
	// 创建secret volume对象
	secretVolume := corev1.Volume{
		Name: d.sidecar.VolumeName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: strings.Join([]string{deployment.Name, "sidecar"}, "-"),
			},
		},
	}
	// 添加卷至Deployment
	deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, secretVolume)
	// 更新Deployment object添加新的sidecar容器和卷
	_, err = d.k8sClient.Kubernetes().AppsV1().Deployments(deployment.Namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		logging.Logger.Error("更新deployment " + deployment.Name + "失败,错误信息," + err.Error())
		errMsg = err
	}
	logging.Logger.Info("更新deployment " + deployment.Name + "成功!")
	return errMsg
}
