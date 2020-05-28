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

// Holds the commonalities between scalar and struct serializers.
type baseSerializer struct {
	hashFunction           IonHasher
	depth                  int
	hasContainerAnnotation bool
}

func (baseSerializer *baseSerializer) stepOut() error {
	err := baseSerializer.endMarker()
	if err != nil {
		return err
	}

	err = baseSerializer.handleAnnotationsEnd(nil, true)
	if err != nil {
		return err
	}

	return nil
}

func (baseSerializer *baseSerializer) stepIn(ionValue interface{}) error {
	panic("implement me")
}

func (baseSerializer *baseSerializer) sum(b []byte) []byte {
	panic("implement me")
}

func (baseSerializer *baseSerializer) handleFieldName(ionValue interface{}) error {
	panic("implement me")
}

func (baseSerializer *baseSerializer) update(bytes []byte) error {
	panic("implement me")
}

func (baseSerializer *baseSerializer) beginMarker() error {
	panic("implement me")
}

func (baseSerializer *baseSerializer) endMarker() error {
	panic("implement me")
}

func (baseSerializer *baseSerializer) handleAnnotationsBegin(ionValue interface{}) error {
	panic("implement me")
}

func (baseSerializer *baseSerializer) handleAnnotationsEnd(ionValue interface{}, isContainer bool) error {
	panic("implement me")
}

func (baseSerializer *baseSerializer) writeSymbol(token string) error {
	panic("implement me")
}

func (baseSerializer *baseSerializer) getBytes(ionType ion.Type, ionValue interface{}, isNull bool) ([]byte, error) {
	panic("implement me")
}

func (baseSerializer *baseSerializer) getLengthFieldLength(bytes []byte) (int, error) {
	panic("implement me")
}

// SymbolToken is currently not available
//func (baseSerializer *baseSerializer)scalarOrNullSplitParts(ionType ion.Type, symbolToken ion.SymbolToken, isNull bool, bytes byte[]) (byte, []byte) {
//	panic("implement me")
//}

func escape(bytes []byte) []byte {
	panic("implement me")
}

func serializers(ionType ion.Type, ionValue interface{}, writer HashWriter) {
	panic("implement me")
}

func tq(ionValue interface{}) []byte {
	panic("implement me")
}
