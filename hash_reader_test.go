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
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/amzn/ion-go/ion"
	"github.com/amzn/ion-hash-go/internal"
)

func TestEmptyString(t *testing.T) {
	ionHashReader, err := NewHashReader(ion.NewReaderStr(""), newIdentityHasherProvider())
	if err != nil {
		t.Fatalf("expected NewHashReader() to successfully create a HashReader; %s", err.Error())
	}

	for i := 0; i < 2; i++ {
		if !ionHashReader.Next() {
			err = ionHashReader.Err()
			if err != nil {
				t.Errorf("expected ionHashReader.Next() to return false without errors; %s", err.Error())
			}
		}

		ionType := ionHashReader.Type()
		if ionType != ion.NoType {
			t.Errorf("expected ionHashReader.Type() to return ion.NoType rather than %s", ionType.String())
		}

		sum, err := ionHashReader.Sum(nil)
		if err != nil {
			t.Fatalf("expected Sum(nil) to execute without errors; %s", err.Error())
		}

		if !reflect.DeepEqual(sum, []byte{}) {
			t.Errorf("expected sum to be %v instead of %v", []byte{}, sum)
		}
	}
}

func TestTopLevelValues(t *testing.T) {
	ionHashReader, err := NewHashReader(ion.NewReaderStr("1 2 3"), newIdentityHasherProvider())
	if err != nil {
		t.Fatalf("expected NewHashReader() to successfully create a HashReader; %s", err.Error())
	}

	expectedTypes := []ion.Type{ion.IntType, ion.IntType, ion.IntType, ion.NoType, ion.NoType}
	expectedSums := [][]byte{[]byte{}, []byte{0x0b, 0x20, 0x01, 0x0e}, []byte{0x0b, 0x20, 0x02, 0x0e},
		[]byte{0x0b, 0x20, 0x03, 0x0e}, []byte{}}

	for i, expectedType := range expectedTypes {
		if !ionHashReader.Next() {
			err = ionHashReader.Err()
			if err != nil {
				t.Errorf("expected ionHashReader.Next() to return true; %s", err.Error())
			} else if expectedType != ion.NoType {
				t.Errorf("expected ionHashReader.Next() to return true")
			}
		}

		ionType := ionHashReader.Type()
		if ionType != expectedType {
			t.Errorf("expected ionHashReader.Type() to return %s rather than %s",
				expectedType.String(), ionType.String())
		}

		sum, err := ionHashReader.Sum(nil)
		if err != nil {
			t.Fatalf("expected Sum(nil) to execute without errors; %s", err.Error())
		}

		if !reflect.DeepEqual(sum, expectedSums[i]) {
			t.Errorf("expected sum to be %v instead of %v", expectedSums[i], sum)
		}
	}
}

func TestConsumeRemainderPartialConsume(t *testing.T) {
	t.Skip() // Skipping test until reader's IsInStruct logic is updated to match dot net

	err := consume(ConsumeRemainderPartialConsume)
	if err != nil {
		t.Error(err)
	}
}

func TestConsumeRemainderStepInStepOutNested(t *testing.T) {
	t.Skip() // Skipping test until reader's IsInStruct logic is updated to match dot net

	err := consume(ConsumeRemainderStepInStepOutNested)
	if err != nil {
		t.Error(err)
	}
}

func TestConsumeRemainderStepInNextStepOut(t *testing.T) {
	t.Skip() // Skipping test until reader's IsInStruct logic is updated to match dot net

	err := consume(ConsumeRemainderStepInNextStepOut)
	if err != nil {
		t.Error(err)
	}
}

func TestConsumeRemainderStepInStepOutTopLevel(t *testing.T) {
	t.Skip() // Skipping test until reader's IsInStruct logic is updated to match dot net

	err := consume(ConsumeRemainderStepInStepOutTopLevel)
	if err != nil {
		t.Error(err)
	}
}

func TestConsumeRemainderNext(t *testing.T) {
	t.Skip() // Skipping test until reader's IsInStruct logic is updated to match dot net

	err := consume(ConsumeRemainderNext)
	if err != nil {
		t.Error(err)
	}
}

