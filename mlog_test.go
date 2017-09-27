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

func TestThresholdLogging(t *testing.T) {
	std.SetThreshold(TO_INVESTIGATE)
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
			std.SetOutput(buffer)
			std.log(tc.level, tc.message)
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
	std.SetFlags(0)
	std.SetThreshold(IN_TESTING)
	std.SetOutput(buffer)

	std.InTesting("debug")
	assert.Equal(t, buffer.String(), "debug\n")
	buffer.Reset()

	std.InProduction("info")
	assert.Equal(t, buffer.String(), "info\n")
	buffer.Reset()

	std.ToInvestigate("warning")
	assert.Equal(t, buffer.String(), "warning\n")
	buffer.Reset()

	std.PageMeNow("error")
	assert.Equal(t, buffer.String(), "error\n")
}

func TestFormattedLogging(t *testing.T) {
	buffer := new(bytes.Buffer)
	std.SetFlags(0)
	std.SetThreshold(IN_TESTING)
	std.SetOutput(buffer)

	std.InTesting("example: %d", 42)
	assert.Equal(t, buffer.String(), "example: 42\n")
}

func TestFlagSet(t *testing.T) {
	var validator = regexp.MustCompile(`^\[[A-Z]+\] .+\.go:[0-9]+: test message\n?$`)
	std.SetFlags(LEVEL | FILE)
	std.SetThreshold(TO_INVESTIGATE)
	cases := []struct {
		level LogLevel
		name  string
	}{
		{TO_INVESTIGATE, "warn"},
		{PAGE_ME_NOW, "error"},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			buffer := new(bytes.Buffer)
			std.SetOutput(buffer)
			std.log(tc.level, "test message")
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
			if tc.multi == true {
				std.SetOutput(buffer, extra)
			} else {
				std.SetOutput(buffer)
			}
			std.log(tc.level, "test message")
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
	std.SetThreshold(IN_PRODUCTION)
	std.SetOutput(buffer) // Set a fallback buffer
	std.SetOutput()
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
			std.SetOutput(buffer)
			if tc.enabled {
				SetFlags(LEVEL)
			} else {
				SetFlags(FILE)
			}
			std.log(PAGE_ME_NOW, "test message")
			if tc.enabled {
				assert.Contains(t, buffer.String(), "[ERROR]")
			} else {
				assert.NotContains(t, buffer.String(), "[ERROR]")
			}
		})
	}
}
