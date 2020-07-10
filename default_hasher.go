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
	provider testIonHasherProvider
}

func newDefaultIonHasher(algorithm algorithm, provider testIonHasherProvider) (IonHasher, error) {
	cryptoHasher, err := newCryptoHasher(algorithm)
	if err != nil {
		return nil, &InvalidArgumentError{"algorithm", algorithm}
	}
	return &defaultHasher{cryptoHasher: cryptoHasher, provider: provider}, nil
}

func (dh *defaultHasher) Write(b []byte) (n int, err error) {
	dh.provider.updateHashLog = append (dh.provider.updateHashLog, b)
	return dh.cryptoHasher.Write(b)
}

func (dh *defaultHasher) Sum(b []byte) []byte {
	hash := dh.cryptoHasher.Sum(b)
	dh.provider.digestHashLog = append (dh.provider.digestHashLog, hash)
	return hash
}
