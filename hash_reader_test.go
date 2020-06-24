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
	"reflect"
	"testing"

	"github.com/amzn/ion-go/ion"
)

func TestEmptyString(t *testing.T) {
	ionHashReader, err := NewHashReader(ion.NewReaderStr(""), newIdentityHasherProvider())
	if err != nil {
		t.Fatalf("Expected NewHashReader() to successfully create a HashReader; %s", err.Error())
	}

	for i := 0; i < 2; i++ {
		if !ionHashReader.Next() {
			err = ionHashReader.Err()
			if err != nil {
				t.Errorf("Something went wrong executing ionHashReader.Next(); %s", err.Error())
			}
		}

		ionType := ionHashReader.Type()
		if ionType != ion.NoType {
			t.Errorf("ionHashReader.Type() was not as expected;\n"+
				"Expected type: %s\n"+
				"Actual type:   %s",
				ion.NoType.String(), ionType.String())
		}

		sum, err := ionHashReader.Sum(nil)
		if err != nil {
			t.Fatalf("Something went wrong executing Sum(nil); %s", err.Error())
		}

		if !reflect.DeepEqual(sum, []byte{}) {
			t.Errorf("sum did not match expectation;\n"+
				"Expected sum: %v\n"+
				"Actual sum:   %v",
				[]byte{}, sum)
		}
	}
}

func TestTopLevelValues(t *testing.T) {
	ionHashReader, err := NewHashReader(ion.NewReaderStr("1 2 3"), newIdentityHasherProvider())
	if err != nil {
		t.Fatalf("Expected NewHashReader() to successfully create a HashReader; %s", err.Error())
	}

	expectedTypes := []ion.Type{ion.IntType, ion.IntType, ion.IntType, ion.NoType, ion.NoType}
	expectedSums := [][]byte{[]byte{}, []byte{0x0b, 0x20, 0x01, 0x0e}, []byte{0x0b, 0x20, 0x02, 0x0e},
		[]byte{0x0b, 0x20, 0x03, 0x0e}, []byte{}}

	for i, expectedType := range expectedTypes {
		if !ionHashReader.Next() {
			err = ionHashReader.Err()
			if err != nil {
				t.Errorf("Something went wrong executing ionHashReader.Next(); %s", err.Error())
			} else if expectedType != ion.NoType {
				t.Errorf("Expected ionHashReader.Next() to return true")
			}
		}

		ionType := ionHashReader.Type()
		if ionType != expectedType {
			t.Errorf("ionHashReader.Type() was not as expected;\n"+
				"Expected type: %s\n"+
				"Actual type:   %s",
				expectedType.String(), ionType.String())
		}

		sum, err := ionHashReader.Sum(nil)
		if err != nil {
			t.Fatalf("Something went wrong executing Sum(nil); %s", err.Error())
		}

		if !reflect.DeepEqual(sum, expectedSums[i]) {
			t.Errorf("sum did not match expectation;\n"+
				"Expected sum: %v\n"+
				"Actual sum:   %v",
				expectedSums[i], sum)
		}
	}
}

func TestConsumeRemainderPartialConsume(t *testing.T) {
	t.Skip() // Skipping test until reader's IsInStruct logic is updated to match dot net

	consume(t, ConsumeRemainderPartialConsume)
}

func TestConsumeRemainderStepInStepOutNested(t *testing.T) {
	t.Skip() // Skipping test until reader's IsInStruct logic is updated to match dot net

	consume(t, ConsumeRemainderStepInStepOutNested)
}

func TestConsumeRemainderStepInNextStepOut(t *testing.T) {
	t.Skip() // Skipping test until reader's IsInStruct logic is updated to match dot net

	consume(t, ConsumeRemainderStepInNextStepOut)
}

func TestConsumeRemainderStepInStepOutTopLevel(t *testing.T) {
	t.Skip() // Skipping test until reader's IsInStruct logic is updated to match dot net

	consume(t, ConsumeRemainderStepInStepOutTopLevel)
}

func TestConsumeRemainderNext(t *testing.T) {
	t.Skip() // Skipping test until reader's IsInStruct logic is updated to match dot net

	consume(t, ConsumeRemainderNext)
}

