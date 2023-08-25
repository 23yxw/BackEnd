package utils

import (
	"io"
	"strings"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var onceForLogger sync.Once
var ErrorLogger *zap.SugaredLogger

func NewLoggerHelper() {
	onceForLogger.Do(func() {
		Loginit()
	})
}

func Loginit() {
	// 设置一些基本日志格式 具体含义还比较好理解，直接看zap源码也不难懂
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})
	// 实现两个判断日志等级的interface
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	// 获取 info、error日志文件的io.Writer 抽象 getWriter() 在下方实现
	infoWriter := getWriter(ParamsInstance.InfoLogPath)
	errorWriter := getWriter(ParamsInstance.ErrorLogPath)
	// 最后创建具体的Logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel),
	)
	log := zap.New(core, zap.AddCaller()) // 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数, 有点小坑
	ErrorLogger = log.Sugar()
}
func getWriter(filename string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	// 保存7天内的日志，每1*24小时(整点)分割一次日志
	hook, err := rotatelogs.New(
		strings.Replace(filename, ".log", "", -1)+"-%Y%m%d%H.log", // 没有使用go风格反人类的format格式
		//rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		panic(err)
	}
	return hook
}
func Debug(args ...interface{}) {
	ErrorLogger.Debug(args...)
}
func Debugf(template string, args ...interface{}) {
	ErrorLogger.Debugf(template, args...)
}
func Info(args ...interface{}) {
	ErrorLogger.Info(args...)
}
func Infof(template string, args ...interface{}) {
	ErrorLogger.Infof(template, args...)
}
func Warn(args ...interface{}) {
	ErrorLogger.Warn(args...)
}
func Warnf(template string, args ...interface{}) {
	ErrorLogger.Warnf(template, args...)
}
func Error(args ...interface{}) {
	ErrorLogger.Error(args...)
}
func Errorf(template string, args ...interface{}) {
	ErrorLogger.Errorf(template, args...)
}
func DPanic(args ...interface{}) {
	ErrorLogger.DPanic(args...)
}
func DPanicf(template string, args ...interface{}) {
	ErrorLogger.DPanicf(template, args...)
}
func Panic(args ...interface{}) {
	ErrorLogger.Panic(args...)
}
func Panicf(template string, args ...interface{}) {
	ErrorLogger.Panicf(template, args...)
}
func Fatal(args ...interface{}) {
	ErrorLogger.Fatal(args...)
}
func Fatalf(template string, args ...interface{}) {
	ErrorLogger.Fatalf(template, args...)
}
