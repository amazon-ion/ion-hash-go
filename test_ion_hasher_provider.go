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

import "strings"

type testIonHasherProvider struct {

	algorithm string
	updateHashLog [][]byte
	digestHashLog [][]byte
}

func newTestIonHasherProvider(algorithm string) *testIonHasherProvider {
	return &testIonHasherProvider{algorithm: algorithm}
}

func (tiop *testIonHasherProvider)getInstance() IonHasherProvider {
	if tiop.algorithm == "identity" {
		return newIdentityHasherProvider(tiop)
	} else {
		return newDefaultHasherProvider(strings.ToUpper(tiop.algorithm), *tiop)
	}
}

func (tiop *testIonHasherProvider)getUpdateHashLog() [][]byte {
	return tiop.updateHashLog
}
func (tiop *testIonHasherProvider)getDigestHashLog() [][]byte {
	return tiop.digestHashLog
}
