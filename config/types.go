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

package config

// Config 定义全局config结构体
type config struct {
	LoggingConfig       logging             `json:"loggingConfig,omitempty" yaml:"loggingConfig,omitempty" xml:"loggingConfig,omitempty"`
	Sidecar             sidecar             `json:"sidecar,omitempty" yaml:"sidecar,omitempty" xml:"sidecar,omitempty"`
	NamespacesWhiteList namespacesWhiteList `json:"namespacesWhiteList,omitempty" xml:"namespacesWhiteList,omitempty" yaml:"namespacesWhiteList,omitempty"`
	DeploymentWhiteList deploymentWhiteList `json:"deploymentWhiteList,omitempty" xml:"deploymentWhiteList,omitempty" yaml:"deploymentWhiteList,omitempty"`
	FluentBitConfig     fluentBitConfig     `json:"fluentBitConfig,omitempty" yaml:"fluentBitConfig,omitempty" xml:"fluentBitConfig,omitempty"`
}

// Logging 日志日志结构体
type logging struct {
	LogPath    string `json:"logPath,omitempty" yaml:"logPath,omitempty" xml:"logPath,omitempty" description:"日志写入路径" example:"/tmp"`
	WriteLog   bool   `json:"writeLog,omitempty" yaml:"writeLog,omitempty" xml:"writeLog,omitempty"description:"是否将日志写入文件中"`
	MaxSize    int    `json:"maxSize,omitempty" yaml:"maxSize,omitempty" xml:"maxSize,omitempty" description:"日志文件最大大小，单位MB"`
	MaxBackups int    `json:"maxBackups,omitempty" yaml:"maxBackups,omitempty" xml:"maxBackups,omitempty" description:"日志文件保留数量"`
	MaxAge     int    `json:"maxAge,omitempty" yaml:"maxAge,omitempty" xml:"maxAge,omitempty" description:"日志文件保留天数"`
	Filename   string `json:"filename,omitempty" yaml:"filename,omitempty" xml:"filename,omitempty" description:"日志文件名称"`
}

// sidecar 定义sidecar容器配置信息
type sidecar struct {
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

// 定义fluentBit全局配置结构体
type fluentBitConfig struct {
	ServiceLogLevel      string `json:"serviceLogLevel,omitempty" yaml:"serviceLogLevel,omitempty" xml:"serviceLogLevel,omitempty" describe:"fluentBit日志level,默认info"`
	InputMemBufLimit     string `json:"inputMemBufLimit,omitempty" yaml:"inputMemBufLimit,omitempty" xml:"inputMemBufLimit,omitempty" describe:"采集日志的应用Tag"`
	InputRefreshInterval int    `json:"inputRefreshInterval,omitempty" yaml:"inputRefreshInterval,omitempty" xml:"inputRefreshInterval,omitempty" describe:"采集日志刷新间隔"`
}

// 定义namespacesWhiteList结构体
type namespacesWhiteList struct {
	Names []string `json:"names,omitempty" xml:"names,omitempty" yaml:"names,omitempty"`
}

// deploymentWhiteList白名单
type deploymentWhiteList struct {
	Names []string `json:"names,omitempty" xml:"names,omitempty" yaml:"names,omitempty"`
}
