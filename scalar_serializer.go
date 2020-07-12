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

import "github.com/amzn/ion-go/ion"

type scalarSerializer struct {
	baseSerializer
}

func newScalarSerializer(hashFunction IonHasher, depth int) serializer {
	return &scalarSerializer{baseSerializer{hashFunction: hashFunction, depth: depth}}
}

func (scalarSerializer *scalarSerializer) scalar(ionValue interface{}) error {
	ionHashValue := ionValue.(hashValue)

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
	if ionHashValue.isNull() {
		ionVal = nil
		ionType = ion.NoType
	} else {
		ionVal, err = ionHashValue.value()
		if err != nil {
			return err
		}
		ionType = ionHashValue.ionType()
	}

	scalarBytes, err := scalarSerializer.getBytes(ionHashValue.ionType(), ionVal, ionHashValue.isNull())
	if err != nil {
		return err
	}

	if ionHashValue.ionType() != ion.SymbolType {
		ionVal = nil
	}

	tq, representation, err :=
		scalarSerializer.scalarOrNullSplitParts(ionType, ionHashValue.isNull(), scalarBytes)
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

func (scalarSerializer *scalarSerializer) stepIn(ionValue interface{}) error {
	return scalarSerializer.baseSerializer.stepIn(ionValue.(hashValue))
}

func (scalarSerializer *scalarSerializer) handleFieldName(ionValue interface{}) error {
	return scalarSerializer.baseSerializer.handleFieldName(ionValue.(hashValue))
}

func (scalarSerializer *scalarSerializer) handleAnnotationsBegin(ionValue interface{}) error {
	return scalarSerializer.baseSerializer.handleAnnotationsBegin(ionValue.(hashValue), false)
}

func (scalarSerializer *scalarSerializer) handleAnnotationsEnd(ionValue interface{}, isContainer bool) error {
	return scalarSerializer.baseSerializer.handleAnnotationsEnd(ionValue.(hashValue), isContainer)
}
