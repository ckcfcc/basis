package logx

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	ext "github.com/ckcfcc/basis/logx/lgext"
	"github.com/sirupsen/logrus"
)

var _glog = NewLog("GLOBAL")
var _logidx = uint32(0)
var _glogLv = logrus.InfoLevel

type Level uint32

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

type Logger interface {
	SetLevel(Level)
	LogErrorf(string, ...interface{}) error
	With(string, interface{}) Logger
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})
	Panic(...interface{})

	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
	Panicf(string, ...interface{})

	Dump([]byte, ...interface{})
	Dumpf([]byte, string, ...interface{})
}

// 包导出函数
func NewLog(name string) Logger {
	if name == "" {
		name = fmt.Sprintf("L_%d", atomic.AddUint32(&_logidx, 1))
	}
	result := &logger{name: name, lgr: logrus.New()}

	ntf := new(ext.TextFormatter)
	ntf.ForceFormatting = true
	ntf.TimestampFormat = "060102_150405.000"
	ntf.FullTimestamp = true
	ntf.SpacePadding = 160

	result.lgr.Formatter = ntf
	result.lgr.Level = _glogLv
	ssHook := ext.StdSSHook()
	ssHook.SetWrapperFile("ckcfcc/basis/logx/logx.")
	result.lgr.Hooks.Add(ssHook)
	return result
}

// TODO:暂时弃用 性能低下
func Bytes2Str(buff []byte, suffix int, padfix int) (hexStr string) {
	var head, tail string
	var bl int = len(buff)
	preSpace := strings.Repeat(" ", suffix)
	padSpace := strings.Repeat(" ", padfix)

	tb := []byte{}
	for i := 1; i <= bl; i++ {
		// 算出尾部 字符串部分
		c := buff[i-1]

		if c < 0x20 || c > 0x7e {
			c = '.'
		}

		tb = append(tb, c)

		if len(tb) >= 16 {
			tail, tb = string(tb), []byte{}
		}

		// 开头的行号
		if i%16 == 1 {
			head = fmt.Sprintf("%08X  ", i-1)
			hexStr += preSpace + head
		}

		// 中间的16进制输出部分
		if i%4 == 1 {
			hexStr += fmt.Sprintf("%s %02x", padSpace, buff[i-1])
		} else {
			hexStr += fmt.Sprintf(" %02x", buff[i-1])
		}

		if i%16 == 0 {
			hexStr += "   " + tail + "\n"
		} else if i == bl {
			// 最后一排补齐
			sub := 16 - bl%16
			if sub != 16 {
				for j := sub; j > 0; j-- {
					hexStr += "   "
				}
				// 补齐pad差值数
				for j := sub / 4; j > 0; j-- {
					hexStr += padSpace
				}
			}
			//fmt.Printf("最后一排补齐 sub:%d runes:%v tail:%s\n", sub, runes, tail)
			tail = string(tb)
			hexStr += "   " + tail + "\n"
		}
	}
	return
}

func SetLevel(lv Level) {
	_glogLv = logrus.Level(lv)
	_glog.(*logger).lgr.Level = logrus.Level(lv)
}

// output err string and return error
func LogErrorf(format string, args ...interface{}) (err error) {
	err = fmt.Errorf(format, args...)
	_glog.Error(err)
	return
}

// use global _glog
func Dump(buff []byte, msgs ...interface{}) {
	_glog.Dump(buff, msgs...)
}

func Dumpf(buff []byte, format string, msgs ...interface{}) {
	_glog.Dumpf(buff, format, msgs...)
}

func With(k string, v interface{}) Logger {
	return _glog.With(k, v)
}

func Debug(msgs ...interface{}) {
	_glog.Debug(msgs...)
}

func Info(msgs ...interface{}) {
	_glog.Info(msgs...)
}

func Warn(msgs ...interface{}) {
	_glog.Warn(msgs...)
}

func Error(msgs ...interface{}) {
	_glog.Error(msgs...)
}

func Fatal(msgs ...interface{}) {
	_glog.Fatal(msgs...)
}

