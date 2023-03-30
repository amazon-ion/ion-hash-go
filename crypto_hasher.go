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
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/md4"       //nolint: staticcheck
	"golang.org/x/crypto/ripemd160" //nolint: staticcheck
	"golang.org/x/crypto/sha3"
)

// Algorithm is the name of the hash algorithm used to calculate the hash.
type Algorithm string

// Constants for each of the algorithm names supported.
const (
	MD4        Algorithm = "MD4"
	MD5        Algorithm = "MD5"
	SHA1       Algorithm = "SHA1"
	SHA224     Algorithm = "SHA224"
	SHA256     Algorithm = "SHA256"
	SHA384     Algorithm = "SHA384"
	SHA512     Algorithm = "SHA512"
	RIPEMD160  Algorithm = "RIPMD160"
	SHA3s224   Algorithm = "SHA3_224"
	SHA3s256   Algorithm = "SHA3_256"
	SHA3s384   Algorithm = "SHA3_384"
	SHA3s512   Algorithm = "SHA3_512"
	SHA512s224 Algorithm = "SHA512_224"
	SHA512s256 Algorithm = "SHA512_256"
	BLAKE2s256 Algorithm = "BLAKE2s_256"
	BLAKE2b256 Algorithm = "BLAKE2b_256"
	BLAKE2b384 Algorithm = "BLAKE2b_384"
	BLAKE2b512 Algorithm = "BLAKE2b_512"
)

// cryptoHasher computes the hash using given algorithm.
type cryptoHasher struct {
	hashAlgorithm hash.Hash
}

// newCryptoHasher returns a new cryptoHasher. Returns an error if the algorithm name provided is unknown.
// Here is a list of available hash functions: https://golang.org/pkg/crypto/#Hash.
func newCryptoHasher(algorithm Algorithm) (IonHasher, error) {
	var hashAlgorithm hash.Hash

	switch algorithm {
	case MD4:
		hashAlgorithm = md4.New()
	case MD5:
		hashAlgorithm = md5.New()
	case SHA1:
		hashAlgorithm = sha1.New()
	case SHA224:
		hashAlgorithm = sha256.New()
	case SHA256:
		hashAlgorithm = sha256.New()
	case SHA384:
		hashAlgorithm = sha512.New()
	case SHA512:
		hashAlgorithm = sha512.New()
	case RIPEMD160:
		hashAlgorithm = ripemd160.New()
	case SHA3s224:
		hashAlgorithm = sha3.New224()
	case SHA3s256:
		hashAlgorithm = sha3.New256()
	case SHA3s384:
		hashAlgorithm = sha3.New384()
	case SHA3s512:
		hashAlgorithm = sha3.New512()
	case SHA512s224:
		hashAlgorithm = sha512.New512_224()
	case SHA512s256:
		hashAlgorithm = sha512.New512_256()
	case BLAKE2s256:
		hashAlgorithm, _ = blake2s.New256(nil)
	case BLAKE2b256:
		hashAlgorithm, _ = blake2b.New256(nil)
	case BLAKE2b384:
		hashAlgorithm, _ = blake2b.New384(nil)
	case BLAKE2b512:
		hashAlgorithm, _ = blake2b.New512(nil)
	default:
		return nil, &InvalidArgumentError{"algorithm", algorithm}
	}

	ch := &cryptoHasher{hashAlgorithm}
	return ch, nil
}

// Write adds more data to the running hash.
func (ch *cryptoHasher) Write(b []byte) (n int, err error) {
	return ch.hashAlgorithm.Write(b)
}

// Sum appends the current hash to b and returns the resulting slice.
// It does not change the underlying hash state.
func (ch *cryptoHasher) Sum(b []byte) []byte {
	return ch.hashAlgorithm.Sum(b)
}

// Reset resets the Hash to its initial state.
func (ch *cryptoHasher) Reset() {
	ch.hashAlgorithm.Reset()
}