func TestReaderUnresolvedSid(t *testing.T) {
	t.Skip() // Skipping test until SymbolToken is implemented

	ionReader := ion.NewReaderBytes([]byte{0xd3, 0x8a, 0x21, 0x01})

	ionHashReader, err := NewHashReader(ionReader, newIdentityHasherProvider())
	if err != nil {
		t.Fatalf("expected NewHashReader() to successfully create a HashReader; %s", err.Error())
	}

	if ionHashReader.Next() {
		t.Error("expected ionHashReader.Next() to return false")
	}

	if ionHashReader.Next() {
		t.Error("expected ionHashReader.Next() to return false")
	} else {
		err := ionHashReader.Err()
		_, ok := err.(*internal.UnknownSymbolError)
		if !ok {
			t.Error("expected ionHashReader.Next() to result in an UnknownSymbolError")
		}
	}
}

func TestIonReaderContract(t *testing.T) {
	t.Skip()

	file, err := ioutil.ReadFile("ion-hash-test/ion_hash_tests.ion")
	if err != nil {
		t.Fatal("expected ion_hash_tests.ion to load properly")
	}

	ionReader := ion.NewReaderBytes(file)

	ionHashReader, err := NewHashReader(ionReader, newIdentityHasherProvider())
	if err != nil {
		t.Fatalf("expected NewHashReader() to successfully create a HashReader; %s", err.Error())
	}

	compare, err := compareReaders(ionReader, ionHashReader)
	if !compare {
		if err != nil {
			t.Errorf("expected compareReaders(ionReader, ionHashReader) to return true; %s", err.Error())
		} else {
			t.Errorf("expected compareReaders(ionReader, ionHashReader) to return true")
		}
	}
}

func ConsumeRemainderPartialConsume(ionHashReader HashReader) error {
	ionHashReader.Next()
	err := ionHashReader.StepIn()
	if err != nil {
		return fmt.Errorf("expected ionHashReader.StepIn() to successfully run without errors; %s", err.Error())
	}

	ionHashReader.Next()
	ionHashReader.Next()
	ionHashReader.Next()
	err = ionHashReader.StepIn()
	if err != nil {
		return fmt.Errorf("expected ionHashReader.StepIn() to successfully run without errors; %s", err.Error())
	}

	ionHashReader.Next()
	err = ionHashReader.StepOut() // we've only partially consumed the struct
	if err != nil {
		return fmt.Errorf("expected ionHashReader.StepOut() to successfully run without errors; %s", err.Error())
	}

	err = ionHashReader.StepOut() // we've only partially consumed the list
	if err != nil {
		return fmt.Errorf("expected ionHashReader.StepOut() to successfully run without errors; %s", err.Error())
	}

	return nil
}

func ConsumeRemainderStepInStepOutNested(ionHashReader HashReader) error {
	ionHashReader.Next()
	err := ionHashReader.StepIn()
	if err != nil {
		return fmt.Errorf("expected ionHashReader.StepIn() to successfully run without errors; %s", err.Error())
	}

	ionHashReader.Next()
	ionHashReader.Next()
	ionHashReader.Next()
	err = ionHashReader.StepIn()
	if err != nil {
		return fmt.Errorf("expected ionHashReader.StepIn() to successfully run without errors; %s", err.Error())
	}

	err = ionHashReader.StepOut() // we haven't consumed ANY of the struct
	if err != nil {
		return fmt.Errorf("expected ionHashReader.StepOut() to successfully run without errors; %s", err.Error())
	}

	err = ionHashReader.StepOut() // we've only partially consumed the list
	if err != nil {
		return fmt.Errorf("expected ionHashReader.StepOut() to successfully run without errors; %s", err.Error())
	}

	return nil
}