func Panic(msgs ...interface{}) {
	_glog.Panic(msgs...)
}

func Debugf(format string, msgs ...interface{}) {
	_glog.Debugf(format, msgs...)
}

func Infof(format string, msgs ...interface{}) {
	_glog.Infof(format, msgs...)
}

func Warnf(format string, msgs ...interface{}) {
	_glog.Warnf(format, msgs...)
}

func Errorf(format string, msgs ...interface{}) {
	_glog.Errorf(format, msgs...)
}

func Fatalf(format string, msgs ...interface{}) {
	_glog.Fatalf(format, msgs...)
}

func Panicf(format string, msgs ...interface{}) {
	_glog.Panicf(format, msgs...)
}

type logger struct {
	name string
	lgr  *logrus.Logger
}

type entry logrus.Entry

// logger 对象相关函数
func (this *logger) SetLevel(lv Level) {
	this.lgr.Level = logrus.Level(lv)
}

func (this *logger) SetLogFile() {
	// TODO: 需要实现一个文件hook
}

func (this *logger) LogErrorf(format string, args ...interface{}) (err error) {
	err = fmt.Errorf(format, args...)
	this.Error(err)
	return
}

func (this *logger) Dump(buff []byte, msgs ...interface{}) {
	if this.lgr.Level != logrus.DebugLevel {
		return
	}

	bl := len(buff)
	//dumpStr := Bytes2Str(buff, 10, 1)
	dumpStr := hex.Dump(buff)
	this.With("dump", dumpStr).Debug("Dump Size:", strconv.Itoa(bl), "(0x", strconv.FormatInt(int64(bl), 16), ") ", msgs[:])
}

func (this *logger) Dumpf(buff []byte, format string, msgs ...interface{}) {
	if this.lgr.Level != logrus.DebugLevel {
		return
	}

	bl := len(buff)
	//dumpStr := Bytes2Str(buff, 10, 1)
	dumpStr := hex.Dump(buff)
	this.With("dump", dumpStr).Debugf("Dump Size:%d(0x%0x) %s", bl, bl, fmt.Sprintf(format, msgs...))
}

func (this *logger) With(k string, v interface{}) Logger {
	if this.name != "" {
		return (*entry)(this.lgr.WithField("prefix", this.name).WithField(k, v))
	} else {
		return (*entry)(this.lgr.WithField(k, v))
	}
}

func (this *logger) Debug(msgs ...interface{}) {
	if this.name != "" {
		this.lgr.WithField("prefix", this.name).Debug(msgs...)
	} else {
		this.lgr.Debug(msgs...)
	}
}

func (this *logger) Info(msgs ...interface{}) {
	if this.name != "" {
		this.lgr.WithField("prefix", this.name).Info(msgs...)
	} else {
		this.lgr.Info(msgs...)
	}
}

func (this *logger) Warn(msgs ...interface{}) {
	if this.name != "" {
		this.lgr.WithField("prefix", this.name).Warn(msgs...)
	} else {
		this.lgr.Warn(msgs...)
	}
}

func (this *logger) Error(msgs ...interface{}) {
	if this.name != "" {
		this.lgr.WithField("prefix", this.name).Error(msgs...)
	} else {
		this.lgr.Error(msgs...)
	}
}

func (this *logger) Fatal(msgs ...interface{}) {
	if this.name != "" {
		this.lgr.WithField("prefix", this.name).Fatal(msgs...)
	} else {
		this.lgr.Fatal(msgs...)
	}
}

func (this *logger) Panic(msgs ...interface{}) {
	if this.name != "" {
		this.lgr.WithField("prefix", this.name).Panic(msgs...)
	} else {
		this.lgr.Panic(msgs...)
	}
}

func (this *logger) Debugf(format string, msgs ...interface{}) {
	if this.name != "" {
		this.lgr.WithField("prefix", this.name).Debugf(format, msgs...)
	} else {
		this.lgr.Debugf(format, msgs...)
	}
}

