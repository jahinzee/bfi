/*
 * bfi: Program (source code) runtime and interpreting functions,
 * including instruction logic.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"bufio"
	"fmt"
	"os"
)

type behaviourArgs struct {
	// If true, exit gracefully when an EOF byte is read from standard input
	// on a ',' instruction.
	HaltOnEOF bool
}

type process struct {
	// Basic process runtime data.
	Program        *program
	ProgramCounter int
	Memory         map[int]int
	MemoryCounter  int

	// If Halt is true, the process should stop reading instructions and halt.
	// If Error is not nil in this state, it will be propagated up. Otherwise,
	// the process just halts gracefully.
	RuntimeHalt  bool
	RuntimeError error

	// Environment arguments for customising behaviour.
	BehaviourArgs behaviourArgs

	// Additional things for handling I/O.
	LastCharOutput rune
	StdinReader    *bufio.Reader
}

type OperationMap = map[rune]func(*process)

var OPERATIONS OperationMap = OperationMap{
	'+': func(proc *process) {
		proc.Memory[proc.MemoryCounter]++
	},

	'-': func(proc *process) {
		proc.Memory[proc.MemoryCounter]--
	},

	'>': func(proc *process) {
		proc.MemoryCounter++
	},

	'<': func(proc *process) {
		proc.MemoryCounter--
	},
	'.': func(proc *process) {
		char := rune(proc.Memory[proc.MemoryCounter])
		fmt.Printf("%c", char)
		proc.LastCharOutput = char
	},
	',': func(proc *process) {
		char, _, err := proc.StdinReader.ReadRune()
		if err != nil {
			isEOF := err.Error() == "EOF"
			haltOnEOF := args.HaltOnEOF

			// Halt on EOF, unless we're not supposed to.
			if isEOF && haltOnEOF {
				proc.RuntimeHalt = true
			} else if !isEOF {
				proc.RuntimeHalt = true
				proc.RuntimeError = fmt.Errorf("error reading standard input: %s", err)
			}
		}
		proc.Memory[proc.MemoryCounter] = int(char)
	},
	'[': func(proc *process) {
		if proc.Memory[proc.MemoryCounter] == 0 {
			jump, err := proc.Program.LoopBounds.GetOther(proc.ProgramCounter)
			if err != nil {
				proc.RuntimeHalt = true
				proc.RuntimeError = fmt.Errorf("unmatched '[' instruction")
			}
			proc.ProgramCounter = jump
		}
	},
	']': func(proc *process) {
		if proc.Memory[proc.MemoryCounter] != 0 {
			jump, err := proc.Program.LoopBounds.GetOther(proc.ProgramCounter)
			if err != nil {
				proc.RuntimeHalt = true
				proc.RuntimeError = fmt.Errorf("unmatched ']' instruction")
			}
			proc.ProgramCounter = jump
		}
	},
}

func RunProgram(prog *program, args behaviourArgs) error {
	proc := process{
		Program:        prog,
		ProgramCounter: 0,
		Memory:         make(map[int]int),
		MemoryCounter:  0,
		RuntimeHalt:    false,
		RuntimeError:   nil,
		BehaviourArgs:  args,
		LastCharOutput: '\n',
		StdinReader:    bufio.NewReader(os.Stdin),
	}

	for !proc.RuntimeHalt {
		instruction := rune(prog.Text[proc.ProgramCounter])
		operation, exists := OPERATIONS[instruction]

		// If the character doesn't exist in the operation map, it's a comment.
		if exists {
			operation(&proc)
		}

		// Post-operation cleanup; wrap values and increment ProgramCounter.
		proc.MemoryCounter %= 30000
		proc.Memory[proc.MemoryCounter] %= 256
		proc.ProgramCounter++

		// Are we out of tape?
		if proc.ProgramCounter >= proc.Program.Length {
			proc.RuntimeHalt = true
		}
	}
	// Just so we're not leaving the program on an hanging line.
	if proc.LastCharOutput != '\n' {
		fmt.Println()
	}
	if proc.RuntimeError != nil {
		return fmt.Errorf("char %d: %s", proc.ProgramCounter, proc.RuntimeError)
	}
	return nil
}
