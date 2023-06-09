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

package secret

import (
	"bytes"
	"encoding/base64"
	"text/template"

	lg "kube-sidecar/utils/logging"
)

// fluentBit全局配置
const (
	fluentBit = `
[SERVICE]
    HTTP_Server on
    HTTP_Listen 0.0.0.0
    HTTP_Port 2020
	Health_Check On
	HC_Errors_Count 5
	HC_Retry_Failure_Count 5
	HC_Period 5
	Log_Level {{.ServiceLogLevel}}
[INPUT]
    Name tail
    Path {{.InputLogPath}}/*.logging
    Parser docker
    Tag {{.InputAppName}}.logging
    DB /var/logging/{{.InputAppName}}/flb-db
    Mem_Buf_Limit {{.InputMemBufLimit}}
    Skip_Long_Lines On
    Refresh_Interval {{.InputRefreshInterval}}
[OUTPUT]
    Name es
    Match {{.InputAppName}}.logging
    Host {{.OutputEsHost}}
    Port {{.OutputEsPort}}
    Index {{.OutputEsIndex}}
    User {{.OutputEsUser}}
    Password {{.OutputEsPassword}}
[OUTPUT]
    Name kafka
    Match {{.InputAppName}}.logging
    brokers {{.OutputKafkaHost}}:{{.OutputKafkaPort}}
    Topic {{.OutputKafkaTopic}}
    User {{.OutputKafkaUser}}
    Password {{.OutputKafkaPassword}}
    Security_Protocol SASL_PLAINTEXT
`
	fluentBitKafka = `
[SERVICE]
	HTTP_Server on
	HTTP_Listen 0.0.0.0
	HTTP_Port 2020
	Health_Check On
	HC_Errors_Count 5
	HC_Retry_Failure_Count 5
	HC_Period 5
	Log_Level {{.ServiceLogLevel}}
[INPUT]
	Name tail
	Path {{.InputLogPath}}/*.logging
	Parser docker
	Tag {{.InputAppName}}.logging
	DB /var/logging/{{.InputAppName}}/flb-db
	Mem_Buf_Limit {{.InputMemBufLimit}}
	Skip_Long_Lines On
	Refresh_Interval {{.InputRefreshInterval}}
[OUTPUT]
	Name kafka
	Match {{.InputAppName}}.logging
	brokers {{.OutputKafkaHost}}:{{.OutputKafkaPort}}
	Topic {{.OutputKafkaTopic}}
	User {{.OutputKafkaUser}}
	Password {{.OutputKafkaPassword}}
	Security_Protocol SASL_PLAINTEXT
`
	fluentBitES = `
[SERVICE]
	HTTP_Server on
	HTTP_Listen 0.0.0.0
	HTTP_Port 2020
	Health_Check On
	HC_Errors_Count 5
	HC_Retry_Failure_Count 5
	HC_Period 5
	Log_Level {{.ServiceLogLevel}}
[INPUT]
	Name tail
	Path {{.InputLogPath}}/*.logging
	Parser docker
	Tag {{.InputAppName}}.logging
	DB /var/logging/{{.InputAppName}}/flb-db
	Mem_Buf_Limit {{.InputMemBufLimit}}
	Skip_Long_Lines On
	Refresh_Interval {{.InputRefreshInterval}}
[OUTPUT]
	Name es
	Match {{.InputAppName}}.logging
	Host {{.OutputEsHost}}
	Port {{.OutputEsPort}}
	Index {{.OutputEsIndex}}
	User {{.OutputEsUser}}
	Password {{.OutputEsPassword}}
`
)

// OutputKafka output为kafka的fluentBit配置结构体
type OutputKafka struct {
	ServiceLogLevel string `json:"serviceLogLevel,omitempty" yaml:"serviceLogLevel,omitempty" xml:"serviceLogLevel,omitempty" describe:"fluentBit日志level,默认info"`
	// Service             FluentBitService `json:"service,omitempty" xml:"service,omitempty" yaml:"service,omitempty" describe:"fluentBit日志level,默认info"`
	// Input               FluentBitInput `json:"input,omitempty" xml:"input,omitempty" yaml:"input,omitempty" describe:"fluentBit INPUT"`
	InputAppName         string `json:"inputAppName,omitempty" yaml:"inputAppName,omitempty" xml:"inputAppName,omitempty" describe:"采集日志的应用名称"`
	InputLogPath         string `json:"inputLogPath,omitempty" yaml:"inputLogPath,omitempty" xml:"inputLogPath,omitempty" describe:"采集日志路劲"`
	InputAppTag          string `json:"inputAppTag,omitempty" yaml:"inputAppTag,omitempty" xml:"inputAppTag,omitempty" describe:"采集日志的应用Tag"`
	InputMemBufLimit     string `json:"inputMemBufLimit,omitempty" yaml:"inputMemBufLimit,omitempty" xml:"inputMemBufLimit,omitempty" describe:"采集日志的应用Tag"`
	InputRefreshInterval int    `json:"inputRefreshInterval,omitempty" yaml:"inputRefreshInterval,omitempty" xml:"inputRefreshInterval,omitempty" describe:"采集日志刷新间隔"`
	OutputKafkaHost      string `json:"outputKafkaHost,omitempty" yaml:"outputKafkaHost,omitempty" xml:"outputKafkaHost,omitempty" describe:"kafka数据库地址"`
	OutputKafkaPort      string `json:"outputKafkaPort,omitempty" yaml:"outputKafkaPort,omitempty" xml:"outputKafkaPort,omitempty" describe:"kafka数据库端口"`
	OutputKafkaTopic     string `json:"outputKafkaTopic,omitempty" yaml:"outputKafkaTopic,omitempty" xml:"outputKafkaTopic,omitempty" describe:"kafka topic"`
	OutputKafkaUser      string `json:"outputKafkaUser,omitempty" yaml:"outputKafkaUser,omitempty" xml:"outputKafkaUser,omitempty" describe:"kafka数据库user"`
	OutputKafkaPassword  string `json:"outputKafkaPassword,omitempty" yaml:"outputKafkaPassword,omitempty" xml:"outputKafkaPassword,omitempty" describe:"kafka数据库password"`
}

