/*
 * Copyright © 2016 Alex Mikitik.
 *
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package mlog

import (
  //"bytes"
  "testing"

  "github.com/stretchr/testify/assert"
)


func TestLevels(t *testing.T) {
	SetThreshold(LEVEL_ERROR)
	assert.Equal(t, Threshold(), LEVEL_ERROR)

	SetThreshold(LEVEL_CRITICAL)
	assert.Equal(t, Threshold(), LEVEL_CRITICAL)

	SetThreshold(LEVEL_WARN)
	assert.Equal(t, Threshold(), LEVEL_WARN)
}


func TestDefaultLogging(t *testing.T) {

}


func TestCustomLogging(t *testing.T) {

}
