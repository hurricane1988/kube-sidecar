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

package sidecar

// Options 定义sidecar容器配置信息
type Options struct {
	Name            string `json:"name,omitempty" xml:"name,omitempty" yaml:"name,omitempty"`
	Image           string `json:"image,omitempty" xml:"image,omitempty" yaml:"image,omitempty"`
	ImagePullPolicy string `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty" xml:"imagePullPolicy,omitempty"`
	RequestsCPU     string `json:"requestsCPU,omitempty" yaml:"requestsCPU,omitempty" xml:"requestsCPU,omitempty"`
	RequestsMemory  string `json:"requestsMemory,omitempty" yaml:"requestsMemory,omitempty" xml:"requestsMemory,omitempty"`
	LimitCPU        string `json:"limitCPU,omitempty" yaml:"limitCPU,omitempty" xml:"limitCPU,omitempty"`
	LimitMemory     string `json:"limitMemory,omitempty" yaml:"limitMemory,omitempty" xml:"limitMemory,omitempty"`
	VolumeName      string `json:"volumeName,omitempty" yaml:"volumeName,omitempty" xml:"volumeName,omitempty"`
	VolumeMount     string `json:"volumeMount,omitempty" yaml:"volumeMount,omitempty" xml:"volumeMount,omitempty"`
	ReadOnly        bool   `json:"readOnly,omitempty" yaml:"readOnly,omitempty" xml:"readOnly,omitempty"`
}

// NewSidecarOptions 容器配置
func NewSidecarOptions() *Options {
	return &Options{
		Name:            "sidecar",
		Image:           "fluent/fluent-bit:2.1.0",
		ImagePullPolicy: "IfNotPresent",
		RequestsCPU:     "250m",
		RequestsMemory:  "512Mi",
		LimitCPU:        "250m",
		LimitMemory:     "512Mi",
		ReadOnly:        true,
	}
}