// OutputElasticsearch output为elasticsearch的fluentBit配置结构体
type OutputElasticsearch struct {
	ServiceLogLevel string `json:"serviceLogLevel,omitempty" yaml:"serviceLogLevel,omitempty" xml:"serviceLogLevel,omitempty" describe:"fluentBit日志level,默认info"`
	// Service          FluentBitService `json:"service,omitempty" xml:"service,omitempty" yaml:"service,omitempty" describe:"fluentBit日志level,默认info"`
	// Input            FluentBitInput `json:"input,omitempty" xml:"input,omitempty" yaml:"input,omitempty" describe:"fluentBit INPUT"`
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
}

// FluentBitConf 定义FluentBit配置结构体
type FluentBitConf struct {
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

// FluentBitService FluentBit Service配置块结构体
// type FluentBitService struct {
// 	ServiceLogLevel string `json:"serviceLogLevel,omitempty" yaml:"serviceLogLevel,omitempty" xml:"serviceLogLevel,omitempty" describe:"fluentBit日志level,默认info"`
// }

// FluentBitInput FluentBitInput配置块结构体
type FluentBitInput struct {
	InputAppName         string `json:"inputAppName,omitempty" yaml:"inputAppName,omitempty" xml:"inputAppName,omitempty" describe:"采集日志的应用名称"`
	InputLogPath         string `json:"inputLogPath,omitempty" yaml:"inputLogPath,omitempty" xml:"inputLogPath,omitempty" describe:"采集日志路劲"`
	InputAppTag          string `json:"inputAppTag,omitempty" yaml:"inputAppTag,omitempty" xml:"inputAppTag,omitempty" describe:"采集日志的应用Tag"`
	InputMemBufLimit     string `json:"inputMemBufLimit,omitempty" yaml:"inputMemBufLimit,omitempty" xml:"inputMemBufLimit,omitempty" describe:"采集日志的应用Tag"`
	InputRefreshInterval int    `json:"inputRefreshInterval,omitempty" yaml:"inputRefreshInterval,omitempty" xml:"inputRefreshInterval,omitempty" describe:"采集日志刷新间隔"`
}

// TextToSecret  创建用户生产secret base64加密的方法
func TextToSecret(plainText []byte) string {
	// 将byte编码为base64并返回
	return base64.StdEncoding.EncodeToString(plainText)
}

// FluentBitTemplate 根据配置参数生产fluentBit文件模版
func FluentBitTemplate(backend string, data FluentBitConf) ([]byte, error) {
	// 解析模版
	var (
		tpl *template.Template
		err error
		// 设置bytes.Buffer对象
		buf bytes.Buffer
	)
	switch backend {
	case "kafka":
		// 应用模版并输出保存文件中
		k := OutputKafka{
			ServiceLogLevel:      data.ServiceLogLevel,
			InputLogPath:         data.InputLogPath,
			InputAppName:         data.InputAppName,
			InputMemBufLimit:     data.InputMemBufLimit,
			InputRefreshInterval: data.InputRefreshInterval,
			OutputKafkaHost:      data.OutputKafkaHost,
			OutputKafkaPort:      data.OutputKafkaPort,
			OutputKafkaUser:      data.OutputKafkaUser,
			OutputKafkaTopic:     data.OutputKafkaTopic,
			OutputKafkaPassword:  data.OutputKafkaPassword,
		}
		tpl, err = template.New("fluentBit-kafka").Parse(fluentBitKafka)
		if err != nil {
			lg.Logger.Error(err.Error())
			return nil, err
		}
		err = tpl.Execute(&buf, k)
		if err != nil {
			lg.Logger.Error(err.Error())
			return nil, err
		}
	case "elasticsearch":
		// 应用模版并输出保存文件中
		r := OutputElasticsearch{
			ServiceLogLevel:      data.ServiceLogLevel,
			InputLogPath:         data.InputLogPath,
			InputAppName:         data.InputAppName,
			InputMemBufLimit:     data.InputMemBufLimit,
			InputRefreshInterval: data.InputRefreshInterval,
			OutputEsHost:         data.OutputEsHost,
			OutputEsPort:         data.OutputEsPort,
			OutputEsIndex:        data.OutputEsIndex,
			OutputEsUser:         data.OutputEsUser,
			OutputEsPassword:     data.OutputEsPassword,
		}
		tpl, err = template.New("fluentBit-elasticsearch").Parse(fluentBitES)
		if err != nil {
			lg.Logger.Error(err.Error())
			return nil, err
		}
		err = tpl.Execute(&buf, r)
		if err != nil {
			lg.Logger.Error(err.Error())
			return nil, err
		}

	default:
		tpl, err = template.New("fluentBit").Parse(fluentBit)
		if err != nil {
			lg.Logger.Error(err.Error())
			return nil, err
		}
	}
	// buf.Bytes()将bytes.Buffer转为[]byte
	return buf.Bytes(), nil
}
