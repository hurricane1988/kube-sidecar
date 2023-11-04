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

package fluent

// Options 定义FluentBit配置结构体
type Options struct {
	ServiceLogLevel string `json:"serviceLogLevel,omitempty" yaml:"serviceLogLevel,omitempty" xml:"serviceLogLevel,omitempty" describe:"fluentBit日志level,默认info"`
	// Service             FluentBitService `json:"service,omitempty" xml:"service,omitempty" yaml:"service,omitempty" describe:"fluentBit日志level,默认info"`
	// Input               FluentBitInput `json:"input,omitempty" xml:"input,omitempty" yaml:"input,omitempty" describe:"fluentBit INPUT"`
	InputAppName         string `json:"inputAppName,omitempty" yaml:"inputAppName,omitempty" xml:"inputAppName,omitempty" describe:"采集日志的应用名称"`
	InputLogPath         string `json:"inputLogPath,omitempty" yaml:"inputLogPath,omitempty" xml:"inputLogPath,omitempty" describe:"采集日志路劲"`
	InputAppTag          string `json:"inputAppTag,omitempty" yaml:"inputAppTag,omitempty" xml:"inputAppTag,omitempty" describe:"采集日志的应用Tag"`
	InputMemBufLimit     string `json:"inputMemBufLimit,omitempty" yaml:"inputMemBufLimit,omitempty" xml:"inputMemBufLimit,omitempty" describe:"采集日志的应用Tag"`
	InputRefreshInterval int    `json:"inputRefreshInterval,omitempty" yaml:"inputRefreshInterval,omitempty" xml:"inputRefreshInterval,omitempty" describe:"采集日志刷新间隔"`
	OutputEsHost         string `json:"outputEsHost,omitempty" yaml:"outputEsHost,omitempty" xml:"outputEsHost,omitempty" describe:"elasticsearch数据库地址"`
	OutputEsPort         string `json:"outputEsPort,omitempty" yaml:"outputEsPort,omitempty" xml:"outputEsPort,omitempty" describe:"elasticsearch数据库端口"`
	OutputEsIndex        string `json:"outputEsIndex,omitempty" yaml:"outputEsIndex,omitempty" xml:"outputEsIndex,omitempty" describe:"elasticsearch数据库index"`
	OutputEsUser         string `json:"outputEsUser,omitempty" yaml:"outputEsUser,omitempty" xml:"outputEsUser,omitempty" describe:"elasticsearch数据库user"`
	OutputEsPassword     string `json:"outputEsPassword,omitempty" yaml:"outputEsPassword,omitempty" xml:"outputEsPassword,omitempty" describe:"elasticsearch数据库password"`
	OutputKafkaHost      string `json:"outputKafkaHost,omitempty" yaml:"outputKafkaHost,omitempty" xml:"outputKafkaHost,omitempty" describe:"kafka数据库地址"`
	OutputKafkaPort      string `json:"outputKafkaPort,omitempty" yaml:"outputKafkaPort,omitempty" xml:"outputKafkaPort,omitempty" describe:"kafka数据库端口"`
	OutputKafkaTopic     string `json:"outputKafkaTopic,omitempty" yaml:"outputKafkaTopic,omitempty" xml:"outputKafkaTopic,omitempty" describe:"kafka topic"`
	OutputKafkaUser      string `json:"outputKafkaUser,omitempty" yaml:"outputKafkaUser,omitempty" xml:"outputKafkaUser,omitempty" describe:"kafka数据库user"`
	OutputKafkaPassword  string `json:"outputKafkaPassword,omitempty" yaml:"outputKafkaPassword,omitempty" xml:"outputKafkaPassword,omitempty" describe:"kafka数据库password"`
}

// NewFluentBitOptions 获取FluentBit配置方法
func NewFluentBitOptions() *Options {
	return &Options{
		ServiceLogLevel:      "info",
		InputMemBufLimit:     "20MB",
		InputRefreshInterval: 20,
	}
}
