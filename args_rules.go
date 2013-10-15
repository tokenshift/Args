package args

import "bytes"
import "fmt"
import "strings"

// Implementation of the Args interface.
type argv struct {
	// Will contain a copy of the arguments passed to Args.
	args []string

	// Tracks whether an argument at the matching index has been consumed.
	consumed []bool

	// Contains any errors that occur while matching arguments.
	errors []error

	// Map of flags that were checked.
	flags map[string]bool

	// Map of options that were set.
	options map[string]string

	// List of positional parameters.
	parameters []string

	// Map of parameter names to their index in
	// the list of parameters.
	namedParameters map[string]int
}

// Makes a copy of the structure.
func (orig argv) clone() (out argv) {
	out.args			= make([]string, len(orig.args))
	out.consumed		= make([]bool, len(orig.consumed))
	out.errors			= make([]error, len(orig.errors))
	out.flags			= make(map[string]bool)
	out.options			= make(map[string]string)
	out.parameters		= make([]string, len(orig.parameters))
	out.namedParameters	= make(map[string]int)

	copy(out.args, orig.args)
	copy(out.consumed, orig.consumed)
	copy(out.errors, orig.errors)
	copy(out.parameters, orig.parameters)
	
	for key, val := range(orig.flags) {
		out.flags[key] = val
	}

	for key, val := range(orig.options) {
		out.options[key] = val
	}

	for key, val := range(orig.namedParameters) {
		out.namedParameters[key] = val
	}

	return
}

// Outputs the current state of the argument context.
// Useful for debugging.
func (args argv) String() string {
	var out bytes.Buffer
	// ..., ..., ...
	// xxx,    , xxx
	// name		=> value
	// --option	=> value
	// --flag
	for i, arg := range(args.args) {
		if i > 0 {
			fmt.Fprint(&out, ", ")
		}
		fmt.Fprint(&out, arg)
	}

	fmt.Fprintln(&out, "")

	for i, arg := range(args.args) {
		if i > 0 {
			fmt.Fprint(&out, "  ")
		}

		var filler string
		if args.consumed[i] {
			filler = "X"
		} else {
			filler = " "
		}
		for i2 := 0; i2 < len(arg); i2 += 1 {
			fmt.Fprint(&out, filler)
		}
	}

	fmt.Fprintln(&out, "")

	for i, val := range(args.parameters) {
		fmt.Fprintf(&out, "%d => %s\n", i, val)
	}

	for name, i := range(args.namedParameters) {
		fmt.Fprintf(&out, "%s => %d\n", name, i)
	}

	return out.String()
}

// Consumes a single flag from the command-line arguments, if found.
//
// name: The name of the flag to look for.
func (chain argv) AllowFlag(name string, alts ...string) Args {
	chain, _ = chain.getFlag(name, alts)
	return chain
}

// Consumes a single option and its value from the command-line arguments,
// if found.
//
// name: The name of the option, as it will appear in the form "--name".
// alts: Any other names that could be used to refer to the same option,
// including single-character names. Any single-character names should
// appear in the form "-n".
//
// The option can only be accessed by its name or position, not by
// any of the alternate names.
func (chain argv) AllowOption(name string, alts ...string) Args {
	chain, _, _ = chain.getOption(name, alts)

	return chain
}

// Consumes the next argument from the command-line as a parameter.
// 
// If there are no more arguments to consume, nothing will be consumed.
func (chain argv) AllowParam() Args {
	chain, _, _ = chain.getParam()

	return chain
}

// Consumes the next argument from the command-line as a parameter,
// giving it the specified name.
// 
// If there are no more arguments to consume, nothing will be consumed,
// and the named parameter will not be present in the result.
// name: The name that will be assigned to the parameter.
func (chain argv) AllowParamNamed(name string) Args {
	chain, index, found := chain.getParam()

	if found {
		chain.namedParameters[name] = index
	}

	return chain
}

