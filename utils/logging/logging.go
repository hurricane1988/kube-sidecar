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

package middleware

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"time"

	"kube-sidecar/config"
)

// 设置全局时间格式
const (
	logTmFmtWithMS = "2006-01-02 15:04:05.000"
)

// Logger 定义全局logger变量
var Logger *zap.Logger

// 初始化NewLogger
func init() {
	NewLogger()
}

func NewLogger() {
	logPath := config.Config.LoggingConfig.LogPath
	writeable := config.Config.LoggingConfig.WriteLog
	// 定义全局的newCore
	var newCore = zapcore.NewTee()
	Encoder := GetEncoder()
	WriteSyncer := GetWriteSyncer(logPath)
	LevelEnabler := GetLevelEnabler()
	ConsoleEncoder := GetConsoleEncoder()
	switch writeable {
	case true:
		newCore = zapcore.NewTee(
			zapcore.NewCore(Encoder, WriteSyncer, LevelEnabler),                          // 写入文件
			zapcore.NewCore(ConsoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel), // 写入控制台
		)
	default:
		// 默认日志不写入文件中
		newCore = zapcore.NewTee(
			zapcore.NewCore(ConsoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel), // 写入控制台
		)
	}
	Logger = zap.New(newCore, zap.AddCaller())
	zap.ReplaceGlobals(Logger)
}

// GetEncoder 自定义的Encoder
func GetEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(
		zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller_line",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     "\n",
			EncodeLevel:    EncodeLevel,
			EncodeTime:     EncodeTime,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   EncodeCaller,
		})
}

// GetConsoleEncoder 输出日志到控制台
func GetConsoleEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
}

// GetWriteSyncer 自定义的WriteSyncer
func GetWriteSyncer(logPath string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   strings.Join([]string{logPath, config.Config.LoggingConfig.Filename}, "/"),
		MaxSize:    config.Config.LoggingConfig.MaxSize,
		MaxBackups: config.Config.LoggingConfig.MaxBackups,
		MaxAge:     config.Config.LoggingConfig.MaxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// GetLevelEnabler 自定义的LevelEnabler
func GetLevelEnabler() zapcore.Level {
	return zapcore.InfoLevel
}

// EncodeLevel 自定义日志级别显示
func EncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

// EncodeTime 自定义时间格式显示
func EncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format(logTmFmtWithMS) + "]")
}

// EncodeCaller 自定义行号显示
func EncodeCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + caller.TrimmedPath() + "]")
}
