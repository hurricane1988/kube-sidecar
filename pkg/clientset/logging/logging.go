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

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"time"
)

// TimeFormat 设置全局时间格式
const (
	TimeFormat = "2006-01-02 15:04:05.000"
)

// Logger 定义全局logger变量
var Logger *zap.Logger

type Logging interface {
	Encoder() zapcore.Encoder
	ConsoleEncoder() zapcore.Encoder
	WriteSyncer() zapcore.WriteSyncer
	LevelEnabler() zapcore.Level
	EncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder)
	EncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder)
	EncodeCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder)
}

func NewLogging(option Options) Logging {
	return &option
}

func (o *Options) Logger() *zap.Logger {
	// 定义全局的newCore
	var (
		newCore = zapcore.NewTee()
	)
	Encoder := o.Encoder()
	WriteSyncer := o.WriteSyncer()
	LevelEnabler := o.LevelEnabler()
	ConsoleEncoder := o.ConsoleEncoder()
	switch o.WriteLog {
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
	return Logger
}

// Encoder 自定义的Encoder
func (o *Options) Encoder() zapcore.Encoder {
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
			EncodeLevel:    o.EncodeLevel,
			EncodeTime:     o.EncodeTime,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   o.EncodeCaller,
		})
}

// ConsoleEncoder 输出日志到控制台
func (o *Options) ConsoleEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
}

// WriteSyncer 自定义的WriteSyncer
func (o *Options) WriteSyncer() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   strings.Join([]string{o.LogPath, o.Filename}, "/"),
		MaxSize:    o.MaxSize,
		MaxBackups: o.MaxBackups,
		MaxAge:     o.MaxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// LevelEnabler 自定义的LevelEnabler
func (o *Options) LevelEnabler() zapcore.Level {
	return zapcore.InfoLevel
}

// EncodeLevel 自定义日志级别显示
func (o *Options) EncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

// EncodeTime 自定义时间格式显示
func (o *Options) EncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format(TimeFormat) + "]")
}

// EncodeCaller 自定义行号显示
func (o *Options) EncodeCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + caller.TrimmedPath() + "]")
}
