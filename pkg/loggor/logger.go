package loggor

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type Log struct {
	internal *logrus.Logger
}

var Default = New("debug")

func New(level string) *Log {
	// 解析日志级别
	parseLevel, err := logrus.ParseLevel(level)
	if err != nil {
		fmt.Printf("解析日志级别失败: %v，使用默认级别 debug\n", err)
		parseLevel = logrus.DebugLevel
	}

	// 初始化 logrus.Logger
	logrusInst := logrus.New()
	logrusInst.Level = parseLevel
	logrusInst.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}
	logrusInst.Out = logrus.StandardLogger().Out

	return &Log{
		internal: logrusInst,
	}
}

func (l *Log) Debug(args ...interface{})                 { l.internal.Debug(args...) }
func (l *Log) Debugf(format string, args ...interface{}) { l.internal.Debugf(format, args...) }
func (l *Log) Debugln(args ...interface{})               { l.internal.Debugln(args...) }
func (l *Log) Error(args ...interface{})                 { l.internal.Error(args...) }
func (l *Log) Errorf(format string, args ...interface{}) { l.internal.Errorf(format, args...) }
func (l *Log) Errorln(args ...interface{})               { l.internal.Errorln(args...) }
func (l *Log) Info(args ...interface{})                  { l.internal.Info(args...) }
func (l *Log) Infof(format string, args ...interface{})  { l.internal.Infof(format, args...) }
func (l *Log) Infoln(args ...interface{})                { l.internal.Infoln(args...) }
func (l *Log) Warn(args ...interface{})                  { l.internal.Warn(args...) }
func (l *Log) Warnf(format string, args ...interface{})  { l.internal.Warnf(format, args...) }
func (l *Log) Warnln(args ...interface{})                { l.internal.Warnln(args...) }
