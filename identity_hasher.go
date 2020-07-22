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

type identityHasher struct {
	IonHasher

	identityHash []byte
	provider     *testIonHasherProvider
}

func newIdentityIonHasher(provider *testIonHasherProvider) IonHasher {
	return &identityHasher{identityHash: []byte{}, provider: provider}
}

func (ih *identityHasher) Write(bytes []byte) (int, error) {
	for _, b := range bytes {
		ih.identityHash = append(ih.identityHash, b)
	}

	if bytes != nil {
		ih.provider.updateHashLog = append(ih.provider.updateHashLog, bytes)
	}
	return len(bytes), nil
}

func (ih *identityHasher) Sum(bytes []byte) []byte {
	// We ignore the error here because we know this particular Write() implementation does not error
	_, _ = ih.Write(bytes)

	identityHash := ih.identityHash
	ih.identityHash = []byte{}

	ih.provider.digestHashLog = append(ih.provider.digestHashLog, identityHash)
	return identityHash
}

func (ih *identityHasher) Reset() {
}
