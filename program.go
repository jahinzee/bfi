/*
 * bfi: Program (source code) pre-processing functions.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"errors"
	"fmt"
)

type program struct {
	Text       string
	Length     int
	LoopBounds loopBounds
}

type loopBounds []struct {
	Left  int
	Right int
}

func (lb loopBounds) GetOther(element int) (int, error) {
	for _, lb := range lb {
		if lb.Right == element {
			return lb.Left, nil
		}
		if lb.Left == element {
			return lb.Right, nil
		}
	}
	return 0, errors.New("could not find other element")
}

func (prog *program) MakeLoopBounds() error {
	stack := MakeStack[int]()
	output := make(loopBounds, 0)

	for idx, rune := range prog.Text {
		if rune == '[' {
			stack.Push(idx)
		}
		if rune == ']' {
			result, err := stack.Pop()
			if err != nil {
				return fmt.Errorf("char %d: unmatched ']' instruction", idx)
			}
			output = append(output, struct {
				Left  int
				Right int
			}{result, idx})
		}
	}

	// If, after we've exhausted the program text, we still have anything
	// left in the stack, atleast one unmatched '[' instruction remains.
	result, err := stack.Pop()
	if err == nil {
		return fmt.Errorf("char %d: unmatched '[' instruction", result)
	}

	prog.LoopBounds = output
	return nil
}

func MakeProgram(text string) (*program, error) {
	var output program
	output.Text = text
	output.Length = len(text)

	err := output.MakeLoopBounds()
	if err != nil {
		return nil, err
	}

	return &output, nil
}
