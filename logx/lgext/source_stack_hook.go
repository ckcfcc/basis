package lgext

import (
	"strings"

	"github.com/ckcfcc/basis/sysx/stack"
	"github.com/sirupsen/logrus"
)

const logrus_path = "github.com/sirupsen/logrus"
const defaultFrameNum = 4

type LogrusStackHook struct {
	CallerLevels []logrus.Level
	StackLevels  []logrus.Level
	FrameNum     int
	Wrapper      string
}

func (hook *LogrusStackHook) SetWrapperFile(file string) {
	hook.Wrapper = file
}

func NewSSHook(callerLevels []logrus.Level, stackLevels []logrus.Level) *LogrusStackHook {
	return &LogrusStackHook{
		CallerLevels: callerLevels,
		StackLevels:  stackLevels,
	}
}

func StdSSHook() *LogrusStackHook {
	return &LogrusStackHook{
		CallerLevels: logrus.AllLevels,
		StackLevels:  []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel},
	}
}

func (hook *LogrusStackHook) Levels() []logrus.Level {
	return hook.CallerLevels
}

func (hook *LogrusStackHook) Fire(entry *logrus.Entry) error {
	skipFrames := defaultFrameNum

	// if len(entry.Data) == 0 {
	// 	skipFrames = hook.FrameNum // 不带f的函数用7 例如 Info()
	// 	//skipFrames = hook.FrameNum + 2 // 带f的函数用8 例如 Infof()
	// } else {
	// 	skipFrames = hook.FrameNum
	// }

	var frames stack.Stack

	_frames := stack.Callers(skipFrames)

	for _, frame := range _frames {
		if !strings.Contains(frame.File, logrus_path) &&
			!strings.Contains(frame.File, hook.Wrapper) {
			frames = append(frames, frame)
		}
	}

	if len(frames) > 0 {
		for _, level := range hook.CallerLevels {
			if entry.Level == level {
				entry.Data["caller"] = frames[0]
				break
			}
		}

		for _, level := range hook.StackLevels {
			if entry.Level == level {
				entry.Data["stack"] = stack.Callers(-2)
				break
			}
		}
	}
	return nil
}
