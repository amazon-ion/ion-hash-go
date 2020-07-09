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


type defaultHasher struct {
	cryptoHasher IonHasher
	updateHashLog [][]byte
	digestHashLog [][]byte
}

func newDefaultIonHasher(algorithm algorithm) (IonHasher, error) {
	cryptoHasher, err := newCryptoHasher(algorithm)
	if err != nil {
		return nil, &InvalidArgumentError{"algorithm", algorithm}
	}
	return &defaultHasher{cryptoHasher: cryptoHasher}, nil
}

func (dh *defaultHasher) Write(b []byte) (n int, err error) {
	dh.updateHashLog = append(dh.updateHashLog, b)
	return dh.cryptoHasher.Write(b)
}

func (dh *defaultHasher) Sum(b []byte) []byte {
	hash := dh.cryptoHasher.Sum(b)
	dh.digestHashLog = append(dh.digestHashLog, hash)
	return hash
}

func (dh *defaultHasher) GetUpdateHashLog() [][]byte {
	return dh.updateHashLog
}

func (dh *defaultHasher) GetDigestHashLog() [][]byte {
	return dh.digestHashLog
}
