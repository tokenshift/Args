package args

import (
	"fmt"
	"strings"
)

/* Entry point for command-line argument parsing.
 *
 * Constructs a new Expectation object. */
func Args(args []string) (parser Expectation) {
	exp := expectation {
		args: make([]string, len(args)),
		consumed: make([]bool, len(args)),
		flags: make(map[string]bool),
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
}

type processor func(exp expectation) (err error)

/* Consumes a single flag from the command-line arguments, if found.
 *
 * name: The name of the flag to look for. */
func (old expectation) AllowFlag(name string, alts ...string) (Expectation) {
	old, _ = old.getFlag(name)
	return old
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
func (old expectation) AllowOption(name string, alts ...string) (Expectation) {
	return old
}

/* Consumes the next argument from the command-line as a parameter.
 * 
 * If there are no more arguments to consume, nothing will be consumed. */
func (old expectation) AllowParam() (Expectation) {
	return old
}

/* Consumes the next argument from the command-line as a parameter,
 * giving it the specified name.
 * 
 * If there are no more arguments to consume, nothing will be consumed,
 * and the named parameter will not be present in the result.
 * name: The name that will be assigned to the parameter. */
func (old expectation) AllowNamedParam(name string) (Expectation) {
	return old
}

/* Consumes a single flag from the command-line arguments.
 * The flag must be present, otherwise validation will fail.
 * name: The name of the flag to look for. */
func (old expectation) ExpectFlag(name string, alts ...string) (Expectation) {
	old, found := old.getFlag(name)

	if !found {
		old.errors = append(old.errors, fmt.Errorf("Flag '%v' was expected and not found."))
	}

	return old
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
func (old expectation) ExpectOption(name string, alts ...string) (Expectation) {
	return old
}

/* Consumes the next argument from the command-line as a parameter.
 *
 * If there are no more arguments to consume, validation will fail. */
func (old expectation) ExpectParam() (Expectation) {
	return old
}

/* Consumes the next argument from the command-line as a parameter,
 * giving it the specified name.
 *
 * If there are no more arguments to consume, validation will fail.
 * name: The name that will be assigned to the parameter. */
func (old expectation) ExpectNamedParam(name string) (Expectation) {
	return old
}

/* Discards any remaining, unconsumed arguments, so that they will
 * not cause validation to fail.
 *
 * Alternately, can be used to force the next Allow or Expect to fail. */
func (old expectation) Chop() (Expectation) {
	return old
}

/* Discards any remaining, unconsumed arguments and calls Validate. */
func (final expectation) ChopAndValidate() (result Expectation, err error) {
	err = fmt.Errorf("Not yet implemented.")
	result = final
	return
}

/* Called once all expectations have been specified, to parse and
 * validate the arguments. */
func (final expectation) Validate() (result Expectation, err error) {
	count := 0

	for _, consumed := range final.consumed {
		if !consumed {
			count++
		}
	}

	if count > 0 {
		final.errors = append(final.errors,fmt.Errorf("%v arguments were not properly consumed.", count))
	}

	if len(final.errors) > 0 {
		err = ArgsError {final.errors}
	}
	
	result = final

	return
}

type ArgsError struct {
	Errors []error
}


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

/* Gets the value of the named option.
 * name: The name of the option. */
func (final expectation) Option(name string) (value string, err error) {
	err = fmt.Errorf("Not yet implemented.")
	return
}

/* Gets the value of the parameter at the specified position.
 * i: The 0-based index of the parameter. */
func (final expectation) ParamAt(index int) (value string, err error) {
	err = fmt.Errorf("Not yet implemented.")
	value = ""
	return
}

/* Gets the value of the named parameter.
 * name: The name of the parameter. */
func (final expectation) ParamNamed(name string) (value string, err error) {
	err = fmt.Errorf("Not yet implemented.")
	value = ""
	return
}


/* Helper Methods */

/* Checks whether the specified flag is present. */
func (old expectation) getFlag(name string) (chain expectation, present bool) {
	for i, arg := range old.args {
		if old.consumed[i] {
			continue
		}

		if strings.HasPrefix(arg, "--") && arg[2:] == name {
			present = true
			old.consumed[i] = true
			break
		}
	}

	old.flags[name] = present

	chain = old

	return
}