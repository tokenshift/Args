# Args.go #

Command-line argument parsing and validation for Go.

## Use ##

Args.go recognizes three types of arguments: parameters, options and flags. Options and flags are arguments of the form "--foo" or "-f"; options are immediately followed by an option value, while flags are standalone, their presence being interpreted as a bool. A parameter is simply a single argument that is neither an option nor a flag.

Options and flags can occur anywhere in the list of arguments, as long as they have not already been consumed. Parameters must be consumed in the order that they are included on the command line.

Args.go validates and parses a set of command-line arguments simultaneously.