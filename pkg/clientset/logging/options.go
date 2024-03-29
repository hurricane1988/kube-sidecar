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

package logging

// Options 日志日志结构体
type Options struct {
	LogPath    string `json:"logPath,omitempty" yaml:"logPath,omitempty" xml:"logPath,omitempty" description:"日志写入路径" example:"/tmp"`
	WriteLog   bool   `json:"writeLog,omitempty" yaml:"writeLog,omitempty" xml:"writeLog,omitempty"description:"是否将日志写入文件中"`
	MaxSize    int    `json:"maxSize,omitempty" yaml:"maxSize,omitempty" xml:"maxSize,omitempty" description:"日志文件最大大小，单位MB"`
	MaxBackups int    `json:"maxBackups,omitempty" yaml:"maxBackups,omitempty" xml:"maxBackups,omitempty" description:"日志文件保留数量"`
	MaxAge     int    `json:"maxAge,omitempty" yaml:"maxAge,omitempty" xml:"maxAge,omitempty" description:"日志文件保留天数"`
	Filename   string `json:"filename,omitempty" yaml:"filename,omitempty" xml:"filename,omitempty" description:"日志文件名称"`
}

// NewLoggingOptions 日志配置
func NewLoggingOptions() *Options {
	return &Options{
		LogPath:    "/tmp",
		Filename:   "kube-sidecar.log",
		WriteLog:   true,
		MaxSize:    10,
		MaxBackups: 10,
		MaxAge:     10,
	}
}
