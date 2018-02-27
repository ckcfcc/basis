package lgext

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/mgutz/ansi"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

const defaultTimestampFormat = time.RFC3339

var (
	baseTimestamp      time.Time    = time.Now()
	defaultColorScheme *ColorScheme = &ColorScheme{
		InfoLevelStyle:  "green",
		WarnLevelStyle:  "yellow",
		ErrorLevelStyle: "red",
		FatalLevelStyle: "red",
		PanicLevelStyle: "red",
		DebugLevelStyle: "blue",
		PrefixStyle:     "cyan",
		TimestampStyle:  "black+h",
	}
	noColorsColorScheme *compiledColorScheme = &compiledColorScheme{
		InfoLevelColor:  ansi.ColorFunc(""),
		WarnLevelColor:  ansi.ColorFunc(""),
		ErrorLevelColor: ansi.ColorFunc(""),
		FatalLevelColor: ansi.ColorFunc(""),
		PanicLevelColor: ansi.ColorFunc(""),
		DebugLevelColor: ansi.ColorFunc(""),
		PrefixColor:     ansi.ColorFunc(""),
		TimestampColor:  ansi.ColorFunc(""),
	}
	defaultCompiledColorScheme *compiledColorScheme = compileColorScheme(defaultColorScheme)
)

func miniTS() int {
	return int(time.Since(baseTimestamp) / time.Second)
}

type ColorScheme struct {
	InfoLevelStyle  string
	WarnLevelStyle  string
	ErrorLevelStyle string
	FatalLevelStyle string
	PanicLevelStyle string
	DebugLevelStyle string
	PrefixStyle     string
	TimestampStyle  string
}

type compiledColorScheme struct {
	InfoLevelColor  func(string) string
	WarnLevelColor  func(string) string
	ErrorLevelColor func(string) string
	FatalLevelColor func(string) string
	PanicLevelColor func(string) string
	DebugLevelColor func(string) string
	PrefixColor     func(string) string
	TimestampColor  func(string) string
}

type TextFormatter struct {
	ForceColors bool

	DisableColors bool

	ForceFormatting bool

	DisableTimestamp bool

	DisableUppercase bool

	FullTimestamp bool

	TimestampFormat string

	DisableSorting bool

	QuoteEmptyFields bool

	QuoteCharacter string

	SpacePadding int

	colorScheme *compiledColorScheme

	isTerminal bool

	sync.Once
}

func getCompiledColor(main string, fallback string) func(string) string {
	var style string
	if main != "" {
		style = main
	} else {
		style = fallback
	}
	return ansi.ColorFunc(style)
}

func compileColorScheme(s *ColorScheme) *compiledColorScheme {
	return &compiledColorScheme{
		InfoLevelColor:  getCompiledColor(s.InfoLevelStyle, defaultColorScheme.InfoLevelStyle),
		WarnLevelColor:  getCompiledColor(s.WarnLevelStyle, defaultColorScheme.WarnLevelStyle),
		ErrorLevelColor: getCompiledColor(s.ErrorLevelStyle, defaultColorScheme.ErrorLevelStyle),
		FatalLevelColor: getCompiledColor(s.FatalLevelStyle, defaultColorScheme.FatalLevelStyle),
		PanicLevelColor: getCompiledColor(s.PanicLevelStyle, defaultColorScheme.PanicLevelStyle),
		DebugLevelColor: getCompiledColor(s.DebugLevelStyle, defaultColorScheme.DebugLevelStyle),
		PrefixColor:     getCompiledColor(s.PrefixStyle, defaultColorScheme.PrefixStyle),
		TimestampColor:  getCompiledColor(s.TimestampStyle, defaultColorScheme.TimestampStyle),
	}
}

func (f *TextFormatter) init(entry *logrus.Entry) {
	if len(f.QuoteCharacter) == 0 {
		f.QuoteCharacter = "\""
	}
	if entry.Logger != nil {
		f.isTerminal = f.checkIfTerminal(entry.Logger.Out)
	}
}

func (f *TextFormatter) checkIfTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}

func (f *TextFormatter) SetColorScheme(colorScheme *ColorScheme) {
	f.colorScheme = compileColorScheme(colorScheme)
}

