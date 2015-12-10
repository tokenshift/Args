# Args.go #

Command-line argument parsing and validation for Go.

## Arguments ##

**Args.go** recognizes three types of arguments: parameters, options and flags.
Flags are standalone arguments of the form "--flag" or "-f"; options are the
same format, but immediately followed by a (single) value. A parameter is simply
a single argument, without a leading hyphen.

Flags and options will be consumed wherever they are found in the arguments.
Parameters are consumed in the order they are included on the command line.

## Operations

The Args library works on an array of command-line arguments (usually
`os.Args[1:]`). Each operation returns the argument requested, as well as the
array with that argument removed, so that it can be passed to additional Arg
calls.