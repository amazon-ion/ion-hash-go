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

// CryptoHasherProvider struct for crypto hasher provider
type CryptoHasherProvider struct {
	IonHasherProvider

	algorithm algorithm
}

// NewCryptoHasherProvider returns a new CryptoHasherProvider.
func NewCryptoHasherProvider(algorithm algorithm) *CryptoHasherProvider {
	return &CryptoHasherProvider{algorithm: algorithm}
}

// NewHasher returns a new cryptoHasher.
func (chp *CryptoHasherProvider) NewHasher() (IonHasher, error) {
	return newCryptoHasher(chp.algorithm)
}
