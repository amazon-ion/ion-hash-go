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

import "testing"

var skipList = []string{
	"TestConsumeRemainderPartialConsume",        // hash_reader_test
	"TestConsumeRemainderStepInStepOutNested",   // hash_reader_test
	"TestConsumeRemainderStepInNextStepOut",     // hash_reader_test
	"TestConsumeRemainderStepInStepOutTopLevel", // hash_reader_test
	"TestConsumeRemainderNext",                  // hash_reader_test
	"TestReaderUnresolvedSid",                   // hash_reader_test
	"TestIonReaderContract",                     // hash_reader_test
	"TestMiscMethods",                           // hash_writer_test
	"TestIonWriterContractWriteValue",           // hash_writer_test
	"TestIonWriterContractWriteValues",          // hash_writer_test
	"TestWriterUnresolvedSid",                   // hash_writer_test
}

func checkTestToSkip(t *testing.T) {
	for _, fileName := range skipList {
		if fileName == t.Name() {
			t.Skip()
		}
	}
}
