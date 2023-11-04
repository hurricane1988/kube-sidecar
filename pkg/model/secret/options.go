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
	"kube-sidecar/pkg/clientset/elastic"
	"kube-sidecar/pkg/clientset/fluent"
	"kube-sidecar/pkg/clientset/kafka"
	"kube-sidecar/pkg/clientset/logging"
	"text/template"
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

// FluentBitService FluentBit Service配置块结构体
// type FluentBitService struct {
// 	ServiceLogLevel string `json:"serviceLogLevel,omitempty" yaml:"serviceLogLevel,omitempty" xml:"serviceLogLevel,omitempty" describe:"fluentBit日志level,默认info"`
// }

// TextToSecret  创建用户生产secret base64加密的方法
func TextToSecret(plainText []byte) string {
	// 将byte编码为base64并返回
	return base64.StdEncoding.EncodeToString(plainText)
}

// FluentBitTemplate 根据配置参数生产fluentBit文件模版
func FluentBitTemplate(backend string, fluent fluent.Options) ([]byte, error) {
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
		k := kafka.OutputKafka{
			ServiceLogLevel:      fluent.ServiceLogLevel,
			InputLogPath:         fluent.InputLogPath,
			InputAppName:         fluent.InputAppName,
			InputMemBufLimit:     fluent.InputMemBufLimit,
			InputRefreshInterval: fluent.InputRefreshInterval,
			OutputKafkaHost:      fluent.OutputKafkaHost,
			OutputKafkaPort:      fluent.OutputKafkaPort,
			OutputKafkaUser:      fluent.OutputKafkaUser,
			OutputKafkaTopic:     fluent.OutputKafkaTopic,
			OutputKafkaPassword:  fluent.OutputKafkaPassword,
		}
		tpl, err = template.New("fluentBit-kafka").Parse(fluentBitKafka)
		if err != nil {
			logging.Logger.Error(err.Error())
			return nil, err
		}
		err = tpl.Execute(&buf, k)
		if err != nil {
			logging.Logger.Error(err.Error())
			return nil, err
		}
	case "elasticsearch":
		// 应用模版并输出保存文件中
		r := elastic.OutputElasticsearch{
			ServiceLogLevel:      fluent.ServiceLogLevel,
			InputLogPath:         fluent.InputLogPath,
			InputAppName:         fluent.InputAppName,
			InputMemBufLimit:     fluent.InputMemBufLimit,
			InputRefreshInterval: fluent.InputRefreshInterval,
			OutputEsHost:         fluent.OutputEsHost,
			OutputEsPort:         fluent.OutputEsPort,
			OutputEsIndex:        fluent.OutputEsIndex,
			OutputEsUser:         fluent.OutputEsUser,
			OutputEsPassword:     fluent.OutputEsPassword,
		}
		tpl, err = template.New("fluentBit-elasticsearch").Parse(fluentBitES)
		if err != nil {
			logging.Logger.Error(err.Error())
			return nil, err
		}
		err = tpl.Execute(&buf, r)
		if err != nil {
			logging.Logger.Error(err.Error())
			return nil, err
		}

	default:
		tpl, err = template.New("fluentBit").Parse(fluentBit)
		if err != nil {
			logging.Logger.Error(err.Error())
			return nil, err
		}
	}
	// buf.Bytes()将bytes.Buffer转为[]byte
	return buf.Bytes(), nil
}
