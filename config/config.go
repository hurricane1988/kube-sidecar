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

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/wonderivan/logger"
	"k8s.io/client-go/util/homedir"
	"kube-sidecar/pkg/clientset/fluent"
	"kube-sidecar/pkg/clientset/jaeger"
	"kube-sidecar/pkg/clientset/logging"
	"kube-sidecar/pkg/clientset/sidecar"
	"kube-sidecar/pkg/clientset/version"
	"kube-sidecar/pkg/clientset/workload"
	"path/filepath"
)

const (
	// WorkDir = "/opt"
	WorkDir = "/Users/hurricane/github/kube-sidecar"
)

// Config 定义全局config结构体
type Config struct {
	LoggingConfig   *logging.Options  `json:"loggingConfig,omitempty" yaml:"loggingConfig,omitempty" xml:"loggingConfig,omitempty" mapstructure:"loggingConfig"`
	JaegerConfig    *jaeger.Options   `yaml:"jaegerConfig,omitempty" xml:"jaegerConfig,omitempty" json:"jaegerConfig,omitempty" mapstructure:"jaegerConfig"`
	Sidecar         *sidecar.Options  `json:"sidecar,omitempty" yaml:"sidecar,omitempty" xml:"sidecar,omitempty" mapstructure:"sidecar"`
	WhiteList       *workload.Options `json:"whiteList,omitempty" xml:"whiteList,omitempty" yaml:"whiteList,omitempty" mapstructure:"whiteList"`
	FluentBitConfig *fluent.Options   `json:"fluentBitConfig,omitempty" yaml:"fluentBitConfig,omitempty" xml:"fluentBitConfig,omitempty" mapstructure:"fluentBitConfig"`
	Version         *version.Options  `json:"version,omitempty" xml:"version,omitempty" yaml:"version,omitempty" mapstructure:"version"`
}

// LoadConfigFromFile 初始化配置文件
func LoadConfigFromFile() (*Config, error) {
	// 获取当前程序执行路径
	// workDir, _ := os.Getwd()
	workDir := filepath.Join(WorkDir, "config", "conf")
	// 加载viper获取配置路径
	viper.AddConfigPath(workDir)
	viper.AddConfigPath(homedir.HomeDir())
	viper.AddConfigPath(".")
	// 设置读取配置文件
	// 设置读取的文件名
	viper.SetConfigName("config")
	// 设置读取的文件后缀
	viper.SetConfigType("yaml")
	// 匹配环境变量
	viper.AutomaticEnv()

	// 执行读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			logger.Error("配置文件" + workDir + "/conf.yaml 不存在!")
		}
	}
	conf := New()
	/* viper动态加载配置 */
	// 监视配置文件是否发生更改
	viper.WatchConfig()
	// 处理配置变更事件
	viper.OnConfigChange(func(in fsnotify.Event) {
		// TODO: 联调测试
		logger.Info("配置文件" + in.Name + "发生变更")
		logger.Warn(conf.LoggingConfig.WriteLog)
		// 重新初始化并生效配置
		err := viper.Unmarshal(conf)
		if err != nil {
			logger.Error("数据反序列化失败,错误信息", err)
			return
		}
	})
	if err := viper.Unmarshal(conf); err != nil {
		return nil, err
	}
	return conf, nil
}

func New() *Config {
	return &Config{
		LoggingConfig:   logging.NewLoggingOptions(),
		JaegerConfig:    jaeger.NewJaegerOptions(),
		Sidecar:         sidecar.NewSidecarOptions(),
		Version:         version.NewVersionOptions(),
		WhiteList:       workload.NewWhiteListOptions(),
		FluentBitConfig: fluent.NewFluentBitOptions(),
	}
}
