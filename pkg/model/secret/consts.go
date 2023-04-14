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
)

// fluentBit默认常量
const fluentBit = `
[SERVICE]
    Flush 1
    Log_Level {{.ServiceLogLevel}}
    Parsers_File parsers.conf
    HTTP_Server on
    HTTP_Listen 0.0.0.0
    HTTP_Port 2020
[INPUT]
    Name tail
    Path {{.InputLogPath}}/*.log
    Parser docker
    Tag {{.InputAppName}}.log
    DB /var/log/{{.InputAppName}}/flb-db
    Mem_Buf_Limit {{.InputMemBufLimit}}
    Skip_Long_Lines On
    Refresh_Interval {{.InputRefreshInterval}}
[OUTPUT]
    Name es
    Match {{.InputAppName}}.log
    Host {{.OutputEsHost}}
    Port {{.OutputEsPort}}
    Index {{.OutputEsIndex}}
    User {{.OutputEsUser}}
    Password {{.OutputEsPassword}}
[OUTPUT]
    Name kafka
    Match {{.InputAppName}}.log
    brokers {{.OutputKafkaHost}}:{{.OutputKafkaPort}}
    Topic {{.OutputKafkaTopic}}
    User {{.OutputKafkaUser}}
    Password {{.OutputKafkaPassword}}
    Security_Protocol SASL_PLAINTEXT
`

// FluentBitConf 定义FluentBit配置结构体
type FluentBitConf struct {
    ServiceLogLevel      string `json:"serviceLogLevel,omitempty" yaml:"serviceLogLevel" xml:"serviceLogLevel" describe:"fluentBit日志level,默认info"`
    InputAppName         string `json:"inputAppName,omitempty" yaml:"inputAppName" xml:"inputAppName" describe:"采集日志的应用名称"`
    InputLogPath         string `json:"inputLogPath,omitempty" yaml:"inputLogPath" xml:"inputLogPath" describe:"采集日志路劲"`
    InputAppTag          string `json:"inputAppTag,omitempty" yaml:"inputAppTag" xml:"inputAppTag" describe:"采集日志的应用Tag"`
    InputMemBufLimit     string `json:"inputMemBufLimit,omitempty" yaml:"inputMemBufLimit" xml:"inputMemBufLimit" describe:"采集日志的应用Tag"`
    InputRefreshInterval int    `json:"inputRefreshInterval,omitempty" yaml:"inputRefreshInterval" xml:"inputRefreshInterval" describe:"采集日志刷新间隔"`
    OutputEsHost         string `json:"outputEsHost,omitempty" yaml:"outputEsHost" xml:"outputEsHost" describe:"elasticsearch数据库地址"`
    OutputEsPort         string `json:"outputEsPort,omitempty" yaml:"outputEsPort" xml:"outputEsPort" describe:"elasticsearch数据库端口"`
    OutputEsIndex        string `json:"outputEsIndex,omitempty" yaml:"outputEsIndex" xml:"outputEsIndex" describe:"elasticsearch数据库index"`
    OutputEsUser         string `json:"outputEsUser,omitempty" yaml:"outputEsUser" xml:"outputEsUser" describe:"elasticsearch数据库user"`
    OutputEsPassword     string `json:"outputEsPassword,omitempty" yaml:"outputEsPassword" xml:"outputEsPassword" describe:"elasticsearch数据库password"`
    OutputKafkaHost      string `json:"outputKafkaHost,omitempty" yaml:"outputKafkaHost" xml:"outputKafkaHost" describe:"kafka数据库地址"`
    OutputKafkaPort      string `json:"outputKafkaPort,omitempty" yaml:"outputKafkaPort" xml:"outputKafkaPort" describe:"kafka数据库端口"`
    OutputKafkaTopic     string `json:"outputKafkaTopic,omitempty" yaml:"outputKafkaTopic" xml:"outputKafkaTopic" describe:"kafka topic"`
    OutputKafkaUser      string `json:"outputKafkaUser,omitempty" yaml:"outputKafkaUser" xml:"outputKafkaUser" describe:"kafka数据库user"`
    OutputKafkaPassword  string `json:"outputKafkaPassword,omitempty" yaml:"outputKafkaPassword" xml:"outputKafkaPassword" describe:"kafka数据库password"`
}

// TextToSecret  创建用户生产secret base64加密的方法
func TextToSecret(plainText []byte) string {
    // 将byte编码为base64并返回
    return base64.StdEncoding.EncodeToString(plainText)
}

// FluentBitTemplate 根据配置参数生产fluentBit文件模版
func FluentBitTemplate(data FluentBitConf) ([]byte, error) {
    // 解析模版
    tpl, err := template.New("fluentBit").Parse(fluentBit)
    if err != nil {
        panic(err)
        return nil, err
    }

    // 设置bytes.Buffer对象
    var buf bytes.Buffer
    // 应用模版并输出保存文件中
    err = tpl.Execute(&buf, data)
    if err != nil {
        panic(err)
        return nil, err
    }
    // buf.Bytes()将bytes.Buffer转为[]byte
    return buf.Bytes(), nil
}
