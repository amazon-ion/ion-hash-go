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
	"io/ioutil"
	"testing"

	"github.com/amzn/ion-go/ion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptyString(t *testing.T) {
	ionHashReader, err := NewHashReader(ion.NewReaderStr(""), newIdentityHasherProvider())
	require.NoError(t, err, "Expected NewHashReader() to successfully create a HashReader")

	for i := 0; i < 2; i++ {
		assert.False(t, ionHashReader.Next())
		assert.NoError(t, ionHashReader.Err(), "Something went wrong executing ionHashReader.Next()")

		assert.Equal(t, ion.NoType.String(), ionHashReader.Type().String(), "ionHashReader.Type() was not as expected")

		sum, err := ionHashReader.Sum(nil)
		require.NoError(t, err, "Something went wrong executing ionHashReader.Sum(nil)")

		assert.Equal(t, []byte{}, sum, "sum did not match expectation")
	}
}

func TestTopLevelValues(t *testing.T) {
	ionHashReader, err := NewHashReader(ion.NewReaderStr("1 2 3"), newIdentityHasherProvider())
	require.NoError(t, err, "Expected NewHashReader() to successfully create a HashReader")

	expectedTypes := []ion.Type{ion.IntType, ion.IntType, ion.IntType, ion.NoType, ion.NoType}
	expectedSums := [][]byte{[]byte{}, []byte{0x0b, 0x20, 0x01, 0x0e}, []byte{0x0b, 0x20, 0x02, 0x0e},
		[]byte{0x0b, 0x20, 0x03, 0x0e}, []byte{}}

	for i, expectedType := range expectedTypes {
		if expectedType == ion.NoType {
			assert.False(t, ionHashReader.Next())
			assert.NoError(t, ionHashReader.Err(), "Something went wrong executing ionHashReader.Next()")
		} else {
			assert.True(t, ionHashReader.Next())
		}

		assert.Equal(t, expectedType.String(), ionHashReader.Type().String(), "ionHashReader.Type() was not as expected")

		sum, err := ionHashReader.Sum(nil)
		require.NoError(t, err, "Something went wrong executing ionHashReader.Sum(nil)")

		assert.Equal(t, expectedSums[i], sum, "sum did not match expectation")
	}
}

func TestConsumeRemainderPartialConsume(t *testing.T) {
	consume(t, ConsumeRemainderPartialConsume)
}

func TestConsumeRemainderStepInStepOutNested(t *testing.T) {
	consume(t, ConsumeRemainderStepInStepOutNested)
}

func TestConsumeRemainderStepInNextStepOut(t *testing.T) {
	consume(t, ConsumeRemainderStepInNextStepOut)
}

func TestConsumeRemainderStepInStepOutTopLevel(t *testing.T) {
	consume(t, ConsumeRemainderStepInStepOutTopLevel)
}

func TestConsumeRemainderNext(t *testing.T) {
	consume(t, ConsumeRemainderNext)
}

func TestReaderUnresolvedSid(t *testing.T) {
	t.Skip() // Skipping test until SymbolToken is implemented

	reader := ion.NewReaderBytes([]byte{0xd3, 0x8a, 0x21, 0x01})

	ionHashReader, err := NewHashReader(reader, newIdentityHasherProvider())
	require.NoError(t, err, "Expected NewHashReader() to successfully create a HashReader")

	assert.False(t, ionHashReader.Next())

	require.False(t, ionHashReader.Next())
	require.Error(t, ionHashReader.Err())
	assert.IsType(t, &UnknownSymbolError{}, ionHashReader.Err())
}

func TestIonReaderContract(t *testing.T) {
	t.Skip()

	file, err := ioutil.ReadFile("ion-hash-test/ion_hash_tests.ion")
	require.NoError(t, err, "Something went wrong loading ion_hash_tests.ion")

	reader := ion.NewReaderBytes(file)

	ionHashReader, err := NewHashReader(reader, newIdentityHasherProvider())
	require.NoError(t, err, "Expected NewHashReader() to successfully create a HashReader")

	compareReaders(t, reader, ionHashReader)
}

