package ionhash

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/amzn/ion-go/ion"
)

func TestEmptyString(t *testing.T) {
	ionHashReader, err := NewHashReader(ion.NewReaderStr(""), NewCryptoHasherProvider(SHA256))
	if err != nil {
		t.Fatal(err)
	}

	if !ionHashReader.Next() {
		t.Error("expected ionHashReader.Next() to return true")
	}

	ionType := ionHashReader.Type()
	if ionType != ion.NoType {
		t.Errorf("expected ionHashReader.Type() to return ion.NoType rather than %s", ionType.String())
	}

	sum, err := ionHashReader.Sum(nil)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(sum, []byte{}) {
		t.Error("sums don't match")
	}
}

func TestTopLevelValues(t *testing.T) {
	ionHashReader, err := NewHashReader(ion.NewReaderStr("1 2 3"), NewCryptoHasherProvider(SHA256))
	if err != nil {
		t.Fatal(err)
	}

	expectedTypes := []ion.Type{ion.IntType, ion.IntType, ion.IntType, ion.NoType, ion.NoType}
	expectedSums := [][]byte{[]byte{}, []byte{0x0b, 0x20, 0x01, 0x0e}, []byte{0x0b, 0x20, 0x02, 0x0e},
		[]byte{0x0b, 0x20, 0x03, 0x0e}, []byte{}}

	for i, expectedType := range expectedTypes {
		if !ionHashReader.Next() {
			t.Error("expected ionHashReader.Next() to return true")
		}

		ionType := ionHashReader.Type()
		if ionType != expectedType {
			t.Errorf("expected ionHashReader.Type() to return %s rather than %s",
				expectedType.String(), ionType.String())
		}

		sum, err := ionHashReader.Sum(nil)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(sum, expectedSums[i]) {
			t.Error("sums don't match")
		}
	}
}

func TestConsumeRemainderPartialConsume(t *testing.T) {
	err := consume(ConsumeRemainderPartialConsume)
	if err != nil {
		t.Error(err)
	}
}

func TestConsumeRemainderStepInStepOutNested(t *testing.T) {
	err := consume(ConsumeRemainderStepInStepOutNested)
	if err != nil {
		t.Error(err)
	}
}

func TestConsumeRemainderStepInNextStepOut(t *testing.T) {
	err := consume(ConsumeRemainderStepInNextStepOut)
	if err != nil {
		t.Error(err)
	}
}

func TestConsumeRemainderStepInStepOutTopLevel(t *testing.T) {
	err := consume(ConsumeRemainderStepInStepOutTopLevel)
	if err != nil {
		t.Error(err)
	}
}

func TestConsumeRemainderSingleNext(t *testing.T) {
	err := consume(ConsumeRemainderSingleNext)
	if err != nil {
		t.Error(err)
	}
}

func TestUnresolvedSid(t *testing.T) {
	ionReader := ion.NewReaderBytes([]byte{0xd3, 0x8a, 0x21, 0x01})

	ionHashReader, err := NewHashReader(ionReader, NewCryptoHasherProvider(SHA256))
	if err != nil {
		t.Error(err)
	}

	if ionHashReader.Next() {
		t.Error("expected ionHashReader.Next() to return false")
	} else {
		err := ionHashReader.Err()
		_, ok := err.(*UnknownSymbolError)
		if !ok {
			t.Error("expected ionHashReader.Next() to result in an UnknownSymbolError")
		}
	}
}

func TestIonReaderContract(t *testing.T) {
	file, err := ioutil.ReadFile("ion_hash_tests.ion")
	if err != nil {
		t.Fatal(err)
	}

	ionReader := ion.NewReaderBytes(file)

	ionHashReader, err := NewHashReader(ionReader, NewCryptoHasherProvider(SHA256))
	if err != nil {
		t.Fatal(err)
	}

	err = Compare(ionReader, ionHashReader)
	if err != nil {
		t.Error(err)
	}
}

