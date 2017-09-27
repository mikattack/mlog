/*
 * Copyright © 2016 Alex Mikitik.
 *
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * NOTE:  This is a derived work. Original license may
 *				be found in the LICENSE-GO file.
 */

package mlog

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

const default_call_level = 2

type LogLevel string

// These flags define additional fields to include in each log entry.
const (
	// Bits or'ed together to control which fields are included in log entries.
	// There is no control over the field value's format.
	DATE  = 1 << iota // Date and time in UTC: 2009/01/23 01:23:45
	FILE              // File name element and line number: example.go:36
	LEVEL             // Log level: [INFO]
)

// Logging levels
const (
	IN_TESTING     LogLevel = "DEBUG"
	IN_PRODUCTION  LogLevel = "INFO"
	TO_INVESTIGATE LogLevel = "WARNING"
	PAGE_ME_NOW    LogLevel = "ERROR"
)

var (
	// Fallback ordering of default levels defined by the stock logger.
	defaultLevels []LogLevel

	// Level-to-number mapping of the configured logging levels.
	levelOrder map[LogLevel]int
)

// Set up log levels and ordering
func init() {
	defaultLevels = []LogLevel{
		IN_TESTING,
		IN_PRODUCTION,
		TO_INVESTIGATE,
		PAGE_ME_NOW,
	}

	levelOrder = map[LogLevel]int{}
	for index, level := range defaultLevels {
		levelOrder[level] = index
	}
}

type Logger struct {
	buffer []byte     // Reusable empty `entry`
	level  LogLevel   // Logging level threshold (default: InProduction)
	mu     sync.Mutex // Ensures atomic writes
	out    io.Writer  // Output destination
	flag   int        // Output flags indicating extra logging information
}

// Creates a new logger. The writer argument sets the destination to which
// log data will be written. The flag argument defines which fields are
// included into log messages by default.
func New(writer io.Writer, flag int) *Logger {
	logger := &Logger{
		level: IN_PRODUCTION,
		out:   writer,
		flag:  flag,
	}
	return logger
}

var std = New(os.Stdout, LEVEL|DATE|FILE) // Default logger

// Returns the flag(s) for the logger's default fields.
func (l *Logger) Flags() int {
	return l.flag
}

// Returns the logging threshold level.
func (l *Logger) Threshold() LogLevel {
	return l.level
}

// Sets which fields are included into log messages by default.
func (l *Logger) SetFlags(flag int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.flag = flag
}

// Sets the output destination for the logger.
func (l *Logger) SetOutput(writers ...io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()

	switch len(writers) {
	case 0:
		InProduction("SetOutput: no io.Writer(s) provided for output")
		return
	case 1:
		l.out = writers[0]
	default:
		l.out = io.MultiWriter(writers...)
	}
}

// Sets the level the logger should emit messages at.
func (l *Logger) SetThreshold(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := levelOrder[level]; ok == false {
		return // Ignore invalid log levels
	}
	l.level = level
}

// Logs a message "in testing". This level is intended for noisy and
// verbose output, often during development.
func (l *Logger) InTesting(format string, args ...interface{}) {
	l.log(IN_TESTING, format, args)
}

// Logs a message "in production". This level is intended for information
// needed to debug production issues.
func (l *Logger) InProduction(format string, args ...interface{}) {
	l.log(IN_PRODUCTION, format, args)
}

// Logs a message "to investigate later". This level is intended for
// important events which require special, but not immediate attention.
func (l *Logger) ToInvestigate(format string, args ...interface{}) {
	l.log(TO_INVESTIGATE, format, args)
}

// Logs a message of such importance, it should wake someone up in the
// the middle of the night. This level is intended for events which require
// immediate attention.
func (l *Logger) PageMeNow(format string, args ...interface{}) {
	l.log(PAGE_ME_NOW, format, args)
}

