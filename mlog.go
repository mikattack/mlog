/*
 * Copyright © 2016 Alex Mikitik.
 *
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package mlog

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)


type mlogger struct {
	logger	*log.Logger
	prefix	string
	writer	io.Writer
}

func (ml *mlogger) Println(v ...interface{}) {
	ml.logger.Println(v...)
	return
}

func (ml *mlogger) Printf(format string, v ...interface{}) {
	ml.logger.Printf(format, v...)
	return
}


const (
	LEVEL_TRACE				= "trace"
	LEVEL_DEBUG				= "debug"
	LEVEL_INFO				= "info"
	LEVEL_WARN				= "warn"
	LEVEL_ERROR				= "error"
	LEVEL_CRITICAL		= "critical"
	LEVEL_FATAL				= "fatal"
	DEFAULT_THRESHOLD	= LEVEL_WARN

	NONE	= 0
	DATE	= log.Ldate
	TIME	= log.Ltime
	SFILE	= log.Lshortfile
	LFILE	= log.Llongfile
	MSEC	= log.Lmicroseconds
)

var (
	flags = DATE | TIME | SFILE
	threshold	string = DEFAULT_THRESHOLD
	levelsEnum map[string]int

	DISCARD		io.Writer = ioutil.Discard

	TRACE			*mlogger = &mlogger{ prefix:"TRACE: ", writer:os.Stdout }
	DEBUG			*mlogger = &mlogger{ prefix:"DEBUG: ", writer:os.Stdout }
	INFO			*mlogger = &mlogger{ prefix:"INFO: ", writer:os.Stdout }
	WARN			*mlogger = &mlogger{ prefix:"WARN: ", writer:os.Stdout }
	ERROR			*mlogger = &mlogger{ prefix:"ERROR: ", writer:os.Stdout }
	CRITICAL	*mlogger = &mlogger{ prefix:"CRITICAL: ", writer:os.Stdout }
	FATAL			*mlogger = &mlogger{ prefix:"FATAL: ", writer:os.Stdout }

	loggers map[string]*mlogger = map[string]*mlogger {
		"trace":		TRACE,
		"debug":		DEBUG,
		"info":		INFO,
		"warn":		WARN,
		"error":		ERROR,
		"critical":	CRITICAL,
		"fatal":		FATAL,
	}
)


func init() {
	levelsEnum = make(map[string]int)

	// Enumerate the core logging levels so that thresholding can be applied
	for index, name := range []string{"trace", "debug", "info", "warn", "error", "critical", "fatal"} {
		levelsEnum[name] = index
	}

	// Initialize default loggers using the values encoded in 'loggers' map
	for _, l := range loggers {
		l.logger = log.New(l.writer, l.prefix, flags)
	}

	SetThreshold(DEFAULT_THRESHOLD)
}


func Threshold() string {
	return threshold
}


func NewLogger(name string, prefix string) {
	loggers[name] = &mlogger{
		logger:		log.New(os.Stdout, prefix, flags),
		prefix:		prefix,
		writer:		os.Stdout,
	}
}


func Println(logger string, v ...interface{}) {
	if handler, ok := loggers[logger]; ok == true {
		handler.logger.Println(v...)
	}
	return
}


func Printf(logger string, format string, v ...interface{}) {
	if handler, ok := loggers[logger]; ok == true {
		handler.logger.Printf(format, v...)
	}
	return
}


// Set the log flags for loggers (available: DATE, TIME, SFILE, LFILE, MSEC,
// and NONE).  If no list of loggers is provided, then all logger's flags
// are set.
func SetFlags(flagset int, loggerList ...string) {
	flags = flagset
	if len(loggerList) == 0 {
		// Change flags for ALL loggers
		for _, ml := range loggers {
			ml.logger.SetFlags(flags)
		}
	} else {
		// Change flags for named loggers
		for _, logger := range loggerList {
			if ml, ok := loggers[logger]; ok == true {
				ml.logger.SetFlags(flags)
			}
		}
	}
}


func SetOutput(logger string, writers ...io.Writer) {
	handler, ok := loggers[logger]
	if ok == false {
		WARN.Printf("cannot set output of unknown logger '%s'", logger)
		return
	}

	// Set the logger's writer (ignoring thresholding)
	switch len(writers) {
	case 0:
		WARN.Println("no io.Writer(s) provided for output", logger)
		return
	case 1:
		handler.writer = writers[0]
	default:
		handler.writer = io.MultiWriter(writers...)
	}

	// Update the output handler of the logger
	if _, ok := levelsEnum[logger]; ok == true {
		/* 
		 * If the logger is subject to thresholding, we only update the handler
		 * when it falls within the current threshold.
		 */
		if levelsEnum[logger] >= levelsEnum[threshold] {
			handler.logger.SetOutput(handler.writer)
		}
	} else {
		// Loggers not subject to thresholding should be immediately updated
		handler.logger.SetOutput(handler.writer)
	}
}


func SetThreshold(level string) {
	// Update overall level
	if _, ok := levelsEnum[level]; ok == false {
		WARN.Printf("ignoring invalid log level '%s'", level)
		return
	}
	threshold = level

	// Re-evaluate each default logger's threshold
	enum := levelsEnum[threshold]
	for key, l := range levelsEnum {
		logger := *loggers[key].logger
		if l < enum {
			// Apply discard writer
			logger.SetOutput(DISCARD)
		} else {
			// Restore configured writer
			logger.SetOutput(loggers[key].writer)
		}
	}
}