func TestReaderUnresolvedSid(t *testing.T) {
	t.Skip() // Skipping test until SymbolToken is implemented

	reader := ion.NewReaderBytes([]byte{0xd3, 0x8a, 0x21, 0x01})

	ionHashReader, err := NewHashReader(reader, newIdentityHasherProvider())
	if err != nil {
		t.Fatalf("Expected NewHashReader() to successfully create a HashReader; %s", err.Error())
	}

	if ionHashReader.Next() {
		t.Error("Expected ionHashReader.Next() to return false")
	}

	if ionHashReader.Next() {
		t.Error("Expected ionHashReader.Next() to return false")
	} else {
		err := ionHashReader.Err()
		_, ok := err.(*UnknownSymbolError)
		if !ok {
			t.Error("Expected ionHashReader.Next() to result in an UnknownSymbolError")
		}
	}
}

func TestIonReaderContract(t *testing.T) {
	t.Skip()

	file, err := ioutil.ReadFile("ion-hash-test/ion_hash_tests.ion")
	if err != nil {
		t.Fatalf("Something went wrong loading ion_hash_tests.ion; %s", err.Error())
	}

	reader := ion.NewReaderBytes(file)

	ionHashReader, err := NewHashReader(reader, newIdentityHasherProvider())
	if err != nil {
		t.Fatalf("Expected NewHashReader() to successfully create a HashReader; %s", err.Error())
	}

	compare, err := compareReaders(reader, ionHashReader)
	if !compare {
		if err != nil {
			t.Errorf("Something went wrong executing compareReaders(reader, ionHashReader); %s", err.Error())
		} else {
			t.Errorf("Expected compareReaders(reader, ionHashReader) to return true")
		}
	}
}

func ConsumeRemainderPartialConsume(t *testing.T, ionHashReader HashReader) {
	ionHashReader.Next()
	err := ionHashReader.StepIn()
	if err != nil {
		t.Errorf("Something went wrong executing ionHashReader.StepIn(); %s", err.Error())
	}

	ionHashReader.Next()
	ionHashReader.Next()
	ionHashReader.Next()
	err = ionHashReader.StepIn()
	if err != nil {
		t.Errorf("Something went wrong executing ionHashReader.StepIn(); %s", err.Error())
	}

	ionHashReader.Next()
	err = ionHashReader.StepOut() // we've only partially consumed the struct
	if err != nil {
		t.Errorf("Something went wrong executing ionHashReader.StepOut(); %s", err.Error())
	}

	err = ionHashReader.StepOut() // we've only partially consumed the list
	if err != nil {
		t.Errorf("Something went wrong executing ionHashReader.StepOut(); %s", err.Error())
	}
}

func ConsumeRemainderStepInStepOutNested(t *testing.T, ionHashReader HashReader) {
	ionHashReader.Next()
	err := ionHashReader.StepIn()
	if err != nil {
		t.Errorf("Something went wrong executing ionHashReader.StepIn(); %s", err.Error())
	}

	ionHashReader.Next()
	ionHashReader.Next()
	ionHashReader.Next()
	err = ionHashReader.StepIn()
	if err != nil {
		t.Errorf("Something went wrong executing ionHashReader.StepIn(); %s", err.Error())
	}

	err = ionHashReader.StepOut() // we haven't consumed ANY of the struct
	if err != nil {
		t.Errorf("Something went wrong executing ionHashReader.StepOut(); %s", err.Error())
	}

	err = ionHashReader.StepOut() // we've only partially consumed the list
	if err != nil {
		t.Errorf("Something went wrong executing ionHashReader.StepOut(); %s", err.Error())
	}
}

func ConsumeRemainderStepInNextStepOut(t *testing.T, ionHashReader HashReader) {
	ionHashReader.Next()
	err := ionHashReader.StepIn()
	if err != nil {
		t.Errorf("Something went wrong executing ionHashReader.StepIn(); %s", err.Error())
	}

	ionHashReader.Next()
	err = ionHashReader.StepOut() // we've partially consumed the list
	if err != nil {
		t.Errorf("Something went wrong executing ionHashReader.StepOut(); %s", err.Error())
	}
}