func ConsumeRemainderPartialConsume(ionHashReader HashReader) error {
	ionHashReader.Next()
	err := ionHashReader.StepIn()
	if err != nil {
		return err
	}

	ionHashReader.Next()
	ionHashReader.Next()
	ionHashReader.Next()
	err = ionHashReader.StepIn()
	if err != nil {
		return err
	}

	ionHashReader.Next()
	err = ionHashReader.StepOut() // we've only partially consumed the struct
	if err != nil {
		return err
	}

	err = ionHashReader.StepOut() // we've only partially consumed the list
	if err != nil {
		return err
	}

	return nil
}

func ConsumeRemainderStepInStepOutNested(ionHashReader HashReader) error {
	ionHashReader.Next()
	err := ionHashReader.StepIn()
	if err != nil {
		return err
	}

	ionHashReader.Next()
	ionHashReader.Next()
	ionHashReader.Next()
	err = ionHashReader.StepIn()
	if err != nil {
		return err
	}

	err = ionHashReader.StepIn()
	if err != nil {
		return err
	}

	err = ionHashReader.StepOut() // we haven't consumed ANY of the struct
	if err != nil {
		return err
	}

	err = ionHashReader.StepOut() // we've only partially consumed the list
	if err != nil {
		return err
	}

	return nil
}

func ConsumeRemainderStepInNextStepOut(ionHashReader HashReader) error {
	ionHashReader.Next()
	err := ionHashReader.StepIn()
	if err != nil {
		return err
	}

	ionHashReader.Next()
	err = ionHashReader.StepOut() // we've partially consumed the list
	if err != nil {
		return err
	}

	return nil
}

func ConsumeRemainderStepInStepOutTopLevel(ionHashReader HashReader) error {
	ionHashReader.Next()
	sum, err := ionHashReader.Sum(nil)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(sum, []byte{}) {
		return fmt.Errorf("sums don't match'")
	}

	err = ionHashReader.StepIn()
	if err != nil {
		return err
	}

	_, err = ionHashReader.Sum(nil)
	if err != nil {
		_, ok := err.(*InvalidOperationError)
		if !ok {
			return fmt.Errorf("expected ionHashReader.Sum(nil) to return an InvalidOperationError")
		}
	} else {
		return fmt.Errorf("expected ionHashReader.Sum(nil) to return an error")
	}

	err = ionHashReader.StepOut() // we haven't consumed ANY of the list
	if err != nil {
		return err
	}

	return nil
}

func ConsumeRemainderSingleNext(ionHashReader HashReader) error {
	ionHashReader.Next()
	ionHashReader.Next()

	return nil
}

type consumeFunction func(HashReader) error

func consume(function consumeFunction) error {
	ionHashReader, err := NewHashReader(ion.NewReaderStr("[1,2,{a:3,b:4},5]"), NewCryptoHasherProvider(SHA256))
	if err != nil {
		return err
	}

	sum, err := ionHashReader.Sum(nil)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(sum, []byte{}) {
		return fmt.Errorf("sums don't match")
	}

	err = function(ionHashReader)
	if err != nil {
		return err
	}

	sum, err = ionHashReader.Sum(nil)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(sum, []byte{0x0b, 0xb0, 0x0b, 0x20, 0x01, 0x0e, 0x0b, 0x20, 0x02, 0x0e,
		0x0b, 0xd0, 0x0c, 0x0b, 0x70, 0x61, 0x0c, 0x0e, 0x0c, 0x0b, 0x20, 0x03, 0x0c, 0x0e, 0x0c,
		0x0b, 0x70, 0x62, 0x0c, 0x0e, 0x0c, 0x0b, 0x20, 0x04, 0x0c, 0x0e, 0x0e, 0x0b, 0x20, 0x05,
		0x0e, 0x0e}) {
		return fmt.Errorf("sums don't match")
	}

	if !ionHashReader.Next() {
		return fmt.Errorf("expected ionHashReader.Next() to return true")
	}

	ionType := ionHashReader.Type()
	if ionType != ion.NoType {
		return fmt.Errorf("expected ionHashReader.Type() to return ion.NoType rather than %s", ionType.String())
	}

	sum, err = ionHashReader.Sum(nil)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(sum, []byte{}) {
		return fmt.Errorf("sums don't match")
	}

	return nil
}
