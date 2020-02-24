# Slang

Slang (short for Sabre Lang) is a tiny LISP built using [Sabre](https://github.com/spy16/sabre).

## Installation

1. Download a release from [Releases](https://github.com/spy16/slang/releases) for
   your target platfomr.
2. Extract files and add extraction path to `PATH` variable.

## Usage

1. `slang` for REPL
2. `slang -e "(+ 1 2 3)"` for executing string
3. `slang -f "examples/simple.lisp"` for executing file.

> If you specify both -f and -e flags, file will be executed first and then the string
> will be executed in the same scope and you will be dropped into REPL. If REPL not needed,
> use -norepl option.