func ConsumeRemainderPartialConsume(t *testing.T, ionHashReader HashReader) {
	ionHashReader.Next()
	assert.NoError(t, ionHashReader.StepIn(), "Something went wrong executing ionHashReader.StepIn()")

	ionHashReader.Next()
	ionHashReader.Next()
	ionHashReader.Next()
	assert.NoError(t, ionHashReader.StepIn(), "Something went wrong executing ionHashReader.StepIn()")

	ionHashReader.Next()
	// we've only partially consumed the struct
	assert.NoError(t, ionHashReader.StepOut(), "Something went wrong executing ionHashReader.StepOut()")

	// we've only partially consumed the list
	assert.NoError(t, ionHashReader.StepOut(), "Something went wrong executing ionHashReader.StepOut()")
}

func ConsumeRemainderStepInStepOutNested(t *testing.T, ionHashReader HashReader) {
	ionHashReader.Next()
	assert.NoError(t, ionHashReader.StepIn(), "Something went wrong executing ionHashReader.StepIn()")

	ionHashReader.Next()
	ionHashReader.Next()
	ionHashReader.Next()
	assert.NoError(t, ionHashReader.StepIn(), "Something went wrong executing ionHashReader.StepIn()")

	// we haven't consumed ANY of the struct
	assert.NoError(t, ionHashReader.StepOut(), "Something went wrong executing ionHashReader.StepOut()")

	// we've only partially consumed the list
	assert.NoError(t, ionHashReader.StepOut(), "Something went wrong executing ionHashReader.StepOut()")
}

func ConsumeRemainderStepInNextStepOut(t *testing.T, ionHashReader HashReader) {
	ionHashReader.Next()
	assert.NoError(t, ionHashReader.StepIn(), "Something went wrong executing ionHashReader.StepIn()")

	ionHashReader.Next()
	// we've partially consumed the list
	assert.NoError(t, ionHashReader.StepOut(), "Something went wrong executing ionHashReader.StepOut()")
}

func ConsumeRemainderStepInStepOutTopLevel(t *testing.T, ionHashReader HashReader) {
	ionHashReader.Next()
	sum, err := ionHashReader.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashReader.Sum(nil)")

	assert.Equal(t, []byte{}, sum, "sum did not match expectation")

	assert.NoError(t, ionHashReader.StepIn(), "Something went wrong executing ionHashReader.StepIn()")

	_, err = ionHashReader.Sum(nil)
	assert.Error(t, err, "Expected ionHashReader.Sum(nil) to return an error")
	assert.IsType(t, &InvalidOperationError{}, err)

	// we haven't consumed ANY of the list
	assert.NoError(t, ionHashReader.StepOut(), "Something went wrong executing ionHashReader.StepOut()")
}

func ConsumeRemainderNext(t *testing.T, ionHashReader HashReader) {
	assert.True(t, ionHashReader.Next())
	assert.NoError(t, ionHashReader.Err(), "Something went wrong executing ionHashReader.Next()")

	assert.False(t, ionHashReader.Next())
	assert.NoError(t, ionHashReader.Err(), "Something went wrong executing ionHashReader.Next()")
}

type consumeFunction func(*testing.T, HashReader)

func consume(t *testing.T, function consumeFunction) {
	ionHashReader, err := NewHashReader(ion.NewReaderStr("[1,2,{a:3,b:4},5]"), newIdentityHasherProvider())
	require.NoError(t, err, "Expected NewHashReader() to successfully create a HashReader")

	sum, err := ionHashReader.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashReader.Sum(nil)")

	assert.Equal(t, []byte{}, sum, "sum did not match expectation")

	function(t, ionHashReader)

	sum, err = ionHashReader.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashReader.Sum(nil)")

	expectedSum := []byte{0x0b, 0xb0, 0x0b, 0x20, 0x01, 0x0e, 0x0b, 0x20, 0x02, 0x0e, 0x0b, 0xd0, 0x0c, 0x0b, 0x70,
		0x61, 0x0c, 0x0e, 0x0c, 0x0b, 0x20, 0x03, 0x0c, 0x0e, 0x0c, 0x0b, 0x70, 0x62, 0x0c, 0x0e, 0x0c, 0x0b, 0x20,
		0x04, 0x0c, 0x0e, 0x0e, 0x0b, 0x20, 0x05, 0x0e, 0x0e}

	assert.Equal(t, expectedSum, sum, "sum did not match expectation")

	assert.False(t, ionHashReader.Next())
	assert.NoError(t, ionHashReader.Err(), "Something went wrong executing ionHashReader.Next()")

	assert.Equal(t, ion.NoType.String(), ionHashReader.Type().String(), "ionHashReader.Type() was not as expected")

	sum, err = ionHashReader.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashReader.Sum(nil)")

	assert.Equal(t, []byte{}, sum, "sum did not match expectation")
}
