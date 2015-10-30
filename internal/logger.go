// Copyright 2015 Ka-Hing Cheung
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"fmt"
	glog "log"
	"os"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
)

var mu sync.Mutex
var loggers = make(map[string]*logHandle)

var log = GetLogger("main")
var fuseLog = GetLogger("fuse")

type logHandle struct {
	logrus.Logger

	name string
	lvl  *logrus.Level
}

func (l *logHandle) Format(e *logrus.Entry) ([]byte, error) {
	// Mon Jan 2 15:04:05 -0700 MST 2006
	const timeFormat = "2006/01/02 15:04:05.000000"
	lvl := e.Level
	if l.lvl != nil {
		lvl = *l.lvl
	}
	str := fmt.Sprintf("%v %v.%v %v",
		e.Time.Format(timeFormat),
		l.name,
		strings.ToUpper(lvl.String()),
		e.Message)

	if len(e.Data) != 0 {
		str += " " + fmt.Sprint(e.Data)
	}

	str += "\n"
	return []byte(str), nil
}

// for aws.Logger
func (l *logHandle) Log(args ...interface{}) {
	l.Debugln(args...)
}

func NewLogger(name string) *logHandle {
	l := &logHandle{name: name}
	l.Out = os.Stderr
	l.Formatter = l
	l.Level = logrus.InfoLevel
	return l
}

func GetLogger(name string) *logHandle {
	mu.Lock()
	defer mu.Unlock()

	if logger, ok := loggers[name]; ok {
		return logger
	} else {
		logger := NewLogger(name)
		loggers[name] = logger
		return logger
	}
}

func GetStdLogger(l *logHandle, lvl logrus.Level) *glog.Logger {
	mu.Lock()
	defer mu.Unlock()

	w := l.Writer()
	l.Formatter.(*logHandle).lvl = &lvl
	l.Level = lvl
	return glog.New(w, "", 0)
}
