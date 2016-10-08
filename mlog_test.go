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
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestSetLevel(t *testing.T) {
	cases := []struct {
		level			string
		expected	string
	}{
		{LEVEL_INFO, LEVEL_INFO},
		{"whatever", LEVEL_INFO},
		{LEVEL_WARN, LEVEL_WARN},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.level), func (t *testing.T) {
			SetThreshold(tc.level)
			assert.Equal(t, Threshold(), tc.expected)
		})
	}
}


func TestDefaultLogging(t *testing.T) {
	SetThreshold(DEFAULT_THRESHOLD)
	cases := []struct {
		logger		*mlogger
		name			string
		message		string
		expected	bool	
	}{
		{FATAL, "fatal", "fatal error", true},
		{CRITICAL, "critical", "critical error", true},
		{ERROR, "error", "standard error", true},
		{WARN, "warn", "warning message", true},
		{INFO, "info", "information", false},
		{DEBUG, "debug", "debug message", false},
		{TRACE, "trace", "trace message", false},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.message), func (t *testing.T) {
			buffer := new(bytes.Buffer)
			SetOutput(tc.name, buffer)
			tc.logger.Println(tc.message)
			if tc.expected == true {
				assert.Contains(t, buffer.String(), tc.message)
			} else {
				assert.NotContains(t, buffer.String(), tc.message)
			}
		})
	}
}


/*
 * 
func TestCustomLogging(t *testing.T) {
	ln := "test-logger"
	NewLogger(ln, "TEST: ")
	cases := []struct {
		name string
		message string
	}{
		{ln, "custom message"},
		{"warn", "warning message"},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.name), func (t *testing.T) {
			buffer := new(bytes.Buffer)
			SetOutput(tc.name, buffer)
			Println(tc.name, tc.message)
			assert.Contains(t, buffer.String(), tc.message)
		})
	}
}
 *
 */


func TestFlagSet(t *testing.T) {
	ln := "test-logger"
	SetFlags(SFILE)
	SetFlags(NONE, ln)
	cases := []struct {
		name			string
		expected	bool
	}{
		{ln, false},
		{"warn", true},
		{"fatal", true},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.name), func (t *testing.T) {
			buffer := new(bytes.Buffer)
			SetOutput(tc.name, buffer)
			Println(tc.name, "test message")
			if tc.expected == true {
				assert.Contains(t, buffer.String(), ".go:")
			} else {
				assert.NotContains(t, buffer.String(), ".go:")
			}
		})
	}
}