func ConsumeRemainderStepInNextStepOut(ionHashReader HashReader) error {
	ionHashReader.Next()
	err := ionHashReader.StepIn()
	if err != nil {
		return fmt.Errorf("expected ionHashReader.StepIn() to successfully run without errors; %s", err.Error())
	}

	ionHashReader.Next()
	err = ionHashReader.StepOut() // we've partially consumed the list
	if err != nil {
		return fmt.Errorf("expected ionHashReader.StepOut() to successfully run without errors; %s", err.Error())
	}

	return nil
}

func ConsumeRemainderStepInStepOutTopLevel(ionHashReader HashReader) error {
	ionHashReader.Next()
	sum, err := ionHashReader.Sum(nil)
	if err != nil {
		return fmt.Errorf("expected Sum(nil) to execute without errors; %s", err.Error())
	}

	if !reflect.DeepEqual(sum, []byte{}) {
		return fmt.Errorf("expected sum to be %v instead of %v", []byte{}, sum)
	}

	err = ionHashReader.StepIn()
	if err != nil {
		return fmt.Errorf("expected ionHashReader.StepIn() to successfully run without errors; %s", err.Error())
	}

	_, err = ionHashReader.Sum(nil)
	if err != nil {
		_, ok := err.(*internal.InvalidOperationError)
		if !ok {
			return fmt.Errorf("expected ionHashReader.Sum(nil) to return an InvalidOperationError")
		}
	} else {
		return fmt.Errorf("expected ionHashReader.Sum(nil) to return an error")
	}

	err = ionHashReader.StepOut() // we haven't consumed ANY of the list
	if err != nil {
		return fmt.Errorf("expected ionHashReader.StepOut() to successfully run without errors; %s", err.Error())
	}

	return nil
}

func ConsumeRemainderNext(ionHashReader HashReader) error {
	ionHashReader.Next()
	ionHashReader.Next()

	return nil
}

type consumeFunction func(HashReader) error

func consume(function consumeFunction) error {
	ionHashReader, err := NewHashReader(ion.NewReaderStr("[1,2,{a:3,b:4},5]"), newIdentityHasherProvider())
	if err != nil {
		return err
	}

	sum, err := ionHashReader.Sum(nil)
	if err != nil {
		return fmt.Errorf("expected Sum(nil) to execute without errors; %s", err.Error())
	}

	if !reflect.DeepEqual(sum, []byte{}) {
		return fmt.Errorf("expected sum to be %v instead of %v", []byte{}, sum)
	}

	err = function(ionHashReader)
	if err != nil {
		return err
	}

	sum, err = ionHashReader.Sum(nil)
	if err != nil {
		return fmt.Errorf("expected Sum(nil) to execute without errors; %s", err.Error())
	}

	expectedSum := []byte{0x0b, 0xb0, 0x0b, 0x20, 0x01, 0x0e, 0x0b, 0x20, 0x02, 0x0e, 0x0b, 0xd0, 0x0c, 0x0b, 0x70,
		0x61, 0x0c, 0x0e, 0x0c, 0x0b, 0x20, 0x03, 0x0c, 0x0e, 0x0c, 0x0b, 0x70, 0x62, 0x0c, 0x0e, 0x0c, 0x0b, 0x20,
		0x04, 0x0c, 0x0e, 0x0e, 0x0b, 0x20, 0x05, 0x0e, 0x0e}

	if !reflect.DeepEqual(expectedSum, sum) {
		return fmt.Errorf("expected sum to be %v instead of %v", expectedSum, sum)
	}

	if ionHashReader.Next() {
		return fmt.Errorf("expected ionHashReader.Next() to return false")
	}

	ionType := ionHashReader.Type()
	if ionType != ion.NoType {
		return fmt.Errorf("expected ionHashReader.Type() to return ion.NoType rather than %s", ionType.String())
	}

	sum, err = ionHashReader.Sum(nil)
	if err != nil {
		return fmt.Errorf("expected Sum(nil) to execute without errors; %s", err.Error())
	}

	if !reflect.DeepEqual(sum, []byte{}) {
		return fmt.Errorf("expected sum to be %v instead of %v", []byte{}, sum)
	}

	return nil
}
