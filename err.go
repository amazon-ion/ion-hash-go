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

	"github.com/amzn/ion-go/ion"
)

// An InvalidOperationError is returned when a method call is invalid for the struct's current state.
type InvalidOperationError struct {
	structName string
	methodName string
	message    string
}

func (e *InvalidOperationError) Error() string {
	if e.message != "" {
		return fmt.Sprintf(`ionhash: Invalid operation at %v.%v: %v`, e.structName, e.methodName, e.message)
	}

	return fmt.Sprintf(`ionhash: Invalid operation error in %v.%v`, e.structName, e.methodName)
}

// InvalidArgumentError is returned when one of the arguments given to a function was not valid.
type InvalidArgumentError struct {
	argumentName  string
	argumentValue interface{}
}

func (e *InvalidArgumentError) Error() string {
	return fmt.Sprintf(`ionhash: Invalid value: "%v" specified for argument: %s`, e.argumentValue, e.argumentName)
}

// InvalidIonTypeError is returned when processing an unexpected Ion type.
type InvalidIonTypeError struct {
	ionType ion.Type
}

func (e *InvalidIonTypeError) Error() string {
	return fmt.Sprintf(`ionhash: Invalid Ion type: %s`, e.ionType.String())
}

// UnknownSymbolError is returned when processing an unknown field name symbol.
type UnknownSymbolError struct {
	sid int64
}

func (e *UnknownSymbolError) Error() string {
	return fmt.Sprintf(`ionhash: Unknown text for sid %d`, e.sid)
}
