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

package internal

import (
	"github.com/amzn/ion-go/ion"
	"github.com/amzn/ion-hash-go/ihp"
)

type scalarSerializer struct {
	baseSerializer
}

func newScalarSerializer(hashFunction ihp.IonHasher, depth int) serializer {
	return &scalarSerializer{baseSerializer{hashFunction: hashFunction, depth: depth}}
}

func (scalarSerializer *scalarSerializer) scalar(ionValue interface{}) error {
	ionHashValue := ionValue.(HashValue)

	err := scalarSerializer.handleAnnotationsBegin(ionHashValue)
	if err != nil {
		return err
	}

	err = scalarSerializer.beginMarker()
	if err != nil {
		return err
	}

	var ionVal interface{}
	var ionType ion.Type
	if ionHashValue.CurrentIsNull() {
		ionVal = nil
		ionType = ion.NoType
	} else {
		ionVal, err = ionHashValue.Value()
		if err != nil {
			return err
		}
		ionType = ionHashValue.IonType()
	}

	scalarBytes, err := scalarSerializer.getBytes(ionHashValue.IonType(), ionVal, ionHashValue.CurrentIsNull())
	if err != nil {
		return err
	}

	if ionHashValue.IonType() != ion.SymbolType {
		ionVal = nil
	}

	tq, representation, err :=
		scalarSerializer.scalarOrNullSplitParts(ionType, ionHashValue.CurrentIsNull(), scalarBytes)
	if err != nil {
		return err
	}

	err = scalarSerializer.write([]byte{tq})
	if err != nil {
		return err
	}

	if len(representation) > 0 {
		err = scalarSerializer.write(escape(representation))
		if err != nil {
			return err
		}
	}

	err = scalarSerializer.endMarker()
	if err != nil {
		return err
	}

	err = scalarSerializer.handleAnnotationsEnd(ionHashValue, false)
	if err != nil {
		return err
	}

	return nil
}

func (scalarSerializer *scalarSerializer) stepOut() error {
	return scalarSerializer.baseSerializer.stepOut()
}

func (scalarSerializer *scalarSerializer) stepIn(ionValue interface{}) error {
	return scalarSerializer.baseSerializer.stepIn(ionValue.(HashValue))
}

func (scalarSerializer *scalarSerializer) sum(b []byte) []byte {
	return scalarSerializer.baseSerializer.sum(b)
}

func (scalarSerializer *scalarSerializer) handleFieldName(ionValue interface{}) error {
	return scalarSerializer.baseSerializer.handleFieldName(ionValue.(HashValue))
}

func (scalarSerializer *scalarSerializer) write(bytes []byte) error {
	return scalarSerializer.baseSerializer.write(bytes)
}

func (scalarSerializer *scalarSerializer) beginMarker() error {
	return scalarSerializer.baseSerializer.beginMarker()
}

func (scalarSerializer *scalarSerializer) endMarker() error {
	return scalarSerializer.baseSerializer.endMarker()
}

func (scalarSerializer *scalarSerializer) handleAnnotationsBegin(ionValue interface{}) error {
	return scalarSerializer.baseSerializer.handleAnnotationsBegin(ionValue.(HashValue), false)
}

func (scalarSerializer *scalarSerializer) handleAnnotationsEnd(ionValue interface{}, isContainer bool) error {
	return scalarSerializer.baseSerializer.handleAnnotationsEnd(ionValue.(HashValue), isContainer)
}

func (scalarSerializer *scalarSerializer) writeSymbol(token string) error {
	return scalarSerializer.baseSerializer.writeSymbol(token)
}

func (scalarSerializer *scalarSerializer) getBytes(ionType ion.Type, ionValue interface{}, isNull bool) ([]byte, error) {
	return scalarSerializer.baseSerializer.getBytes(ionType, ionValue, isNull)
}

func (scalarSerializer *scalarSerializer) getLengthFieldLength(bytes []byte) (int, error) {
	return scalarSerializer.baseSerializer.getLengthFieldLength(bytes)
}

func (scalarSerializer *scalarSerializer) scalarOrNullSplitParts(
	ionType ion.Type, isNull bool, bytes []byte) (byte, []byte, error) {

	return scalarSerializer.baseSerializer.scalarOrNullSplitParts(ionType, isNull, bytes)
}
