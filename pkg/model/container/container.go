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

package container

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"kube-sidecar/config"
)

// CreateContainer 创建容器方法
func CreateContainer() *corev1.Container {
	sidecar := config.Config.Sidecar
	return &corev1.Container{
		Name:            sidecar.Name,
		Image:           sidecar.Image,
		ImagePullPolicy: corev1.PullPolicy(sidecar.ImagePullPolicy),
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      sidecar.VolumeName,
				MountPath: sidecar.VolumeMount,
			},
		},
		// 设置容器的resource资源
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(sidecar.RequestsCPU),
				corev1.ResourceMemory: resource.MustParse(sidecar.RequestsMemory),
			},
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(sidecar.LimitCPU),
				corev1.ResourceMemory: resource.MustParse(sidecar.LimitMemory),
			},
		},
	}
}
