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
	logger *log.Logger
	prefix string
	writer io.Writer
}

const (
	LEVEL_TEST       			= "debug"
	LEVEL_PRODUCTION 			= "info"
	LEVEL_TOMORROW   			= "warn"
	LEVEL_MIDDLE_OF_NIGHT = "error"

	DEFAULT_THRESHOLD = LEVEL_PRODUCTION

	NONE   = 0
	DATE   = log.Ldate
	TIME   = log.Ltime
	SFILE  = log.Lshortfile
	LFILE  = log.Llongfile
	MSEC   = log.Lmicroseconds
	COMMON = log.LUTC | DATE | TIME | SFILE
)

var (
	flags             = DATE | TIME | SFILE
	threshold  string = DEFAULT_THRESHOLD
	levelsEnum map[string]int

	DISCARD io.Writer = ioutil.Discard

	InTest    									*log.Logger = log.New(os.Stdout, "", COMMON)
	InProd     									*log.Logger = log.New(os.Stdout, "", COMMON)
	ToInvestigateTomorrow     	*log.Logger = log.New(os.Stdout, "", COMMON)
	WakeMeInTheMiddleOfTheNight *log.Logger = log.New(os.Stdout, "", COMMON)

	loggers map[string]*mlogger = map[string]*mlogger{
		"debug":    &mlogger{prefix: "DEBUG: ", writer: os.Stdout, logger: InTest},
		"info":     &mlogger{prefix: "INFO: ", writer: os.Stdout, logger: InProd},
		"warn":     &mlogger{prefix: "WARN: ", writer: os.Stdout, logger: ToInvestigateTomorrow},
		"error":    &mlogger{prefix: "ERROR: ", writer: os.Stdout, logger: WakeMeInTheMiddleOfTheNight},
	}
)

func init() {
	levelsEnum = make(map[string]int)

	// Enumerate the core logging levels so that thresholding can be applied
	for index, name := range []string{"debug", "info", "warn", "error"} {
		levelsEnum[name] = index
	}

	// Initialize default loggers using the values encoded in 'loggers' map
	for _, l := range loggers {
		l.logger.SetPrefix(l.prefix)
	}

	SetThreshold(DEFAULT_THRESHOLD)
}

func Threshold() string {
	return threshold
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
		InProd.Printf("cannot set output of unknown logger '%s'", logger)
		return
	}

	// Set the logger's writer (ignoring thresholding)
	switch len(writers) {
	case 0:
		InProd.Println("no io.Writer(s) provided for output", logger)
		return
	case 1:
		handler.writer = writers[0]
	default:
		handler.writer = io.MultiWriter(writers...)
	}

	// Update the handler when it falls within the current threshold
	if levelsEnum[logger] >= levelsEnum[threshold] {
		handler.logger.SetOutput(handler.writer)
	}
}

func SetThreshold(level string) {
	// Update overall level
	if _, ok := levelsEnum[level]; ok == false {
		InProd.Printf("ignoring invalid log level '%s'", level)
		return
	}

	threshold = level

	// Re-evaluate each default logger's threshold
	enum := levelsEnum[threshold]
	for key, l := range levelsEnum {
		logger := loggers[key].logger
		if l < enum {
			// Apply discard writer
			logger.SetOutput(DISCARD)
		} else {
			// Restore configured writer
			logger.SetOutput(loggers[key].writer)
		}
	}
}

func WithPrefix(enable bool) {
	for _, l := range loggers {
		prefix := ""
		if enable {
			prefix = l.prefix
		}
		l.logger.SetPrefix(prefix)
	}
}
