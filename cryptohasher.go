/*
 * Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/md4"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
	"hash"
)

type algorithm int

const (
	MD4  algorithm = iota
	MD5
	SHA1
	SHA224
	SHA256
	SHA384
	SHA512
	RIPEMD160
	SHA3_224
	SHA3_256
	SHA3_384
	SHA3_512
	SHA512_224
	SHA512_256
	BLAKE2s_256
	BLAKE2b_256
	BLAKE2b_384
	BLAKE2b_512
)

type cryptoHasher struct {
	hashAlgorithm hash.Hash
}

func NewCryptoHasher(algorithm algorithm) (IonHasher, error) {
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
	case SHA3_224:
		hashAlgorithm = sha3.New224()
	case SHA3_256:
		hashAlgorithm = sha3.New256()
	case SHA3_384:
		hashAlgorithm = sha3.New384()
	case SHA3_512:
		hashAlgorithm = sha3.New512()
	case SHA512_224:
		hashAlgorithm = sha512.New512_224()
	case SHA512_256:
		hashAlgorithm = sha512.New512_256()
	case BLAKE2s_256:
		hashAlgorithm, _ = blake2s.New256(nil)
	case BLAKE2b_256:
		hashAlgorithm, _ = blake2b.New256(nil)
	case BLAKE2b_384:
		hashAlgorithm, _ = blake2b.New384(nil)
	case BLAKE2b_512:
		hashAlgorithm, _ = blake2b.New512(nil)
	default:
		return nil, &SyntaxError{"Hashing algorithm is invalid"}
	}

	ch := &cryptoHasher{hashAlgorithm}
	return ch, nil
}

func (ch *cryptoHasher) Write(p []byte) (n int, err error) {
	return ch.hashAlgorithm.Write(p)
}

func (ch *cryptoHasher) Sum(b []byte) []byte {
	return ch.hashAlgorithm.Sum(b)
}
