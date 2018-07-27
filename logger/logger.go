package logger

import (
	"errors"
	"fmt"
	"runtime"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var logChan chan *logInfo

const bufLen int = 1000
const waitTime time.Duration = 5

type logInfo struct {
	level   string
	caller  string
	content string
}

func Init(proDir string) error {
	if proDir == "" {
		return errors.New("Invalid proDir")
	}

	var err error
	logLevel := zap.NewAtomicLevelAt(zap.DebugLevel)
	logEncoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "name",
		CallerKey:      "caller",
		StacktraceKey:  "",
		LineEnding:     zapcore.DefaultLineEnding, //add "\n" in line end automatically
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeTime:     logEncodeTime,   //show style of time
		EncodeCaller:   logEncodeCaller, //caller info
		EncodeName:     zapcore.FullNameEncoder,
	}
	logConfig := zap.Config{
		Level:            logLevel,
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    logEncoderConfig,
		OutputPaths:      []string{"stdout", proDir + "/log/server.log"},
		ErrorOutputPaths: []string{"stderr", proDir + "/log/logger_err.log"},
	}
	logger, err = logConfig.Build()
	if err != nil {
		return err
	}

	logChan = make(chan *logInfo, bufLen)
	go insertLog()

	return nil
}

func logEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func logEncodeCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	pc, file, line, ok := runtime.Caller(6)
	caller.PC = pc
	caller.File = file
	caller.Line = line
	caller.Defined = ok
	enc.AppendString(caller.String())
}

func newLogInfo(level string, description string) *logInfo {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return &logInfo{level, "undefined", description}
	}
	return &logInfo{level, fmt.Sprintf("%s:%d", file, line), description}
}

func insertLog() {
	for {
		logNum := len(logChan)

		if logNum == 0 {
			time.Sleep(waitTime * time.Second)
			continue
		}

		for i := 0; i < logNum; i++ {
			_, ok := <-logChan
			if !ok {
				continue
			}
		}

		runtime.Gosched()
		if len(logChan) < bufLen/2 {
			time.Sleep(waitTime * time.Second)
		}
	}
}

//detail is key and value, it is not necessary.
//use example:logger.Debug("test",zap.String("url","www.heqingfeng.com"))
func Debug(description ...interface{}) {
	des := fmt.Sprint(description...)
	logger.Debug(des)
}

func Info(description ...interface{}) {
	des := fmt.Sprint(description...)
	logger.Info(des)
	logChan <- newLogInfo("info", des)
}

func Warning(description ...interface{}) {
	des := fmt.Sprint(description...)
	logger.Warn(des)
	logChan <- newLogInfo("warn", des)
}

func Error(description ...interface{}) {
	des := fmt.Sprint(description...)
	logger.Error(des)
	logChan <- newLogInfo("error", des)
}

func Fatal(description ...interface{}) {
	des := fmt.Sprint(description...)
	logger.Fatal(des)
	logChan <- newLogInfo("fatal", des)
}

func Sync() {
	logger.Sync()
}
