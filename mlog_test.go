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
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetLevel(t *testing.T) {
	cases := []struct {
		level    string
		expected string
	}{
		{LEVEL_PRODUCTION, LEVEL_PRODUCTION},
		{"whatever", LEVEL_PRODUCTION},
		{LEVEL_TOMORROW, LEVEL_TOMORROW},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.level), func(t *testing.T) {
			SetThreshold(tc.level)
			assert.Equal(t, Threshold(), tc.expected)
		})
	}
}

func TestDefaultThresholdLogging(t *testing.T) {
	SetThreshold(LEVEL_TOMORROW)
	cases := []struct {
		logger   *log.Logger
		name     string
		message  string
		expected bool
	}{
		{wakeMeInTheMiddleOfTheNight, "error", "standard error", true},
		{toInvestigateTomorrow, "warn", "warning message", true},
		{inProd, "info", "information", false},
		{inTest, "debug", "debug message", false},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.message), func(t *testing.T) {
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

func TestFlagSet(t *testing.T) {
	SetFlags(SFILE)
	SetFlags(NONE, LEVEL_MIDDLE_OF_NIGHT)
	cases := []struct {
		logger   *log.Logger
		name     string
		expected bool
	}{
		{toInvestigateTomorrow, "warn", true},
		{wakeMeInTheMiddleOfTheNight, "error", false},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			buffer := new(bytes.Buffer)
			SetOutput(tc.name, buffer)
			tc.logger.Println(tc.name, "test message")
			if tc.expected == true {
				assert.Contains(t, buffer.String(), ".go:")
			} else {
				assert.NotContains(t, buffer.String(), ".go:")
			}
		})
	}
}

func TestWriterOutput(t *testing.T) {
	cases := []struct {
		name     string
		logger	 *log.Logger
		multi    bool
		expected bool
	}{
		{LEVEL_TOMORROW, toInvestigateTomorrow, false, true},
		{"invalid", inTest, false, false},
		{LEVEL_TOMORROW, toInvestigateTomorrow, true, true},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			buffer := new(bytes.Buffer)
			extra := new(bytes.Buffer)
			if tc.multi == true {
				SetOutput(tc.name, buffer, extra)
			} else {
				SetOutput(tc.name, buffer)
			}
			tc.logger.Println(tc.name, "test message")
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

	// Empty logger
	buffer := new(bytes.Buffer)
	SetThreshold(LEVEL_PRODUCTION)
	SetOutput(LEVEL_PRODUCTION, buffer)
	SetOutput(LEVEL_TOMORROW)
	toInvestigateTomorrow.Println("empty")
	assert.Contains(t, buffer.String(), "no io.Writer(s) provided")
}
