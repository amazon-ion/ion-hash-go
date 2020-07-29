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

// testIonHasherProvider struct for test Ion hasher provider.
// Used for testing purposes only.
type testIonHasherProvider struct {
	algorithm     string
	updateHashLog [][]byte
	digestHashLog [][]byte
}

// newTestIonHasherProvider returns a new testIonHasherProvider.
func newTestIonHasherProvider(algorithm string) *testIonHasherProvider {
	return &testIonHasherProvider{algorithm: algorithm}
}

// getInstance returns either an identityHasherProvider or defaultHasherProvider depending on the algorithm.
func (tiop *testIonHasherProvider) getInstance() IonHasherProvider {
	if tiop.algorithm == "identity" {
		return newIdentityHasherProvider(tiop)
	}
	return newDefaultHasherProvider(strings.ToUpper(tiop.algorithm), tiop)
}

// getUpdateHashlog returns the updateHashLog.
func (tiop *testIonHasherProvider) getUpdateHashLog() [][]byte {
	return tiop.updateHashLog
}

// getDigestHashLog returns the digestHashLog.
func (tiop *testIonHasherProvider) getDigestHashLog() [][]byte {
	return tiop.digestHashLog
}

// getFinalDigestHashLog returns the last value in the digestHashLog.
func (tiop *testIonHasherProvider) getFinalDigestHashLog() []byte {
	return tiop.digestHashLog[len(tiop.digestHashLog)-1]
}
