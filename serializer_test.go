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

func TestEscape(t *testing.T) {
	// null case
	assert.Nil(t, escape(nil))

	// happy cases
	var empty []byte
	assert.Equal(t, empty, escape(empty))

	bytes := []byte{0x10, 0x11, 0x12, 0x13}
	assert.Equal(t, bytes, escape(bytes))

	// escape cases
	assert.Equal(t, []byte{0x0C, 0x0B}, escape([]byte{0x0B}))
	assert.Equal(t, []byte{0x0C, 0x0E}, escape([]byte{0x0E}))
	assert.Equal(t, []byte{0x0C, 0x0C}, escape([]byte{0x0C}))

	assert.Equal(t, []byte{0x0C, 0x0B, 0x0C, 0x0E, 0x0C, 0x0C}, escape([]byte{0x0B, 0x0E, 0x0C}))

	assert.Equal(t, []byte{0x0C, 0x0C, 0x0C, 0x0C}, escape([]byte{0x0C, 0x0C}))

	assert.Equal(t, []byte{0x0C, 0x0C, 0x10, 0x0C, 0x0C, 0x11, 0x0C, 0x0C, 0x12, 0x0C, 0x0C},
		escape([]byte{0x0C, 0x10, 0x0C, 0x11, 0x0C, 0x12, 0x0C}))
}