func (this *logger) Infof(format string, msgs ...interface{}) {
	if this.name != "" {
		this.lgr.WithField("prefix", this.name).Infof(format, msgs...)
	} else {
		this.lgr.Infof(format, msgs...)
	}
}

func (this *logger) Warnf(format string, msgs ...interface{}) {
	if this.name != "" {
		this.lgr.WithField("prefix", this.name).Warnf(format, msgs...)
	} else {
		this.lgr.Warnf(format, msgs...)
	}
}

func (this *logger) Errorf(format string, msgs ...interface{}) {
	if this.name != "" {
		this.lgr.WithField("prefix", this.name).Errorf(format, msgs...)
	} else {
		this.lgr.Errorf(format, msgs...)
	}
}

func (this *logger) Fatalf(format string, msgs ...interface{}) {
	if this.name != "" {
		this.lgr.WithField("prefix", this.name).Fatalf(format, msgs...)
	} else {
		this.lgr.Fatalf(format, msgs...)
	}
}

func (this *logger) Panicf(format string, msgs ...interface{}) {
	if this.name != "" {
		this.lgr.WithField("prefix", this.name).Panicf(format, msgs...)
	} else {
		this.lgr.Panicf(format, msgs...)
	}
}

// entry
func (this *entry) SetLevel(lv Level) {
	this.Level = logrus.Level(lv)
}

func (this *entry) LogErrorf(format string, args ...interface{}) (err error) {
	err = fmt.Errorf(format, args...)
	this.Error(err)
	return
}

func (this *entry) Dump(buff []byte, msgs ...interface{}) {
	if this.Level != logrus.DebugLevel {
		return
	}
	bl := len(buff)
	//dumpStr := Bytes2Str(buff, 10, 1)
	dumpStr := hex.Dump(buff)
	this.With("dump", dumpStr).Debug("Dump Size:", strconv.Itoa(bl), "(0x", strconv.FormatInt(int64(bl), 16), ") ", msgs[:])
}

func (this *entry) Dumpf(buff []byte, format string, msgs ...interface{}) {
	if this.Level != logrus.DebugLevel {
		return
	}
	bl := len(buff)
	//dumpStr := Bytes2Str(buff, 10, 1)
	dumpStr := hex.Dump(buff)
	this.With("dump", dumpStr).Debugf("Dump Size:%d(0x%0x) %s", bl, bl, fmt.Sprintf(format, msgs...))
}

func (this *entry) With(k string, v interface{}) Logger {
	return (*entry)((*logrus.Entry)(this).WithField(k, v))
}

func (this *entry) Debug(msgs ...interface{}) {
	(*logrus.Entry)(this).Debug(msgs...)
}

func (this *entry) Info(msgs ...interface{}) {
	(*logrus.Entry)(this).Info(msgs...)
}

func (this *entry) Warn(msgs ...interface{}) {
	(*logrus.Entry)(this).Warn(msgs...)
}

func (this *entry) Error(msgs ...interface{}) {
	(*logrus.Entry)(this).Error(msgs...)
}

func (this *entry) Fatal(msgs ...interface{}) {
	(*logrus.Entry)(this).Fatal(msgs...)
}

func (this *entry) Panic(msgs ...interface{}) {
	(*logrus.Entry)(this).Panic(msgs...)
}

func (this *entry) Debugf(format string, msgs ...interface{}) {
	(*logrus.Entry)(this).Debugf(format, msgs...)
}

func (this *entry) Infof(format string, msgs ...interface{}) {
	(*logrus.Entry)(this).Infof(format, msgs...)
}

func (this *entry) Warnf(format string, msgs ...interface{}) {
	(*logrus.Entry)(this).Warnf(format, msgs...)
}

func (this *entry) Errorf(format string, msgs ...interface{}) {
	(*logrus.Entry)(this).Errorf(format, msgs...)
}

func (this *entry) Fatalf(format string, msgs ...interface{}) {
	(*logrus.Entry)(this).Fatalf(format, msgs...)
}

func (this *entry) Panicf(format string, msgs ...interface{}) {
	(*logrus.Entry)(this).Panicf(format, msgs...)
}
