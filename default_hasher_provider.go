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

// defaultHasherProvider struct for default hasher provider.
// Used for testing purposes only.
type defaultHasherProvider struct {
	IonHasherProvider

	algorithm Algorithm
	provider  *testIonHasherProvider
}

func newDefaultHasherProvider(algo string, provider *testIonHasherProvider) *defaultHasherProvider {
	return &defaultHasherProvider{algorithm: Algorithm(algo), provider: provider}
}

// NewHasher returns a new defaultHasher.
func (dhp *defaultHasherProvider) NewHasher() (IonHasher, error) {
	ionHasher, err := newDefaultHasher(dhp.algorithm, dhp.provider)
	if err != nil {
		return nil, err
	}
	return ionHasher, nil
}
