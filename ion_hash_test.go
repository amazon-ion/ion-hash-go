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
			assert.Equal(t, parameters[i].expectedHashLog.identityUpdateList, provider.getUpdateHashLog(), parameters[i].hasherName+" failed")
		}
		if len(parameters[i].expectedHashLog.identityDigestList) > 0 {
			assert.Equal(t, parameters[i].expectedHashLog.identityDigestList, provider.getDigestHashLog(), parameters[i].hasherName+" failed")
		}
		if len(parameters[i].expectedHashLog.identityFinalDigestList) > 0 {
			assert.Equal(t, parameters[i].expectedHashLog.identityFinalDigestList, provider.getFinalDigestHashLog(), parameters[i].hasherName+" failed")
		}

		if len(parameters[i].expectedHashLog.md5UpdateList) > 0 {
			assert.Equal(t, parameters[i].expectedHashLog.md5UpdateList, provider.getDigestHashLog(), parameters[i].hasherName+" failed")
		}
		if len(parameters[i].expectedHashLog.md5DigestList) > 0 {
			assert.Equal(t, parameters[i].expectedHashLog.md5DigestList, provider.getDigestHashLog(), parameters[i].hasherName+" failed")
		}
	}
}

func Traverse(t *testing.T, reader ion.Reader, provider IonHasherProvider) {
	hr, err := NewHashReader(reader, provider)
	require.NoError(t, err, "Something went wrong executing NewHashReader()")

	TraverseReader(t, hr)

	_, err = hr.Sum(nil)
	require.NoError(t, err, "Something went wrong with executing hr.Sum()")
}

func TraverseReader(t *testing.T, hr HashReader) {
	for hr.Next() {
		if hr.Type() != ion.NoType && hr.IsInStruct() {
			require.NoError(t, hr.StepIn(), "Something went wrong executing hr.StepIn()")

			TraverseReader(t, hr)

			require.NoError(t, hr.StepOut(), "Something went wrong executing hr.StepOut()")
		}
	}
	require.NoError(t, hr.Err(), "Something went wrong executing hr.Next()")
}

func ionHashDataSource(t *testing.T) []testObject {
	var dataList []testObject

	file, err := ioutil.ReadFile("ion_hash_tests.ion")
	require.NoError(t, err, "Something went wrong loading ion_hash_tests.ion")

	reader := ion.NewReaderBytes(file)
	for reader.Next() {
		testName := "unknown"
		if reader.Annotations() != nil {
			testName = reader.Annotations()[0]
		}

		require.NoError(t, reader.StepIn(), "Something went wrong executing reader.StepIn()")

		require.True(t, reader.Next()) // Read the initial Ion value.

		testCase := []byte{}
		fieldName := reader.FieldName()
		if fieldName != nil && *fieldName == "10n" {
			require.NoError(t, reader.StepIn(), "Something went wrong executing reader.StepIn()")

			testCase = append(testCase, []byte{0xE0, 0x01, 0x00, 0xEA}...)
			for reader.Next() {
				intValue, err := reader.Int64Value()
				require.NoError(t, err, "Something went wrong executing reader.IntValue()")
				testCase = append(testCase, byte(intValue))
			}
			require.NoError(t, reader.Err(), "Something went wrong executing reader.Next()")
			require.NoError(t, reader.StepOut(), "Something went wrong executing reader.StepOut()")

		} else {
			// Create textWriter to set testName to Ion text.
			str := strings.Builder{}
			textWriter := ion.NewTextWriter(&str)

			// Create binaryWriter to set testCase to Ion binary.
			buf := bytes.Buffer{}
			binaryWriter := ion.NewBinaryWriter(&buf)

			writeToWriters(t, reader, textWriter, binaryWriter)

			require.NoError(t, textWriter.Finish(), "Something went wrong executing textWriter.Finish().")
			require.NoError(t, binaryWriter.Finish(), "Something went wrong executing binaryWriter.Finish().")

			if testName == "unknown" {
				testName = str.String()
			}
			if len(testCase) == 0 {
				testCase = buf.Bytes()
			}
		}

		require.True(t, reader.Next()) // Iterate through expected/ digest bytes.

		fieldName = reader.FieldName()
		if fieldName != nil && *fieldName == "expect" {
			require.NoError(t, reader.StepIn(), "Something went wrong executing reader.StepIn()")

			for reader.Next() {
				hasherName := reader.FieldName()
				if hasherName == nil {
					continue
				}

				identityUpdateList := [][]byte{}
				identityDigestList := [][]byte{}
				identityFinalDigestList := []byte{}
				md5UpdateList := [][]byte{}
				md5DigestList := [][]byte{}

				switch *hasherName {
				case "identity":
					require.NoError(t, reader.StepIn(), "Something went wrong executing reader.StepIn()")

					for reader.Next() {
						annotations := reader.Annotations()

						if len(annotations) > 0 {
							switch annotations[0] {
							case "update":
								updateBytes := readSexpAndAppendToList(t, reader)
								identityUpdateList = append(identityUpdateList, updateBytes)
							case "digest":
								digestBytes := readSexpAndAppendToList(t, reader)
								identityDigestList = append(identityDigestList, digestBytes)
							case "final_digest":
								digestBytes := readSexpAndAppendToList(t, reader)
								identityFinalDigestList = digestBytes
							}
						}
					}
					require.NoError(t, reader.Err(), "Something went wrong executing reader.Next()")
					require.NoError(t, reader.StepOut(), "Something went wrong executing reader.StepOut()")
				case "md5":
					require.NoError(t, reader.StepIn(), "Something went wrong executing reader.StepIn()")

					for reader.Next() {
						annotations := reader.Annotations()

						if len(annotations) > 0 {
							switch annotations[0] {
							case "update":
								updateBytes := readSexpAndAppendToList(t, reader)
								md5UpdateList = append(md5UpdateList, updateBytes)
							case "digest":
								digestBytes := readSexpAndAppendToList(t, reader)
								md5DigestList = append(md5DigestList, digestBytes)
							}
						}
					}
					require.NoError(t, reader.Err(), "Something went wrong executing reader.Next()")
					require.NoError(t, reader.StepOut(), "Something went wrong executing reader.StepOut()")
				}

				if *hasherName != "identity" {
					testName = testName + "." + *hasherName
				}

				expectedHashLog := hashLog{
					identityUpdateList:      identityUpdateList,
					identityDigestList:      identityDigestList,
					identityFinalDigestList: identityFinalDigestList,
					md5UpdateList:           md5UpdateList,
					md5DigestList:           md5DigestList,
				}

				dataList = append(dataList, testObject{testName, testCase, &expectedHashLog, newTestIonHasherProvider(*hasherName)})
			}
			require.NoError(t, reader.Err(), "Something went wrong executing reader.Next()")
			require.NoError(t, reader.StepOut(), "Something went wrong executing reader.StepOut()")
		}
		require.NoError(t, reader.StepOut(), "Something went wrong executing reader.StepOut()")
	}
	require.NoError(t, reader.Err(), "Something went wrong executing reader.Next()")

	return dataList
}

type testObject struct {
	hasherName      string
	testCase        []byte
	expectedHashLog *hashLog
	provider        *testIonHasherProvider
}

type hashLog struct {
	identityUpdateList      [][]byte
	identityDigestList      [][]byte
	identityFinalDigestList []byte
	md5UpdateList           [][]byte
	md5DigestList           [][]byte
}