// Consumes a single flag from the command-line arguments.
// The flag must be present, otherwise validation will fail.
// name: The name of the flag to look for.
func (chain argv) ExpectFlag(name string, alts ...string) Args {
	chain, found := chain.getFlag(name, alts)

	if !found {
		chain.errors = append(chain.errors, fmt.Errorf("Flag '%v' was expected and not found.", name))
	}

	return chain
}

// Consumes a single option and its value from the command-line arguments.
//
// If the option is not found, validation will fail.
// name: The name of the option, as it will appear in the form "--name".
// alts: Any other names that could be used to refer to the same option,
// including single-character names. Any single-character names should
// appear in the form "-n".
//
// The option can only be accessed by its name or position, not by
// any of the alternate names.
func (chain argv) ExpectOption(name string, alts ...string) Args {
	chain, _, found := chain.getOption(name, alts)

	if !found {
		chain.errors = append(chain.errors, fmt.Errorf("Option '%v' was expected and not found.", name))
	}

	return chain
}

// Consumes the next argument from the command-line as a parameter.
//
// If there are no more arguments to consume, validation will fail.
func (chain argv) ExpectParam() Args {
	chain, _, found := chain.getParam()

	if !found {
		chain.errors = append(chain.errors, fmt.Errorf("No more arguments to consume."))
	}

	return chain
}

// Consumes the next argument from the command-line as a parameter,
// giving it the specified name.
//
// If there are no more arguments to consume, validation will fail.
// name: The name that will be assigned to the parameter.
func (chain argv) ExpectParamNamed(name string) Args {
	chain, index, found := chain.getParam()

	if found {
		chain.namedParameters[name] = index
	} else {
		chain.errors = append(chain.errors, fmt.Errorf("No more arguments to consume."))
	}

	return chain
}

// Checks whether the specified flag is present.
//
// name: The name of the flag to look for.
func (chain argv) getFlag(name string, alts []string) (out argv, present bool) {
	out = chain.clone()

	names := make([]string, 0, len(alts) + 1)
	names = append(names, name)
	names = append(names, alts...)

	for _, n := range names {
		for i, arg := range out.args {
			if out.consumed[i] {
				continue
			}

			if len(n) == 1 {
				if strings.HasPrefix(arg, "-") && arg[1:] == n {
					present = true
				}
			} else {
				if strings.HasPrefix(arg, "--") && arg[2:] == n {
					present = true
				}
			}

			if (present) {
				out.consumed[i] = true
				break
			}
		}
	}

	out.flags[name] = present

	return
}

// Retrieves (if found) an option with the specified name.
//
// name: The primary name of the option.
// alts: Any other names the option could have.
func (chain argv) getOption(name string, alts []string) (out argv, val string, found bool) {
	out = chain

	names := make([]string, 0, len(alts)+1)
	names = append(names, name)
	names = append(names, alts...)

	for _, n := range names {
		for i, arg := range out.args {
			if out.consumed[i] {
				continue
			}

			if len(n) == 1 {
				if strings.HasPrefix(arg, "-") && arg[1:] == n && len(out.args) > i+1 {
					found = true
				}
			} else {
				if strings.HasPrefix(arg, "--") && arg[2:] == n && len(out.args) > i+1 {
					found = true
				}
			}

			if found {
				out = chain.clone()
				val = out.args[i+1]
				out.consumed[i] = true
				out.consumed[i+1] = true
				break
			}
		}
		
		if found {
			break
		}
	}

	if found {
		out.options[name] = val
	}

	return
}

// Retrieves the next positional parameter, if there is one.
func (chain argv) getParam() (out argv, index int, found bool) {
	out = chain

	for i, val := range chain.args {
		if out.consumed[i] {
			continue
		}

		found = true
		out = chain.clone()
		out.consumed[i] = true
		index = len(out.parameters)
		out.parameters = append(out.parameters, val)
		break
	}

	return
}
