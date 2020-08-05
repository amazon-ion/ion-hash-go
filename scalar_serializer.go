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

func (ss *scalarSerializer) scalar(ionValue interface{}) error {
	ionHashValue := ionValue.(hashValue)

	err := ss.handleAnnotationsBegin(ionHashValue)
	if err != nil {
		return err
	}

	err = ss.beginMarker()
	if err != nil {
		return err
	}

	var ionVal interface{}
	var ionType ion.Type
	if ionHashValue.IsNull() {
		ionVal = nil
		ionType = ion.NoType
	} else {
		ionVal, err = ionHashValue.value()
		if err != nil {
			return err
		}
		ionType = ionHashValue.Type()
	}

	scalarBytes, err := ss.getBytes(ionHashValue.Type(), ionVal, ionHashValue.IsNull())
	if err != nil {
		return err
	}

	if ionHashValue.Type() != ion.SymbolType {
		ionVal = nil
	}

	tq, representation, err := ss.scalarOrNullSplitParts(ionType, ionHashValue.IsNull(), scalarBytes)
	if err != nil {
		return err
	}

	err = ss.write([]byte{tq})
	if err != nil {
		return err
	}

	if len(representation) > 0 {
		err = ss.write(escape(representation))
		if err != nil {
			return err
		}
	}

	err = ss.endMarker()
	if err != nil {
		return err
	}

	err = ss.handleAnnotationsEnd(ionHashValue, false)
	if err != nil {
		return err
	}

	return nil
}

func (ss *scalarSerializer) stepIn(ionValue interface{}) error {
	return ss.baseSerializer.stepIn(ionValue.(hashValue))
}

func (ss *scalarSerializer) handleFieldName(ionValue interface{}) error {
	return ss.baseSerializer.handleFieldName(ionValue.(hashValue))
}

func (ss *scalarSerializer) handleAnnotationsBegin(ionValue interface{}) error {
	return ss.baseSerializer.handleAnnotationsBegin(ionValue.(hashValue), false)
}

func (ss *scalarSerializer) handleAnnotationsEnd(ionValue interface{}, isContainer bool) error {
	return ss.baseSerializer.handleAnnotationsEnd(ionValue.(hashValue), isContainer)
}
