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
	"kube-sidecar/pkg/clientset/sidecar"
)

type container struct {
	sidecar sidecar.Options
}

type Container interface {
	Create() *corev1.Container
}

func NewContainer(sidecar sidecar.Options) Container {
	return &container{
		sidecar: sidecar,
	}
}

// Create 创建容器方法
func (s *container) Create() *corev1.Container {
	return &corev1.Container{
		Name:            s.sidecar.Name,
		Image:           s.sidecar.Image,
		ImagePullPolicy: corev1.PullPolicy(s.sidecar.ImagePullPolicy),
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      s.sidecar.VolumeName,
				MountPath: s.sidecar.VolumeMount,
			},
		},
		// 设置容器的resource资源
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(s.sidecar.RequestsCPU),
				corev1.ResourceMemory: resource.MustParse(s.sidecar.RequestsMemory),
			},
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(s.sidecar.LimitCPU),
				corev1.ResourceMemory: resource.MustParse(s.sidecar.LimitMemory),
			},
		},
	}
}
