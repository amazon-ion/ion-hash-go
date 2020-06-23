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

type identityHasher struct {
	ihp.IonHasher

	identityHash []byte
}

func newIdentityIonHasher() ihp.IonHasher {
	return &identityHasher{identityHash: []byte{}}
}

func (identityHasher *identityHasher) Write(bytes []byte) (int, error) {
	for _, b := range bytes {
		identityHasher.identityHash = append(identityHasher.identityHash, b)
	}

	return len(bytes), nil
}

func (identityHasher *identityHasher) Sum(bytes []byte) []byte {
	// We ignore the error here because we know this particular Write() implementation does not error
	_, _ = identityHasher.Write(bytes)

	identityHash := identityHasher.identityHash
	identityHasher.identityHash = []byte{}

	return identityHash
}
