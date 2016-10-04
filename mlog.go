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
	Logger	*log.Logger
	Prefix	string
	Writer	io.Writer
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

	TRACE			*log.Logger
	DEBUG			*log.Logger
	INFO			*log.Logger
	WARN			*log.Logger
	ERROR			*log.Logger
	CRITICAL	*log.Logger
	FATAL			*log.Logger

	loggers map[string]*mlogger = map[string]*mlogger {
		"trace":		&mlogger{ Logger:TRACE, Prefix:"TRACE", Writer:os.Stdout },
		"debug":		&mlogger{ Logger:DEBUG, Prefix:"DEBUG", Writer:os.Stdout },
		"info":			&mlogger{ Logger:INFO, Prefix:"INFO", Writer:os.Stdout },
		"warn":			&mlogger{ Logger:WARN, Prefix:"WARN", Writer:os.Stdout },
		"error":		&mlogger{ Logger:ERROR, Prefix:"ERROR", Writer:os.Stdout },
		"critical":	&mlogger{ Logger:CRITICAL, Prefix:"CRITICAL", Writer:os.Stdout },
		"fatal":		&mlogger{ Logger:FATAL, Prefix:"FATAL", Writer:os.Stdout },
	}
)


func init() {
	levelsEnum = make(map[string]int)

	// Enumerate the core logging levels so that thresholding can be applied
	for index, name := range []string{"trace", "debug", "info", "warn", "error", "critical", "fatal"} {
		levelsEnum[name] = index
	}

	// Initialize loggers
	for _, l := range loggers {
		l.Logger = log.New(os.Stdout, l.Prefix, flags)
	}

	SetThreshold(DEFAULT_THRESHOLD)
}


func applyThreshold() {
	t := levelsEnum[threshold]
	for key, level := range levelsEnum {
		l := loggers[key]
		if level < t {
			// Apply discard writer
			l.Logger.SetOutput(DISCARD)
		} else {
			// Restore configured writer
			l.Logger.SetOutput(l.Writer)
		}
	}
}


func Threshold() string {
	return threshold
}


func NewCustomLogger(name string, prefix string) {
	loggers[name] = &mlogger{
		Logger:		log.New(os.Stdout, prefix, flags),
		Prefix:		prefix,
		Writer:		os.Stdout,
	}
}


func Println(logger string, v ...interface{}) {
	if handler, ok := loggers[logger]; ok == true {
		handler.Logger.Println(v...)
	}
	return
}


func Printf(logger string, format string, v ...interface{}) {
	if handler, ok := loggers[logger]; ok == true {
		handler.Logger.Printf(format, v...)
	}
	return
}


// Set the log flags for all loggers (available: DATE, TIME, SFILE, LFILE, and MSEC).
func SetFlags(flagset int, loggerList ...string) {
	flags = flagset
	if len(loggerList) == 0 {
		// Change flags for ALL loggers
		for _, ml := range loggers {
			ml.Logger.SetFlags(flags)
		}
	} else {
		// Change flags for named loggers
		for _, logger := range loggerList {
			if ml, ok := loggers[logger]; ok == true {
				ml.Logger.SetFlags(flags)
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

	switch len(writers) {
	case 0:
		WARN.Println("no io.Writer(s) provided for output", logger)
		return
	case 1:
		handler.Writer = writers[0]
	default:
		handler.Writer = io.MultiWriter(writers...)
	}

	handler.Logger.SetOutput(handler.Writer)

	if _, ok := loggers[logger]; ok == true {
		applyThreshold()
	}
}


func SetThreshold(level string) {
	if _, ok := levelsEnum[level]; ok == false {
		WARN.Printf("ignoring invalid log level '%s'", level)
		return
	}
	threshold = level
	applyThreshold()
}