// Evaluates whether a given logging level meets the threshold currently
// set by the logger. If a given logging level is unknown (because of
// custom logging levels), the threshold check will always fail.
func (l *Logger) meetsThreshold(level LogLevel) bool {
	threshold, tOk := levelOrder[l.level]
	requestLevel, rOk := levelOrder[level]
	if tOk == false || rOk == false {
		fmt.Fprintf(l.out, "Invalid log level in comparison: %s >= %s", level, l.level)
		return false
	}
	return requestLevel >= threshold
}

// Fixed-width decimal ASCII. Negative width skips zero-padding.
func itoa(buffer *[]byte, i int, width int) {
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || width > 1 {
		width--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	b[bp] = byte('0' + i)
	*buffer = append(*buffer, b[bp:]...)
}

// Add level, date, and file information to a given buffer, if enabled.
func (l *Logger) formatHeader(buffer *[]byte, level LogLevel, t time.Time, file string, line int) {
	if l.flag&LEVEL != 0 {
		*buffer = append(*buffer, '[')
		*buffer = append(*buffer, string(level)...)
		*buffer = append(*buffer, "] "...)
	}
	if l.flag&DATE != 0 {
		t = t.UTC()
		year, month, day := t.Date()
		itoa(buffer, year, 4)
		*buffer = append(*buffer, '/')
		itoa(buffer, int(month), 2)
		*buffer = append(*buffer, '/')
		itoa(buffer, day, 2)
		*buffer = append(*buffer, ' ')

		hour, minute, second := t.Clock()
		itoa(buffer, hour, 2)
		*buffer = append(*buffer, ':')
		itoa(buffer, minute, 2)
		*buffer = append(*buffer, ':')
		itoa(buffer, second, 2)
		*buffer = append(*buffer, '.')
		itoa(buffer, t.Nanosecond()/1e3, 6)
		*buffer = append(*buffer, ' ')
	}
	if l.flag&FILE != 0 {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		*buffer = append(*buffer, short...)
		*buffer = append(*buffer, ':')
		itoa(buffer, line, -1)
		*buffer = append(*buffer, ": "...)
	}
}

// Logs a message at a given level
func (l *Logger) log(level LogLevel, format string, args ...interface{}) error {
	if l.meetsThreshold(level) == false {
		return nil
	}

	var now time.Time
	if l.flag&DATE != 0 {
		now = time.Now()
	}

	var file string
	var line int

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.flag&FILE != 0 {
		l.mu.Unlock()
		var okay bool
		_, file, line, okay = runtime.Caller(default_call_level)
		if !okay {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}

	l.buffer = l.buffer[:0] // Empty buffer
	l.formatHeader(&l.buffer, level, now, file, line)
	if len(args) > 1 {
		l.buffer = append(l.buffer, fmt.Sprintf(format, args...)...)
	} else {
		l.buffer = append(l.buffer, format...)
	}

	if len(format) == 0 || format[len(format)-1] != '\n' {
		l.buffer = append(l.buffer, '\n')
	}

	_, err := l.out.Write(l.buffer)
	return err
}

// Returns the flag(s) for the global logger's default fields.
func Flags() int {
	return std.Flags()
}

// Returns the global logging threshold level.
func Threshold() LogLevel {
	return std.Threshold()
}

// Sets which fields are included by default in the global logger.
func SetFlags(flag int) {
	std.SetFlags(flag)
}

// Sets the output destination for global logger.
func SetOutput(writer io.Writer) {
	std.SetOutput(writer)
}

// Sets the level the global logger should emit messages at.
func SetThreshold(level LogLevel) {
	std.SetThreshold(level)
}

// Logs a message "in testing" for the global logger.
func InTesting(format string, args ...interface{}) {
	std.InTesting(format, args...)
}

// Logs a message "in production" for the global logger.
func InProduction(format string, args ...interface{}) {
	std.InProduction(format, args...)
}

// Logs a message "to investigate later" for the global logger.
func ToInvestigate(format string, args ...interface{}) {
	std.ToInvestigate(format, args...)
}

// Logs a message of the highest importance to the global logger.
func PageMeNow(format string, args ...interface{}) {
	std.PageMeNow(format, args...)
}
