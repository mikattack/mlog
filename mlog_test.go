/*
 * Copyright © 2016 Alex Mikitik.
 *
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package mlog

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobalSetThreshold(t *testing.T) {
	buffer := new(bytes.Buffer)
	SetOutput(buffer)
	cases := []struct {
		level    LogLevel
		expected LogLevel
	}{
		{IN_PRODUCTION, IN_PRODUCTION},
		{"whatever", IN_PRODUCTION},
		{TO_INVESTIGATE, TO_INVESTIGATE},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(string(tc.level)), func(t *testing.T) {
			SetThreshold(tc.level)
			assert.Equal(t, Threshold(), tc.expected)
		})
	}
}

func TestGlobalLogging(t *testing.T) {
	buffer := new(bytes.Buffer)
	SetFlags(0)
	SetOutput(buffer)
	SetThreshold(IN_TESTING)

	assert.Equal(t, Flags(), 0)
	assert.Equal(t, Threshold(), IN_TESTING)

	InTesting("debug")
	assert.Equal(t, buffer.String(), "debug\n")
	buffer.Reset()

	InProduction("info")
	assert.Equal(t, buffer.String(), "info\n")
	buffer.Reset()

	ToInvestigate("warning")
	assert.Equal(t, buffer.String(), "warning\n")
	buffer.Reset()

	PageMeNow("error")
	assert.Equal(t, buffer.String(), "error\n")
}

func TestThresholdLogging(t *testing.T) {
	cases := []struct {
		level    LogLevel
		message  string
		expected bool
	}{
		{PAGE_ME_NOW, "standard error", true},
		{TO_INVESTIGATE, "warning message", true},
		{IN_PRODUCTION, "information", false},
		{IN_TESTING, "debug message", false},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.message), func(t *testing.T) {
			buffer := new(bytes.Buffer)
			logger := New(buffer, 0)
			logger.SetThreshold(TO_INVESTIGATE)
			logger.SetOutput(buffer)
			logger.log(default_call_depth, tc.level, tc.message)
			if tc.expected == true {
				assert.Contains(t, buffer.String(), tc.message)
			} else {
				assert.NotContains(t, buffer.String(), tc.message)
			}
		})
	}
}

func TestLoggingFunctions(t *testing.T) {
	buffer := new(bytes.Buffer)
	logger := New(buffer, 0)
	logger.SetThreshold(IN_TESTING)

	logger.InTesting("debug")
	assert.Equal(t, buffer.String(), "debug\n")
	buffer.Reset()

	logger.InProduction("info")
	assert.Equal(t, buffer.String(), "info\n")
	buffer.Reset()

	logger.ToInvestigate("warning")
	assert.Equal(t, buffer.String(), "warning\n")
	buffer.Reset()

	logger.PageMeNow("error")
	assert.Equal(t, buffer.String(), "error\n")
}

func TestFormattedLogging(t *testing.T) {
	buffer := new(bytes.Buffer)
	logger := New(buffer, 0)
	logger.SetThreshold(IN_TESTING)
	logger.InTesting("example: %d", 42)
	assert.Equal(t, buffer.String(), "example: 42\n")
}

func TestFlagSet(t *testing.T) {
	var validator = regexp.MustCompile(`^\[[A-Z]+\] .+\.go:[0-9]+: test message\n?$`)
	tmp := new(bytes.Buffer)
	logger := New(tmp, 0)
	logger.SetFlags(LEVEL | FILE)
	logger.SetThreshold(TO_INVESTIGATE)
	cases := []struct {
		level  LogLevel
		name   string
		logger *Logger
	}{
		{TO_INVESTIGATE, "warn", logger},
		{PAGE_ME_NOW, "error", logger},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			buffer := new(bytes.Buffer)
			tc.logger.SetOutput(buffer)
			tc.logger.log(default_call_depth, tc.level, "test message")
			validation := validator.MatchString(buffer.String())
			assert.Equal(t, validation, true)
		})
	}
}

func TestWriterOutput(t *testing.T) {
	cases := []struct {
		name     string
		level    LogLevel
		multi    bool
		expected bool
	}{
		{"Investigate", TO_INVESTIGATE, false, true},
		{"invalid", IN_TESTING, false, false},
		{"Investigate", TO_INVESTIGATE, true, true},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			buffer := new(bytes.Buffer)
			extra := new(bytes.Buffer)
			logger := New(buffer, 0)
			if tc.multi == true {
				logger.SetOutput(buffer, extra)
			} else {
				logger.SetOutput(buffer)
			}
			logger.log(default_call_depth, tc.level, "test message")
			if tc.expected == true {
				assert.Contains(t, buffer.String(), "test message")
				if tc.multi == true {
					assert.Contains(t, extra.String(), "test message")
				}
			} else {
				assert.NotContains(t, buffer.String(), "test message")
			}
		})
	}
}

func TestEmptyOutput(t *testing.T) {
	buffer := new(bytes.Buffer)
	logger := New(buffer, 0)
	logger.SetThreshold(IN_PRODUCTION)
	logger.SetOutput(buffer) // Set a fallback buffer
	logger.SetOutput()
	assert.Contains(t, buffer.String(), "SetOutput: no io.Writer(s) provided")
}

func TestPrefix(t *testing.T) {
	cases := []struct {
		name    string
		enabled bool
	}{
		{"with prefix", true},
		{"without prefix", false},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			buffer := new(bytes.Buffer)
			logger := New(buffer, 0)
			if tc.enabled {
				logger.SetFlags(LEVEL)
			} else {
				logger.SetFlags(FILE)
			}
			logger.log(default_call_depth, PAGE_ME_NOW, "test message")
			if tc.enabled {
				assert.Contains(t, buffer.String(), "[ERROR]")
			} else {
				assert.NotContains(t, buffer.String(), "[ERROR]")
			}
		})
	}
}

func TestInvalidLogLevel(t *testing.T) {
	buffer := new(bytes.Buffer)
	logger := New(buffer, 0)
	var NONSENSE LogLevel = "nonsense"

	logger.SetThreshold(NONSENSE)
	assert.Contains(t, buffer.String(), "Invalid log level: nonsense\n")
	buffer.Reset()

	logger.log(default_call_depth, NONSENSE, "invalid log level")
	assert.Contains(t, buffer.String(), "Invalid log level: nonsense\n")
}

func TestLoggingLineNumber(t *testing.T) {
	// Hard-coding line numbers doesn't end up working outside of this file

	buffer := new(bytes.Buffer)
	logger := New(buffer, FILE)
	logger.ToInvestigate("test")
	assert.Contains(t, buffer.String(), "_test.go:232")
	buffer.Reset()

	var NONSENSE LogLevel = "nonsense"
	logger.SetThreshold(NONSENSE)
	assert.Contains(t, buffer.String(), "mlog.go:")
}

func TestInvalidCallDepth(t *testing.T) {
	buffer := new(bytes.Buffer)
	logger := New(buffer, FILE)
	logger.log(200000000000, TO_INVESTIGATE, "invalid call depth")
	assert.Contains(t, buffer.String(), "???:0")
}
