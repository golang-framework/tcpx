// Copyright 2022 YuWenYu  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package z

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	zoErr *zap.SugaredLogger = nil
	z9210 *zap.SugaredLogger = nil

	qoErr = make(chan *soErr, 1024)
	q9210 = make(chan *s9210, 2048)
)

type (
	soErr struct {
		err string
	}

	s9210 struct {
		err string
	}
)

func init() {
	go startOErrInQueue()
	go start9210InQueue()
}

func SetOErrToZ(err string) {
	qoErr <- &soErr{err: err}
}

func startOErrInQueue() {
	for {
		select {
		case d, ok := <-qoErr:
			if ok == false {
				break
			}

			xoErr().Info(d.err)
			break
		}
	}
}

func Set9210ToZ(err string) {
	q9210 <- &s9210{err: err}
}

func start9210InQueue() {
	for {
		select {
		case d, ok := <-q9210:
			if ok == false {
				break
			}

			x9210().Info(d.err)
			break
		}
	}
}

func xoErr() *zap.SugaredLogger {
	zoErr = z("err")
	return zoErr
}

func x9210() *zap.SugaredLogger {
	z9210 = z("9210")
	return z9210
}

func z(mode string) *zap.SugaredLogger {
	var z *zap.SugaredLogger = nil
	var zWriteSyncerFile = mode

	switch mode {
	case "err":
		z = zoErr
		break

	case "9210":
		z = z9210
		break

	default:
		return nil
	}

	if z == nil {
		return initialized(zWriteSyncerFile)
	}

	return z
}

func initialized(zWriteSyncerFile string) *zap.SugaredLogger {
	var mode = strings.Split(zWriteSyncerFile, "-")
	if len(mode) < 1 {
		return nil
	}

	var z *zap.SugaredLogger = nil
	var c zapcore.Core = nil

	c = zapcore.NewCore(encode(), zapTmpWriteSyncer(zWriteSyncerFile), zapcore.InfoLevel)
	z = zap.New(c).Sugar()

	return z
}

func encode() zapcore.Encoder {
	conf := zap.NewProductionEncoderConfig()
	conf.EncodeTime = zapcore.ISO8601TimeEncoder // 修改时间编码器

	// 在日志文件中使用大写字母记录日志级别
	conf.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(conf)
}

func zapTmpWriteSyncer(name string) zapcore.WriteSyncer {
	var buf strings.Builder
	buf.WriteString("./logs/ori/")

	if _, err := os.Stat(buf.String()); os.IsNotExist(err) {
		_ = os.MkdirAll(buf.String(), os.ModePerm)
	}

	buf.WriteString(name)
	buf.WriteString(".log")

	fn := buf.String()
	buf.Reset()

	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   fn,
		MaxSize:    1,
		MaxBackups: 1,
		Compress:   true,
		LocalTime:  true,
	})
}
