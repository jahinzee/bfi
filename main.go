/*
 * bfi: Main command-line entry point.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/wzshiming/ctc"
)

var args struct {
	FILE      string `arg:"positional,required" help:"program source file"`
	HaltOnEOF bool   `arg:"-e,--exit-on-eof" help:"exit program when an EOF byte is read from standard input"`
}

func main() {
	arg.MustParse(&args)

	file, err := os.ReadFile(args.FILE)
	if err != nil {
		fmt.Print(ctc.ForegroundRed)
		fmt.Printf("Cannot open file:\n %s\n", err)
		fmt.Print(ctc.Reset)
		os.Exit(1)
	}

	prog, err := MakeProgram(string(file))
	if err != nil {
		fmt.Print(ctc.ForegroundRed)
		fmt.Printf("Parsing error:\n  %s\n", err)
		fmt.Print(ctc.Reset)
		os.Exit(1)
	}

	result := RunProgram(prog, behaviourArgs{
		HaltOnEOF: args.HaltOnEOF,
	})
	if result != nil {
		fmt.Print(ctc.ForegroundRed)
		fmt.Printf("Runtime error:\n  %s\n", result)
		fmt.Print(ctc.Reset)
		os.Exit(1)
	}
}
