/*
 * Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License").
 * You may not use this file except in compliance with the License.
 * A copy of the License is located at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * or in the "license" file accompanying this file. This file is distributed
 * on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package ionhash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNullOrEmptyComparisons(t *testing.T) {
	assert.Equal(t, 0, compareBytes(nil, nil), "Result of compareBytes(nil, nil) was not as expected")

	assert.Equal(t, 0, compareBytes(nil, []byte{}), "Result of compareBytes(nil, []byte{}) was not as expected")

	assert.Equal(t, 0, compareBytes([]byte{}, nil), "Result of compareBytes([]byte{}, nil) was not as expected")
}

func TestIdentity(t *testing.T) {
	var emptyByteArray []byte
	assert.Equal(t, 0, compareBytes(emptyByteArray, emptyByteArray))

	bytes := []byte{0x01, 0x02, 0x03}
	assert.Equal(t, 0, compareBytes(bytes, bytes))
}

func TestEquals(t *testing.T) {
	assert.Equal(t, 0, compareBytes([]byte{0x01, 0x02, 0x03}, []byte{0x01, 0x02, 0x03}))
}

func TestLessThan(t *testing.T) {
	assert.Equal(t, -1, compareBytes([]byte{0x01, 0x02, 0x03}, []byte{0x01, 0x02, 0x04}))
}

func TestGreaterThan(t *testing.T) {
	assert.Equal(t, 1, compareBytes([]byte{0x01, 0x02, 0x04}, []byte{0x01, 0x02, 0x03}))
}

func TestLessThanDueToLength(t *testing.T) {
	assert.Equal(t, -1, compareBytes([]byte{0x01, 0x02, 0x03}, []byte{0x01, 0x02, 0x03, 0x04}))
}

func TestGreaterThanDueToLength(t *testing.T) {
	assert.Equal(t, 1, compareBytes([]byte{0x01, 0x02, 0x03, 0x04}, []byte{0x01, 0x02, 0x03}))
}
