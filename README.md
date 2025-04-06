# bfi

A very simple interpreter for the Brainfuck programming language, written in Go.

Aiming to be fully-compatible with programs written for [the reference written by Daniel B. Cristofani](https://brainfuck.org/brainfuck.html).

Extensible-ish by design, new instructions can be added by modifying the `OPERATIONS` object in `interpreter.go`.

## Usage

```
Usage: bfi [--exit-on-eof] FILE

Positional arguments:
  FILE                   program source file

Options:
  --exit-on-eof, -e      exit program when an EOF byte is read from standard input
  --help, -h             display this help and exit
```

## Future Plans

- [ ] Support for `#` instruction for debugging
- [ ] Optional unbounded mode (removal of 30k memory limit)
- [ ] Optional extended instructions

[^1]: Excluding the `#` instruction.