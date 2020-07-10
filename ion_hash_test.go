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
	"fmt"
	"github.com/amzn/ion-go/ion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"strings"
	"testing"
)

func TestIonHash(t *testing.T) {
	parameters := ionHashDataSource(t)

	for i := range parameters {
		ionBinary := parameters[i].testCase
		reader := ion.NewReaderBytes(ionBinary)

		provider := parameters[i].provider
		Traverse(t, reader, provider.getInstance())

		if len(parameters[i].expectedHashLog.identityUpdateList) > 0 {
			assert.Equal(t, provider.getUpdateHashLog(), parameters[i].expectedHashLog.identityUpdateList)
		}
		if len(parameters[i].expectedHashLog.identityDigestList) > 0 {
			assert.Equal(t, provider.getDigestHashLog(), parameters[i].expectedHashLog.identityDigestList)
		}
		if len(parameters[i].expectedHashLog.identityFinalDigest) > 0 {
			assert.Equal(t, provider.getDigestHashLog(), parameters[i].expectedHashLog.identityDigestList)
		}

		if len(parameters[i].expectedHashLog.md5UpdateList) > 0 {
			assert.Equal(t, provider.getUpdateHashLog(), parameters[i].expectedHashLog.md5UpdateList)
		}
		if len(parameters[i].expectedHashLog.md5DigestList) > 0 {
			assert.Equal(t, provider.getDigestHashLog(), parameters[i].expectedHashLog.md5DigestList)
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

		err := reader.StepIn()
		assert.NoError(t, err, "Something went wrong executing reader.StepIn()")

		reader.Next() // Read the initial Ion value.

		testCase := []byte{}
		if reader.FieldName() == "10n" {
			err = reader.StepIn()
			assert.NoError(t, err, "Something went wrong executing reader.StepIn()")

			for reader.Next() {
				intValue, err := reader.Int64Value()
				assert.NoError(t, err, "Something went wrong executing reader.IntValue()")
				testCase = append(testCase, byte(intValue))
			}
			err = reader.StepOut()
			assert.NoError(t, err, "Something went wrong executing reader.StepOut()")

		} else {
			// Create textWriter to set testName to Ion text.
			str := strings.Builder{}
			textWriter := ion.NewTextWriter(&str)

			// Create binaryWriter to set testCase to Ion binary.
			buf := bytes.Buffer{}
			binaryWriter := ion.NewBinaryWriter(&buf)

			writeToWriter(t, reader, textWriter, binaryWriter)

			err = textWriter.Finish()
			assert.NoError(t, err, "Something went wrong writing Ion value to text writer.")
			err = binaryWriter.Finish()
			assert.NoError(t, err, "Something went wrong writing Ion value to binary writer.")

			if testName == "unknown" {
				testName = str.String()
			}
			if len(testCase) == 0 {
				testCase = buf.Bytes()
			}
		}

		//todo remove
		fmt.Println("testName: " + testName)
		fmt.Println(testCase)
		//iterate through sexp (expected, digest).
		reader.Next()
		fieldName := reader.FieldName()

		if fieldName == "expect" {
			err = reader.StepIn()
			assert.NoError(t, err, "Something went wrong executing reader.StepIn()")

			for reader.Next() {
				identityUpdateList := [][]byte{}
				identityDigestList := [][]byte{}
				identityFinalDigest := [][]byte{}
				md5UpdateList := [][]byte{}
				md5DigestList := [][]byte{}

				fieldName = reader.FieldName()
				hasherName := fieldName
				if fieldName == "identity" {
					err = reader.StepIn()
					assert.NoError(t, err, "Something went wrong executing reader.StepIn()")

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
					err = reader.StepOut()
					assert.NoError(t, err, "Something went wrong executing reader.StepOut()")
				} else if fieldName == "md5" {
					err = reader.StepIn()
					assert.NoError(t, err, "Something went wrong executing reader.StepIn()")

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
					err = reader.StepOut()
					assert.NoError(t, err, "Something went wrong executing reader.StepOut()")
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
			err = reader.StepOut()
			assert.NoError(t, err, "Something went wrong executing reader.StepOut()")
		}
		err = reader.StepOut()
		assert.NoError(t, err, "Something went wrong executing reader.StepOut()")
	}
	return dataList
}

type testObject struct {
	hasherName      string
	testCase        []byte
	expectedHashLog hashLog
	provider        *testIonHasherProvider
}

func writeToWriter(t *testing.T, reader ion.Reader, textWriter ion.Writer, binaryWriter ion.Writer) {
	ionType := reader.Type()

	if reader.Annotations() != nil {
		err := textWriter.Annotations(reader.Annotations()...)
		require.NoError(t, err)
		err = binaryWriter.Annotations(reader.Annotations()...)
		require.NoError(t, err)
	}

	switch ionType {
	case ion.NullType:
		if reader.FieldName() != "" {
			err := textWriter.FieldName(reader.FieldName())
			require.NoError(t, err)
			err = binaryWriter.FieldName(reader.FieldName())
			require.NoError(t, err)
		}
		if reader.IsNull() {
			err := textWriter.WriteNullType(ion.NullType)
			require.NoError(t, err)
			err = binaryWriter.WriteNullType(ion.NullType)
			require.NoError(t, err)
		} else {
			err := textWriter.WriteNull()
			require.NoError(t, err)
			err = binaryWriter.WriteNull()
			require.NoError(t, err)
		}

	case ion.BoolType:
		if reader.FieldName() != "" {
			err := textWriter.FieldName(reader.FieldName())
			require.NoError(t, err)
			err = binaryWriter.FieldName(reader.FieldName())
			require.NoError(t, err)
		}
		if reader.IsNull() {
			err := textWriter.WriteNullType(ion.BoolType)
			require.NoError(t, err)
			err = binaryWriter.WriteNullType(ion.BoolType)
			require.NoError(t, err)
		} else {
			boolValue, err := reader.BoolValue()
			require.NoError(t, err)
			err = textWriter.WriteBool(boolValue)
			require.NoError(t, err)
			err = binaryWriter.WriteBool(boolValue)
			require.NoError(t, err)
		}

	case ion.BlobType:
		if reader.FieldName() != "" {
			err := textWriter.FieldName(reader.FieldName())
			require.NoError(t, err)
			err = binaryWriter.FieldName(reader.FieldName())
			require.NoError(t, err)
		}
		if reader.IsNull() {
			err := textWriter.WriteNullType(ion.BlobType)
			require.NoError(t, err)
			err = binaryWriter.WriteNullType(ion.BlobType)
			require.NoError(t, err)
		} else {
			byteValue, err := reader.ByteValue()
			require.NoError(t, err)
			err = textWriter.WriteBlob(byteValue)
			require.NoError(t, err)
			err = binaryWriter.WriteBlob(byteValue)
			require.NoError(t, err)
		}

	case ion.ClobType:
		if reader.FieldName() != "" {
			err := textWriter.FieldName(reader.FieldName())
			require.NoError(t, err)
			err = binaryWriter.FieldName(reader.FieldName())
			require.NoError(t, err)
		}
		if reader.IsNull() {
			err := textWriter.WriteNullType(ion.ClobType)
			require.NoError(t, err)
			err = binaryWriter.WriteNullType(ion.ClobType)
			require.NoError(t, err)
		} else {
			byteValue, err := reader.ByteValue()
			require.NoError(t, err)
			err = textWriter.WriteClob(byteValue)
			require.NoError(t, err)
			err = binaryWriter.WriteClob(byteValue)
			require.NoError(t, err)
		}

	case ion.DecimalType:
		if reader.FieldName() != "" {
			err := textWriter.FieldName(reader.FieldName())
			require.NoError(t, err)
			err = binaryWriter.FieldName(reader.FieldName())
			require.NoError(t, err)
		}
		if reader.IsNull() {
			err := textWriter.WriteNullType(ion.DecimalType)
			require.NoError(t, err)
			err = binaryWriter.WriteNullType(ion.DecimalType)
			require.NoError(t, err)
		} else {
			decimalValue, err := reader.DecimalValue()
			require.NoError(t, err)
			err = textWriter.WriteDecimal(decimalValue)
			require.NoError(t, err)
			err = binaryWriter.WriteDecimal(decimalValue)
			require.NoError(t, err)
		}

	case ion.FloatType:
		if reader.FieldName() != "" {
			err := textWriter.FieldName(reader.FieldName())
			require.NoError(t, err)
			err = binaryWriter.FieldName(reader.FieldName())
			require.NoError(t, err)
		}
		if reader.IsNull() {
			err := textWriter.WriteNullType(ion.FloatType)
			require.NoError(t, err)
			err = binaryWriter.WriteNullType(ion.FloatType)
			require.NoError(t, err)
		} else {
			floatValue, err := reader.FloatValue()
			require.NoError(t, err)
			err = textWriter.WriteFloat(floatValue)
			require.NoError(t, err)
			err = binaryWriter.WriteFloat(floatValue)
			require.NoError(t, err)
		}

	case ion.IntType:
		intSize, err := reader.IntSize()
		require.NoError(t, err)

		//if reader.FieldName() != "" {
		//	err := textWriter.FieldName(reader.FieldName())
		//	require.NoError(t, err)
		//	err = binaryWriter.FieldName(reader.FieldName())
		//	require.NoError(t, err)
		//}
		if reader.IsNull() {
			err = textWriter.WriteNullType(ion.IntType)
			require.NoError(t, err)
			err = binaryWriter.WriteNullType(ion.IntType)
			require.NoError(t, err)
		} else {
			switch intSize {
			case ion.Int32:

				intValue, err := reader.IntValue()
				require.NoError(t, err)
				err = textWriter.WriteInt(int64(intValue))
				require.NoError(t, err)
				err = binaryWriter.WriteInt(int64(intValue))
				require.NoError(t, err)

			case ion.Int64:
				intValue, err := reader.Int64Value()
				require.NoError(t, err)
				err = textWriter.WriteInt(intValue)
				require.NoError(t, err)
				err = binaryWriter.WriteInt(intValue)
				require.NoError(t, err)

			case ion.Uint64:
				intValue, err := reader.Uint64Value()
				require.NoError(t, err)
				err = textWriter.WriteUint(intValue)
				require.NoError(t, err)
				err = binaryWriter.WriteUint(intValue)
				require.NoError(t, err)

			case ion.BigInt:
				intValue, err := reader.BigIntValue()
				require.NoError(t, err)
				err = textWriter.WriteBigInt(intValue)
				require.NoError(t, err)
				err = binaryWriter.WriteBigInt(intValue)
				require.NoError(t, err)

			default:
				t.Error("Expected intSize to be one of Int32, Int64, Uint64, or BigInt")
			}
		}

	case ion.StringType:
		if reader.FieldName() != "" {
			err := textWriter.FieldName(reader.FieldName())
			require.NoError(t, err)
			err = binaryWriter.FieldName(reader.FieldName())
			require.NoError(t, err)
		}
		if reader.IsNull() {
			err := textWriter.WriteNullType(ion.StringType)
			require.NoError(t, err)
			err = binaryWriter.WriteNullType(ion.StringType)
			require.NoError(t, err)
		} else {
			stringValue, err := reader.StringValue()
			require.NoError(t, err)
			err = textWriter.WriteString(stringValue)
			require.NoError(t, err)
			err = binaryWriter.WriteString(stringValue)
			require.NoError(t, err)
		}

	case ion.SymbolType:
		if reader.FieldName() != "" {
			err := textWriter.FieldName(reader.FieldName())
			require.NoError(t, err)
			err = binaryWriter.FieldName(reader.FieldName())
			require.NoError(t, err)
		}
		if reader.IsNull() {
			err := textWriter.WriteNullType(ion.SymbolType)
			require.NoError(t, err)
			err = binaryWriter.WriteNullType(ion.SymbolType)
			require.NoError(t, err)
		} else {
			stringValue, err := reader.StringValue()
			require.NoError(t, err)
			err = textWriter.WriteSymbol(stringValue)
			require.NoError(t, err)
			err = binaryWriter.WriteSymbol(stringValue)
			require.NoError(t, err)
		}

	case ion.TimestampType:
		if reader.FieldName() != "" {
			err := textWriter.FieldName(reader.FieldName())
			require.NoError(t, err)
			err = binaryWriter.FieldName(reader.FieldName())
			require.NoError(t, err)
		}
		if reader.IsNull() {
			err := textWriter.WriteNullType(ion.TimestampType)
			require.NoError(t, err)
			err = binaryWriter.WriteNullType(ion.TimestampType)
			require.NoError(t, err)
		} else {
			timeValue, err := reader.TimeValue()
			require.NoError(t, err)
			err = textWriter.WriteTimestamp(timeValue)
			require.NoError(t, err)
			err = binaryWriter.WriteTimestamp(timeValue)
			require.NoError(t, err)
		}

	case ion.SexpType:
		if reader.IsNull() {
			err := textWriter.WriteNullType(ion.SexpType)
			require.NoError(t, err)
			err = binaryWriter.WriteNullType(ion.SexpType)
			require.NoError(t, err)
		} else {
			err := reader.StepIn()
			require.NoError(t, err)
			err = textWriter.BeginSexp()
			require.NoError(t, err)
			err = binaryWriter.BeginSexp()
			require.NoError(t, err)
			if reader.FieldName() != "" {
				err = textWriter.FieldName(reader.FieldName())
				require.NoError(t, err)
				err = binaryWriter.FieldName(reader.FieldName())
				require.NoError(t, err)
			}
			for reader.Next() {
				writeToWriter(t, reader, textWriter, binaryWriter)
			}

			err = reader.StepOut()
			require.NoError(t, err)
			err = textWriter.EndSexp()
			require.NoError(t, err)
			err = binaryWriter.EndSexp()
			require.NoError(t, err)
		}

	case ion.ListType:
		if reader.IsNull() {
			err := textWriter.WriteNullType(ion.ListType)
			require.NoError(t, err)
			err = binaryWriter.WriteNullType(ion.ListType)
			require.NoError(t, err)
		} else {
			err := reader.StepIn()
			require.NoError(t, err)
			err = textWriter.BeginList()
			require.NoError(t, err)
			err = binaryWriter.BeginList()
			require.NoError(t, err)
			if reader.FieldName() != "" {
				err = textWriter.FieldName(reader.FieldName())
				require.NoError(t, err)
				err = binaryWriter.FieldName(reader.FieldName())
				require.NoError(t, err)
			}
			for reader.Next() {
				writeToWriter(t, reader, textWriter, binaryWriter)
			}

			err = reader.StepOut()
			require.NoError(t, err)
			err = textWriter.EndList()
			require.NoError(t, err)
			err = binaryWriter.EndList()
			require.NoError(t, err)
		}

	case ion.StructType:
		if reader.IsNull() {
			err := textWriter.WriteNullType(ion.StructType)
			require.NoError(t, err)
			err = binaryWriter.WriteNullType(ion.StructType)
			require.NoError(t, err)
		} else {
			err := reader.StepIn()
			require.NoError(t, err)
			err = textWriter.BeginStruct()
			require.NoError(t, err)
			err = binaryWriter.BeginStruct()
			require.NoError(t, err)
			if reader.FieldName() != "" {
				err = textWriter.FieldName(reader.FieldName())
				require.NoError(t, err)
				err = binaryWriter.FieldName(reader.FieldName())
				require.NoError(t, err)
			}

			for reader.Next() {
				writeToWriter(t, reader, textWriter, binaryWriter)
			}

			err = reader.StepOut()
			require.NoError(t, err)
			err = textWriter.EndStruct()
			require.NoError(t, err)
			err = binaryWriter.EndStruct()
			require.NoError(t, err)
		}

	default:
		t.Fatal(InvalidIonTypeError{ionType})
	}
}

func readSexpAndAppendToList(t *testing.T, reader ion.Reader) []byte {
	err := reader.StepIn()
	require.NoError(t, err)
	updateBytes := []byte{}
	for reader.Next() {
		intValue, err := reader.Int64Value()
		assert.NoError(t, err, "Something went wrong executing reader.Int64Value()")
		updateBytes = append(updateBytes, byte(intValue))
	}
	err = reader.StepOut()
	require.NoError(t, err)
	return updateBytes
}

type hashLog struct {
	identityUpdateList  [][]byte
	identityDigestList  [][]byte
	identityFinalDigest [][]byte
	md5UpdateList       [][]byte
	md5DigestList       [][]byte
}
