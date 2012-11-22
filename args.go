package args

import (
	"fmt"
	"strings"
)

/* Entry point for command-line argument parsing.
 *
 * Constructs a new Expectation object. */
func Args(args []string) (parser Expectation) {
	exp := expectation{
		args:     make([]string, len(args)),
		consumed: make([]bool, len(args)),

		errors: make([]error, 0, len(args)),

		flags:           make(map[string]bool),
		options:         make(map[string]string),
		parameters:      make([]string, 0, len(args)),
		namedParameters: make(map[string]int),
	}

	copy(exp.args, args)

	return exp
}

/* Represents a set of parse-and-validation instructions.
 *
 * Expectation methods are intended to be chained, ala jQuery.
 * The returned object may or may not be the same object as
 * the original; do not depend on either behavior! */
type Expectation interface {
	/* Allowances */
	// These argments MAY be present.

	/* Consumes a single flag from the command-line arguments, if found.
	 *
	 * name: The name of the flag to look for.
	 * alts: Any other names that could be used to refer to the same option,
	 * including single-character names. Any single-character names should
	 * appear in the form "-n".
	 *
	 * The flag can only be accessed by its name or position, not by
	 * any of the alternate names. */
	AllowFlag(name string, alts ...string) (chain Expectation)

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
	AllowOption(name string, alts ...string) (chain Expectation)

	/* Consumes the next argument from the command-line as a parameter.
	 * 
	 * If there are no more arguments to consume, nothing will be consumed. */
	AllowParam() (chain Expectation)

	/* Consumes the next argument from the command-line as a parameter,
	 * giving it the specified name.
	 * 
	 * If there are no more arguments to consume, nothing will be consumed,
	 * and the named parameter will not be present in the result.
	 * name: The name that will be assigned to the parameter. */
	AllowNamedParam(name string) (chain Expectation)

	/* Expectations */
	// These argments MUST be present.

	/* Consumes a single flag from the command-line arguments.
	 * The flag must be present, otherwise validation will fail.
	 * name: The name of the flag to look for.
	 *
	 * The flag can only be accessed by its name or position, not by
	 * any of the alternate names. */
	ExpectFlag(name string, alts ...string) (chain Expectation)

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
	ExpectOption(name string, alts ...string) (chain Expectation)

	/* Consumes the next argument from the command-line as a parameter.
	 *
	 * If there are no more arguments to consume, validation will fail. */
	ExpectParam() (chain Expectation)

	/* Consumes the next argument from the command-line as a parameter,
	 * giving it the specified name.
	 *
	 * If there are no more arguments to consume, validation will fail.
	 * name: The name that will be assigned to the parameter. */
	ExpectNamedParam(name string) (chain Expectation)

	/* Termination */

	/* Discards any remaining, unconsumed arguments, so that they will
	 * not cause validation to fail.
	 *
	 * Alternately, can be used to force the next Allow or Expect to fail. */
	Chop() (chain Expectation)

	/* Discards any remaining, unconsumed arguments and calls Validate. */
	ChopAndValidate() (result Expectation, err error)

	/* Called once all expectations have been specified, to parse and
	 * validate the arguments. */
	Validate() (result Expectation, err error)

	/* Results */

	/* Gets whether the named flag was set.
	 * name: The name of the flag to check. */
	Flag(name string) (value bool, err error)

	/* Checks whether the named flag was checked.
	 *
	 * NOTE: Does not check whether the flag was present
	 * in the arguments; checks only whether it was
	 * Expected or Allowed.
	 * name: The name of the flag. */
	HasFlag(name string) (present bool)

	/* Checks whether the named option was found.
	 *
	 * Use this before calling Option on an Allowed
	 * (not Expected) option.
	 * name: The name of the option. */
	HasOption(name string) (present bool)

	/* Checks whether there is a parameter at
	 * the specified index.
	 * i: The 0-based index of the parameter. */
	HasParamAt(i int) (present bool)

	/* Checks whether there is a parameter with
	 * the specified name.
	 * i: The name of the parameter. */
	HasParamNamed(name string) (present bool)

	/* Gets the value of the named option.
	 * name: The name of the option. */
	Option(name string) (value string, err error)

	/* Gets the value of the parameter at the specified position.
	 * i: The 0-based index of the parameter. */
	ParamAt(i int) (param string, err error)

	/* Gets the value of the named parameter.
	 * name: The name of the parameter. */
	ParamNamed(name string) (param string, err error)
}

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

