package args

// Entry point for command-line argument parsing.
// Constructs a new Expectation object.
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

// Represents a set of parse-and-validation instructions.
//
// Expectation methods are intended to be chained, ala jQuery.
// The returned object may or may not be the same object as
// the original; do not depend on either behavior!
type Expectation interface {
	// Consumes a single flag from the command-line arguments, if found.
	//
	// name: The name of the flag to look for.
	// alts: Any other names that could be used to refer to the same option,
	// including single-character names. Any single-character names should
	// appear in the form "-n".
	//
	// The flag can only be accessed by its name or position, not by
	// any of the alternate names.
	AllowFlag(name string, alts ...string) (chain Expectation)

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
	AllowOption(name string, alts ...string) (chain Expectation)

	// Consumes the next argument from the command-line as a parameter.
	// 
	// If there are no more arguments to consume, nothing will be consumed.
	AllowParam() (chain Expectation)

	// Consumes the next argument from the command-line as a parameter,
	// giving it the specified name.
	// 
	// If there are no more arguments to consume, nothing will be consumed,
	// and the named parameter will not be present in the result.
	// name: The name that will be assigned to the parameter.
	AllowNamedParam(name string) (chain Expectation)

	// Consumes a single flag from the command-line arguments.
	// The flag must be present, otherwise validation will fail.
	// name: The name of the flag to look for.
	//
	// The flag can only be accessed by its name or position, not by
	// any of the alternate names.
	ExpectFlag(name string, alts ...string) (chain Expectation)

	// Consumes a single option and its value from the command-line arguments.
	//
	// If the option is not found, validation will fail.
	// name: The name of the option, as it will appear in the form "--name".
	// alts: Any other names that could be used to refer to the same option,
	// including single-character names. Any single-character names should
	// appear in the form "-n".
	//
	// The option can only be accessed by its name or position, not by
	// any of the alternate names. */
	ExpectOption(name string, alts ...string) (chain Expectation)

	// Consumes the next argument from the command-line as a parameter.
	//
	// If there are no more arguments to consume, validation will fail.
	ExpectParam() (chain Expectation)

	// Consumes the next argument from the command-line as a parameter,
	// giving it the specified name.
	//
	// If there are no more arguments to consume, validation will fail.
	// name: The name that will be assigned to the parameter.
	ExpectNamedParam(name string) (chain Expectation)

	// Discards any remaining, unconsumed arguments, so that they will
	// not cause validation to fail.
	//
	// Alternately, can be used to force the next Allow or Expect to fail.
	Chop() (chain Expectation)

	// Discards any remaining, unconsumed arguments and calls Validate.
	ChopAndValidate() (result Expectation, err error)

	// Called once all expectations have been specified, to parse and
	// validate the arguments.
	Validate() (result Expectation, err error)

	// Gets whether the named flag was set.
	// name: The name of the flag to check.
	Flag(name string) (value bool)

	// Checks whether the named flag was checked.
	//
	// NOTE: Does not check whether the flag was present
	// in the arguments; checks only whether it was
	// Expected or Allowed.
	// name: The name of the flag.
	HasFlag(name string) (present bool)

	// Checks whether the named option was found.
	//
	// Use this before calling Option on an Allowed
	// (not Expected) option.
	// name: The name of the option.
	HasOption(name string) (present bool)

	// Checks whether there is a parameter at
	// the specified index.
	// i: The 0-based index of the parameter.
	HasParamAt(i int) (present bool)

	// Checks whether there is a parameter with
	// the specified name.
	// i: The name of the parameter.
	HasParamNamed(name string) (present bool)

	// Gets the value of the named option.
	// name: The name of the option.
	Option(name string) (value string)

	// Gets the value of the parameter at the specified position.
	// i: The 0-based index of the parameter.
	ParamAt(i int) (param string)

	// Gets the value of the named parameter.
	// name: The name of the parameter.
	ParamNamed(name string) (param string)
}
