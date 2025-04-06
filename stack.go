/*
 * bfi: Simple stack implementation for use in `program.go`.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"errors"
)

type stack[T any] struct {
	data   []T
	length int
}

func MakeStack[T any]() *stack[T] {
	output := stack[T]{make([]T, 0), 0}
	return &output
}

func (stack *stack[T]) Push(element T) {
	stack.data = append(stack.data, element)
	stack.length++
}

func (stack *stack[T]) Pop() (T, error) {
	var output T

	if stack.length == 0 {
		return output, errors.New("stack is empty")
	}

	output = stack.data[stack.length-1]
	stack.data = stack.data[0 : stack.length-1]
	stack.length--

	return output, nil
}

func (stack *stack[T]) Length() int {
	return stack.length
}