func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	var keys []string = make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}
	lastKeyIdx := len(keys) - 1

	if !f.DisableSorting {
		sort.Strings(keys)
	}
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	prefixFieldClashes(entry.Data)

	f.Do(func() { f.init(entry) })

	isFormatted := f.ForceFormatting || f.isTerminal

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}
	if isFormatted {
		isColored := (f.ForceColors || f.isTerminal) && !f.DisableColors
		var colorScheme *compiledColorScheme
		if isColored {
			if f.colorScheme == nil {
				colorScheme = defaultCompiledColorScheme
			} else {
				colorScheme = f.colorScheme
			}
		} else {
			colorScheme = noColorsColorScheme
		}
		f.printColored(b, entry, keys, timestampFormat, colorScheme)
	} else {
		if !f.DisableTimestamp {
			f.appendKeyValue(b, "time", entry.Time.Format(timestampFormat), true)
		}
		f.appendKeyValue(b, "level", entry.Level.String(), true)
		if entry.Message != "" {
			f.appendKeyValue(b, "msg", entry.Message, lastKeyIdx >= 0)
		}
		// b.WriteByte('\t')
		// b.WriteByte('\n')
		for i, key := range keys {
			f.appendKeyValue(b, key, entry.Data[key], lastKeyIdx != i)
		}
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *TextFormatter) printColored(b *bytes.Buffer, entry *logrus.Entry, keys []string, timestampFormat string, colorScheme *compiledColorScheme) {
	var levelColor func(string) string
	var levelText string
	switch entry.Level {
	case logrus.InfoLevel:
		levelColor = colorScheme.InfoLevelColor
	case logrus.WarnLevel:
		levelColor = colorScheme.WarnLevelColor
	case logrus.ErrorLevel:
		levelColor = colorScheme.ErrorLevelColor
	case logrus.FatalLevel:
		levelColor = colorScheme.FatalLevelColor
	case logrus.PanicLevel:
		levelColor = colorScheme.PanicLevelColor
	default:
		levelColor = colorScheme.DebugLevelColor
	}

	if entry.Level != logrus.WarnLevel {
		levelText = entry.Level.String()
	} else {
		levelText = "warn"
	}

	if !f.DisableUppercase {
		levelText = strings.ToUpper(levelText)
	}

	level := levelColor(fmt.Sprintf("%5s", levelText))
	prefix := ""
	message := entry.Message

	if prefixValue, ok := entry.Data["prefix"]; ok {
		prefix = colorScheme.PrefixColor(" " + prefixValue.(string) + ":")
	} else {
		prefixValue, trimmedMsg := extractPrefix(entry.Message)
		if len(prefixValue) > 0 {
			prefix = colorScheme.PrefixColor(" " + prefixValue + ":")
			message = trimmedMsg
		}
	}

	messageFormat := "%s"
	if f.SpacePadding != 0 {
		messageFormat = fmt.Sprintf("%%-%ds", f.SpacePadding)
	}

	if f.DisableTimestamp {
		fmt.Fprintf(b, "%s%s "+messageFormat, level, prefix, message)
	} else {
		var timestamp string
		if !f.FullTimestamp {
			timestamp = fmt.Sprintf("[%04d]", miniTS())
		} else {
			timestamp = fmt.Sprintf("[%s]", entry.Time.Format(timestampFormat))
		}
		fmt.Fprintf(b, "%s %s%s "+messageFormat, colorScheme.TimestampColor(timestamp), level, prefix, message)
	}
	for _, k := range keys {
		if k != "prefix" && k != "stack" && k != "dump" && k != "caller" {
			v := entry.Data[k]
			fmt.Fprintf(b, " %s=%+v", levelColor(k), v)
		}
		if k == "caller" {
			v := entry.Data[k]
			fmt.Fprintf(b, " \n\t%s", colorScheme.TimestampColor(fmt.Sprintf("%s=%s", k, v)))
		}
		if k == "stack" {
			v := entry.Data[k]
			fmt.Fprintf(b, "\n\t\t%s:%+v", levelColor("当前调用栈"), v)
		}
		if k == "dump" {
			v := entry.Data[k]
			fmt.Fprintf(b, "\n\t\t%s:\n%+v", levelColor("dump内容"), v)
		}
	}
}

func (f *TextFormatter) needsQuoting(text string) bool {
	if f.QuoteEmptyFields && len(text) == 0 {
		return true
	}
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.') {
			return true
		}
	}
	return false
}

func extractPrefix(msg string) (string, string) {
	prefix := ""
	regex := regexp.MustCompile("^\\[(.*?)\\]")
	if regex.MatchString(msg) {
		match := regex.FindString(msg)
		prefix, msg = match[1:len(match)-1], strings.TrimSpace(msg[len(match):])
	}
	return prefix, msg
}

func (f *TextFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}, appendSpace bool) {
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)

	if appendSpace {
		b.WriteByte(' ')
	}
}

func (f *TextFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	switch value := value.(type) {
	case string:
		if !f.needsQuoting(value) {
			b.WriteString(value)
		} else {
			fmt.Fprintf(b, "%s%v%s", f.QuoteCharacter, value, f.QuoteCharacter)
		}
	case error:
		errmsg := value.Error()
		if !f.needsQuoting(errmsg) {
			b.WriteString(errmsg)
		} else {
			fmt.Fprintf(b, "%s%v%s", f.QuoteCharacter, errmsg, f.QuoteCharacter)
		}
	default:
		fmt.Fprint(b, value)
	}
}

func prefixFieldClashes(data logrus.Fields) {
	if t, ok := data["time"]; ok {
		data["fields.time"] = t
	}

	if m, ok := data["msg"]; ok {
		data["fields.msg"] = m
	}

	if l, ok := data["level"]; ok {
		data["fields.level"] = l
	}
}
