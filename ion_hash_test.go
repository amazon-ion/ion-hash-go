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

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/amzn/ion-go/ion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIonHash(t *testing.T) {
	parameters := ionHashDataSource(t)

	for i := range parameters {
		ionBinary := parameters[i].testCase
		reader := ion.NewReaderBytes(ionBinary)

		provider := parameters[i].provider
		Traverse(t, reader, provider.getInstance())

		if len(parameters[i].expectedHashLog.identityUpdateList) > 0 {
			assert.Equal(t, provider.getUpdateHashLog(), parameters[i].expectedHashLog.identityUpdateList, parameters[i].hasherName+" failed")
		}
		if len(parameters[i].expectedHashLog.identityDigestList) > 0 {
			assert.Equal(t, provider.getDigestHashLog(), parameters[i].expectedHashLog.identityDigestList, parameters[i].hasherName+" failed")
		}
		if len(parameters[i].expectedHashLog.identityFinalDigest) > 0 {
			assert.Equal(t, provider.getDigestHashLog(), parameters[i].expectedHashLog.identityFinalDigest, parameters[i].hasherName+" failed")
		}

		if len(parameters[i].expectedHashLog.md5UpdateList) > 0 {
			assert.Equal(t, provider.getUpdateHashLog(), parameters[i].expectedHashLog.md5UpdateList, parameters[i].hasherName+" failed")
		}
		if len(parameters[i].expectedHashLog.md5DigestList) > 0 {
			assert.Equal(t, provider.getDigestHashLog(), parameters[i].expectedHashLog.md5DigestList, parameters[i].hasherName+" failed")
		}
	}
}

func Traverse(t *testing.T, reader ion.Reader, provider IonHasherProvider) {
	hr, err := NewHashReader(reader, provider)
	require.NoError(t, err, "Something went wrong executing NewHashReader()")
	TraverseReader(hr)
	hr.Sum(nil)
}

func TraverseReader(hr HashReader) {
	for hr.Next() {
		if hr.Type() != ion.NoType && hr.isInStruct() {
			hr.StepIn()
			TraverseReader(hr)
			hr.StepOut()
		}
	}
}

func ionHashDataSource(t *testing.T) []testObject {
	var dataList []testObject

	//todo revert to original file name
	file, err := ioutil.ReadFile("ion-hash-test/ion_hash_tests_2.ion")
	require.NoError(t, err, "Something went wrong loading ion_hash_tests.ion")

	reader := ion.NewReaderBytes(file)
	for reader.Next() {
		testName := "unknown"
		if reader.Annotations() != nil {
			testName = reader.Annotations()[0]
		}

		assert.NoError(t, reader.StepIn(), "Something went wrong executing reader.StepIn()")

		reader.Next() // Read the initial Ion value.

		testCase := []byte{}
		if reader.FieldName() == "10n" {
			assert.NoError(t, reader.StepIn(), "Something went wrong executing reader.StepIn()")

			for reader.Next() {
				intValue, err := reader.Int64Value()
				assert.NoError(t, err, "Something went wrong executing reader.IntValue()")
				testCase = append(testCase, byte(intValue))
			}
			assert.NoError(t, reader.StepOut(), "Something went wrong executing reader.StepOut()")

		} else {
			// Create textWriter to set testName to Ion text.
			str := strings.Builder{}
			textWriter := ion.NewTextWriter(&str)

			// Create binaryWriter to set testCase to Ion binary.
			buf := bytes.Buffer{}
			binaryWriter := ion.NewBinaryWriter(&buf)

			writeToWriter(t, reader, textWriter, binaryWriter)

			assert.NoError(t, textWriter.Finish(), "Something went wrong executing textWriter.Finish().")
			assert.NoError(t, binaryWriter.Finish(), "Something went wrong executing binaryWriter.Finish().")

			if testName == "unknown" {
				testName = str.String()
			}
			if len(testCase) == 0 {
				testCase = buf.Bytes()
			}
		}

		// Iterate through expected/ digest bytes.
		reader.Next()
		fieldName := reader.FieldName()

		if fieldName == "expect" {
			assert.NoError(t, reader.StepIn(), "Something went wrong executing reader.StepIn()")

			for reader.Next() {
				identityUpdateList := [][]byte{}
				identityDigestList := [][]byte{}
				identityFinalDigest := [][]byte{}
				md5UpdateList := [][]byte{}
				md5DigestList := [][]byte{}

				fieldName = reader.FieldName()
				hasherName := fieldName
				if fieldName == "identity" {
					assert.NoError(t, reader.StepIn(), "Something went wrong executing reader.StepIn()")

					for reader.Next() {
						annotations := reader.Annotations()

						if len(annotations) > 0 {
							if annotations[0] == "update" {
								updateBytes := readSexpAndAppendToList(t, reader)
								identityUpdateList = append(identityUpdateList, updateBytes)

							} else if annotations[0] == "digest" {
								digestBytes := readSexpAndAppendToList(t, reader)
								identityDigestList = append(identityDigestList, digestBytes)
							} else if annotations[0] == "final_digest" {
								digestBytes := readSexpAndAppendToList(t, reader)
								identityFinalDigest = append(identityFinalDigest, digestBytes)
							}
						}
					}
					assert.NoError(t, reader.StepOut(), "Something went wrong executing reader.StepOut()")
				} else if fieldName == "md5" {
					assert.NoError(t, reader.StepIn(), "Something went wrong executing reader.StepIn()")

					for reader.Next() {
						annotations := reader.Annotations()

						if len(annotations) > 0 {
							if annotations[0] == "update" {
								updateBytes := readSexpAndAppendToList(t, reader)
								md5UpdateList = append(md5UpdateList, updateBytes)

							} else if annotations[0] == "digest" {
								digestBytes := readSexpAndAppendToList(t, reader)
								md5DigestList = append(md5DigestList, digestBytes)
							}
						}
					}
					assert.NoError(t, reader.StepOut(), "Something went wrong executing reader.StepOut()")
				}

				expectedHashLog := hashLog{
					identityUpdateList:  identityUpdateList,
					identityDigestList:  identityDigestList,
					identityFinalDigest: identityFinalDigest,
					md5UpdateList:       md5UpdateList,
					md5DigestList:       md5DigestList,
				}

				if hasherName != "identity" {
					testName = testName + "." + hasherName
				}

				dataList = append(dataList, testObject{testName, testCase, expectedHashLog, newTestIonHasherProvider(hasherName)})
			}
			assert.NoError(t, reader.StepOut(), "Something went wrong executing reader.StepOut()")
		}
		assert.NoError(t, reader.StepOut(), "Something went wrong executing reader.StepOut()")
	}
	return dataList
}

type testObject struct {
	hasherName      string
	testCase        []byte
	expectedHashLog hashLog
	provider        *testIonHasherProvider
}

type hashLog struct {
	identityUpdateList  [][]byte
	identityDigestList  [][]byte
	identityFinalDigest [][]byte
	md5UpdateList       [][]byte
	md5DigestList       [][]byte
}
