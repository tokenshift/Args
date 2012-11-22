package args

import (
	"fmt"
	"strings"
)

/* Implementation of the Expectation interface. */
type expectation struct {
	/* Will contain a copy of the arguments passed to Args. */
	args []string

	/* Tracks whether an argument at the matching index has been consumed. */
	consumed []bool

	/* Contains any errors that occur while matching arguments. */
	errors []error

	/* Map of flags that were checked. */
	flags map[string]bool

	/* Map of options that were set. */
	options map[string]string

	/* List of positional parameters. */
	parameters []string

	/* Map of parameter names to their index in
	 * the list of parameters. */
	namedParameters map[string]int
}

/* Consumes a single flag from the command-line arguments, if found.
 *
 * name: The name of the flag to look for. */
func (chain expectation) AllowFlag(name string, alts ...string) Expectation {
	chain, _ = chain.getFlag(name)
	return chain
}

/* Consumes a single option and its value from the command-line arguments,
 * if found.
 *
 * name: The name of the option, as it will appear in the form "--name".
 * alts: Any other names that could be used to refer to the same option,
 * including single-character names. Any single-character names should
 * appear in the form "-n".
 *
 * The option can only be accessed by its name or position, not by
 * any of the alternate names. */
func (chain expectation) AllowOption(name string, alts ...string) Expectation {
	chain, _, _ = chain.getOption(name, alts)

	return chain
}

/* Consumes the next argument from the command-line as a parameter.
 * 
 * If there are no more arguments to consume, nothing will be consumed. */
func (chain expectation) AllowParam() Expectation {
	chain, _, _ = chain.getParam()

	return chain
}

/* Consumes the next argument from the command-line as a parameter,
 * giving it the specified name.
 * 
 * If there are no more arguments to consume, nothing will be consumed,
 * and the named parameter will not be present in the result.
 * name: The name that will be assigned to the parameter. */
func (chain expectation) AllowNamedParam(name string) Expectation {
	chain, index, found := chain.getParam()

	if found {
		chain.namedParameters[name] = index
	}

	return chain
}

/* Consumes a single flag from the command-line arguments.
 * The flag must be present, otherwise validation will fail.
 * name: The name of the flag to look for. */
func (chain expectation) ExpectFlag(name string, alts ...string) Expectation {
	chain, found := chain.getFlag(name)

	if !found {
		chain.errors = append(chain.errors, fmt.Errorf("Flag '%v' was expected and not found.", name))
	}

	return chain
}

/* Consumes a single option and its value from the command-line arguments.
 *
 * If the option is not found, validation will fail.
 * name: The name of the option, as it will appear in the form "--name".
 * alts: Any other names that could be used to refer to the same option,
 * including single-character names. Any single-character names should
 * appear in the form "-n".
 *
 * The option can only be accessed by its name or position, not by
 * any of the alternate names. */
func (chain expectation) ExpectOption(name string, alts ...string) Expectation {
	chain, _, found := chain.getOption(name, alts)

	if !found {
		chain.errors = append(chain.errors, fmt.Errorf("Option '%v' was expected and not found.", name))
	}

	return chain
}

/* Consumes the next argument from the command-line as a parameter.
 *
 * If there are no more arguments to consume, validation will fail. */
func (chain expectation) ExpectParam() Expectation {
	chain, _, found := chain.getParam()

	if !found {
		chain.errors = append(chain.errors, fmt.Errorf("No more arguments to consume."))
	}

	return chain
}

/* Consumes the next argument from the command-line as a parameter,
 * giving it the specified name.
 *
 * If there are no more arguments to consume, validation will fail.
 * name: The name that will be assigned to the parameter. */
func (chain expectation) ExpectNamedParam(name string) Expectation {
	chain, index, found := chain.getParam()

	if found {
		chain.namedParameters[name] = index
	} else {
		chain.errors = append(chain.errors, fmt.Errorf("No more arguments to consume."))
	}

	return chain
}

/* Helper Methods */

/* Checks whether the specified flag is present.
 *
 * name: The name of the flag to look for. */
func (chain expectation) getFlag(name string) (out expectation, present bool) {
	for i, arg := range chain.args {
		if chain.consumed[i] {
			continue
		}

		if strings.HasPrefix(arg, "--") && arg[2:] == name {
			present = true
			chain.consumed[i] = true
			break
		}
	}

	chain.flags[name] = present

	out = chain

	return
}

/* Retrieves (if found) an option with the specified name.
 *
 * name: The primary name of the option.
 * alts: Any other names the option could have. */
func (chain expectation) getOption(name string, alts []string) (out expectation, val string, found bool) {
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
				val = out.args[i+1]
				out.consumed[i] = true
				out.consumed[i+1] = true
				break
			}
		}
	}

	if found {
		chain.options[name] = val
	}

	return
}

/* Retrieves the next positional parameter, if there is one. */
func (chain expectation) getParam() (out expectation, index int, found bool) {
	out = chain

	for i, val := range chain.args {
		if out.consumed[i] {
			continue
		}

		found = true
		out.consumed[i] = true
		index = len(out.parameters)
		out.parameters = append(out.parameters, val)
		break
	}

	return
}
