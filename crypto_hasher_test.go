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
	"github.com/stretchr/testify/require"
)

func TestInvalidAlgorithm(t *testing.T) {
	hasherProvider := NewCryptoHasherProvider("invalid algorithm")

	_, err := hasherProvider.NewHasher()
	assert.Error(t, err, "Expected hasherProvider.NewHasher() to return an error")
	assert.IsType(t, &InvalidArgumentError{}, err, "Expected hasherProvider.NewHasher() to return InvalidArgumentError")
}

func TestHasher(t *testing.T) {
	// Using flawed MD5 algorithm FOR TEST PURPOSES ONLY
	hasherProvider := NewCryptoHasherProvider(MD5)

	h, err := hasherProvider.NewHasher()
	require.NoError(t, err, "Expected NewHasher() to successfully create a Hasher")

	hasher, ok := h.(*cryptoHasher)
	require.True(t, ok)

	emptyHasherDigest := hasher.Sum(nil)

	_, err = hasher.Write([]byte{0x0f})
	assert.NoError(t, err, "Something went wrong executing hasher.Write([]byte{0x0f})")

	expected := []byte{0xd8, 0x38, 0x69, 0x1e, 0x5d, 0x4a, 0xd0, 0x68, 0x79, 0xca, 0x72, 0x14, 0x42, 0xe8, 0x83, 0xd4}

	assert.Equal(t, expected, hasher.Sum(nil), "sum did not match expectation")

	// Verify that the hasher has reset
	assert.Equal(t, emptyHasherDigest, hasher.Sum(nil), "sum did not match expectation")
}
