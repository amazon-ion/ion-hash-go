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

import "errors"

// Stack implementation used internally.
type Stack []interface{}

// IsEmpty returns true if the stack is empty, false otherwise.
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

// Push element into the stack.
func (s *Stack) Push(element interface{}) {
	*s = append(*s, element)
}

// Pop removes and returns the top element of the stack. Returns an error if the stack is empty.
func (s *Stack) Pop() (interface{}, error) {
	if s.IsEmpty() {
		return nil, errors.New("pop function called on an empty stack")
	}

	index := len(*s) - 1   // Get the index of the top most element.
	element := (*s)[index] // Index into the slice and obtain the element.
	*s = (*s)[:index]      // Remove it from the stack by slicing it off.

	return element, nil
}

// Peek returns the top element of the stack. Returns an error if the stack is empty.
func (s *Stack) Peek() (interface{}, error) {
	if s.IsEmpty() {
		return nil, errors.New("peek function called on an empty stack")
	}

	index := len(*s) - 1   // Get the index of the top most element.
	element := (*s)[index] // Index into the slice and obtain the element.

	return element, nil
}

// Size returns the number of elements in the Stack.
func (s *Stack) Size() int {
	return len(*s)
}