/* Discards any remaining, unconsumed arguments, so that they will
 * not cause validation to fail.
 *
 * Alternately, can be used to force the next Allow or Expect to fail. */
func (chain expectation) Chop() Expectation {
	for i, _ := range chain.consumed {
		chain.consumed[i] = true
	}
	return chain
}

/* Discards any remaining, unconsumed arguments and calls Validate. */
func (final expectation) ChopAndValidate() (result Expectation, err error) {
	result, err = final.Chop().Validate()
	return
}

/* Called once all expectations have been specified, to parse and
 * validate the arguments.
 *
 * Validation will fail if:
 * 1. There are unconsumed arguments remaining.
 *    Call Chop() to consume any remaining arguments.
 * 2. Any Expected arguments were not found.
 *    Allowed arguments will not cause validation errors when missing. */
func (final expectation) Validate() (result Expectation, err error) {
	count := 0

	for _, consumed := range final.consumed {
		if !consumed {
			count++
		}
	}

	if count > 0 {
		final.errors = append(final.errors, fmt.Errorf("%v arguments were not properly consumed.", count))
	}

	if len(final.errors) > 0 {
		err = ArgsError{final.errors}
	}

	result = final

	return
}

/* Error type to represent failed argument validation.
 *
 * Errors: The list of individual validation errors
 * that were encountered. */
type ArgsError struct {
	Errors []error
}

/* Display string for ArgsError.
 *
 * Displays the list of validation errors. */
func (argsError ArgsError) Error() string {
	return fmt.Sprintf("Validation failed: %v", argsError.Errors)
}

/* Gets whether the named flag was set.
 * name: The name of the flag to check. */
func (final expectation) Flag(name string) (value bool, err error) {
	value, ok := final.flags[name]

	if !ok {
		err = fmt.Errorf("You must explicitly Expect or Allow the flag '%v'.", name)
	}

	return
}

/* Checks whether the named flag was checked.
 *
 * NOTE: Does not check whether the flag was present
 * in the arguments; checks only whether it was
 * Expected or Allowed.
 * name: The name of the flag. */
func (final expectation) HasFlag(name string) (present bool) {
	_, present = final.flags[name]
	return
}

/* Checks whether the named option was found.
 *
 * Use this before calling Option on an Allowed
 * (not Expected) option.
 * name: The name of the option. */
func (final expectation) HasOption(name string) (present bool) {
	_, present = final.options[name]
	return
}

/* Checks whether there is a parameter at
 * the specified index.
 * i: The 0-based index of the parameter. */
func (final expectation) HasParamAt(i int) (present bool) {
	present = len(final.parameters) > i
	return
}

/* Checks whether there is a parameter with
 * the specified name.
 * i: The name of the parameter. */
func (final expectation) HasParamNamed(name string) (present bool) {
	_, present = final.namedParameters[name]
	return
}

/* Gets the value of the named option.
 * name: The name of the option. */
func (final expectation) Option(name string) (value string, err error) {
	value, present := final.options[name]

	if !present {
		err = fmt.Errorf("Option '%v' was not found.", name)
	}

	return
}

/* Gets the value of the parameter at the specified position.
 * i: The 0-based index of the parameter. */
func (final expectation) ParamAt(index int) (value string, err error) {
	if len(final.parameters) > index {
		value = final.parameters[index]
	} else {
		value = "ERROR"
		err = fmt.Errorf("No parameter present at index %v.", index)
	}
	return
}

/* Gets the value of the named parameter.
 * name: The name of the parameter. */
func (final expectation) ParamNamed(name string) (value string, err error) {
	index, found := final.namedParameters[name]

	if found {
		value = final.parameters[index]
	}

	if !found {
		err = fmt.Errorf("No parameter present with name %v.", name)
	}

	return
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

	names := make([]string, 0, len(alts) + 1)
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