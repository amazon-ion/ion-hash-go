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

import "github.com/amzn/ion-hash-go/ihp"

type cryptoHasherProvider struct {
	ihp.IonHasherProvider

	algorithm algorithm
}

func newCryptoHasherProvider(algorithm algorithm) *cryptoHasherProvider {
	return &cryptoHasherProvider{algorithm: algorithm}
}

func (chp *cryptoHasherProvider) NewHasher() (ihp.IonHasher, error) {
	return newCryptoHasher(chp.algorithm)
}