func ConsumeRemainderStepInStepOutTopLevel(t *testing.T, ionHashReader HashReader) {
	ionHashReader.Next()
	sum, err := ionHashReader.Sum(nil)
	if err != nil {
		t.Fatalf("Something went wrong executing Sum(nil); %s", err.Error())
	}

	if !reflect.DeepEqual(sum, []byte{}) {
		t.Errorf("sum did not match expectation;\n"+
			"Expected sum: %v\n"+
			"Actual sum:   %v",
			[]byte{}, sum)
	}

	err = ionHashReader.StepIn()
	if err != nil {
		t.Errorf("Something went wrong executing ionHashReader.StepIn(); %s", err.Error())
	}

	_, err = ionHashReader.Sum(nil)
	if err != nil {
		_, ok := err.(*InvalidOperationError)
		if !ok {
			t.Errorf("Expected ionHashReader.Sum(nil) to return an InvalidOperationError; %s", err.Error())
		}
	} else {
		t.Errorf("Expected ionHashReader.Sum(nil) to return an error")
	}

	err = ionHashReader.StepOut() // we haven't consumed ANY of the list
	if err != nil {
		t.Errorf("Something went wrong executing ionHashReader.StepOut(); %s", err.Error())
	}
}

func ConsumeRemainderNext(t *testing.T, ionHashReader HashReader) {
	if !ionHashReader.Next() {
		err := ionHashReader.Err()
		if err != nil {
			t.Errorf("Something went wrong executing ionHashReader.Next(); %s", err.Error())
		} else {
			t.Error("Expected ionHashReader.Next() to return true")
		}
	}

	if ionHashReader.Next() {
		t.Error("Expected ionHashReader.Next() to return false")
	} else {
		err := ionHashReader.Err()
		if err != nil {
			t.Errorf("Something went wrong executing ionHashReader.Next(); %s", err.Error())
		}
	}
}

type consumeFunction func(*testing.T, HashReader)

func consume(t *testing.T, function consumeFunction) {
	ionHashReader, err := NewHashReader(ion.NewReaderStr("[1,2,{a:3,b:4},5]"), newIdentityHasherProvider())
	if err != nil {
		t.Fatalf("Expected NewHashReader() to successfully create a HashReader; %s", err.Error())
	}

	sum, err := ionHashReader.Sum(nil)
	if err != nil {
		t.Fatalf("Something went wrong executing Sum(nil); %s", err.Error())
	}

	if !reflect.DeepEqual(sum, []byte{}) {
		t.Errorf("sum did not match expectation;\n"+
			"Expected sum: %v\n"+
			"Actual sum:   %v",
			[]byte{}, sum)
	}

	function(t, ionHashReader)

	sum, err = ionHashReader.Sum(nil)
	if err != nil {
		t.Fatalf("Something went wrong executing Sum(nil); %s", err.Error())
	}

	expectedSum := []byte{0x0b, 0xb0, 0x0b, 0x20, 0x01, 0x0e, 0x0b, 0x20, 0x02, 0x0e, 0x0b, 0xd0, 0x0c, 0x0b, 0x70,
		0x61, 0x0c, 0x0e, 0x0c, 0x0b, 0x20, 0x03, 0x0c, 0x0e, 0x0c, 0x0b, 0x70, 0x62, 0x0c, 0x0e, 0x0c, 0x0b, 0x20,
		0x04, 0x0c, 0x0e, 0x0e, 0x0b, 0x20, 0x05, 0x0e, 0x0e}

	if !reflect.DeepEqual(expectedSum, sum) {
		t.Errorf("sum did not match expectation;\n"+
			"Expected sum: %v\n"+
			"Actual sum:   %v",
			expectedSum, sum)
	}

	if ionHashReader.Next() {
		t.Errorf("Expected ionHashReader.Next() to return false")
	} else {
		err = ionHashReader.Err()
		if err != nil {
			t.Errorf("Something went wrong executing ionHashReader.Next(); %s", err.Error())
		}
	}

	ionType := ionHashReader.Type()
	if ionType != ion.NoType {
		t.Errorf("ionHashReader.Type() was not as expected;\n"+
			"Expected type: %s\n"+
			"Actual type:   %s",
			ion.NoType.String(), ionType.String())
	}

	sum, err = ionHashReader.Sum(nil)
	if err != nil {
		t.Fatalf("Something went wrong executing Sum(nil); %s", err.Error())
	}

	if !reflect.DeepEqual(sum, []byte{}) {
		t.Errorf("sum did not match expectation;\n"+
			"Expected sum: %v\n"+
			"Actual sum:   %v",
			[]byte{}, sum)
	}
}
