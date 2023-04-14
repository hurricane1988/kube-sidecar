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
	"k8s.io/klog/v2"

	"kube-sidecar/utils/clients/k8s"
)

// AddDeploymentSidecar 为deployment添加sidecar容器方法
func AddDeploymentSidecar(container *corev1.Container, deployment *appsv1.Deployment) {
	// 新增sidecar container
	deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, *container)
}

// AddSidecarSecretVolume 为Deployment中的sidecar容器添加secret volume
func AddSidecarSecretVolume(ctx context.Context, client k8s.Client, volumeName, secretName, mountPath string, deployment *appsv1.Deployment) {
	// 为deployment添加Secret卷
	deployment.Spec.Template.Spec.Volumes = append(
		deployment.Spec.Template.Spec.Volumes,
		corev1.Volume{
			Name: volumeName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: secretName,
				},
			},
		})

	// 挂载卷到sidecar容器
	deployment.Spec.Template.Spec.Containers[len(deployment.Spec.Template.Spec.Containers)-1].VolumeMounts = append(
		deployment.Spec.Template.Spec.Containers[len(deployment.Spec.Template.Spec.Containers)-1].VolumeMounts,
		corev1.VolumeMount{
			Name:      volumeName,
			MountPath: mountPath,
		},
	)
	// 更新deployment
	_, err := client.Kubernetes().
		AppsV1().
		Deployments(deployment.Name).
		Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		klog.Errorf("更新Deployment "+deployment.Name+"失败,错误信息", err)
	}
}

// OnUpdate 创建处理Deployment更新方法
func OnUpdate(container corev1.Container) {
	const (
		annotationKey = "deployment.kubernetes.io/sidecar"
	)
	deployment := &appsv1.Deployment{}
	annotations := deployment.Annotations
	if val, ok := annotations[annotationKey]; ok && val == "true" {
		AddDeploymentSidecar(&container, deployment)
	}
}
