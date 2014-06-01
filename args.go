package args

// Entry point for command-line argument parsing.
// Constructs a new Expectation object.
func Load(args []string) (Args) {
	exp := argv{
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
type Args interface {
	// Consumes a single flag from the command-line arguments, if found.
	//
	// name: The name of the flag to look for.
	// alts: Any other names that could be used to refer to the same option,
	// including single-character names. Any single-character names should
	// appear in the form "-n".
	//
	// The flag can only be accessed by its primary name, not by any of the
	// alternate names.
	AllowFlag(name string, alts ...string) (chain Args)

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
	AllowOption(name string, alts ...string) (chain Args)

	// Consumes the next argument from the command-line as a parameter.
	// 
	// If there are no more arguments to consume, nothing will be consumed.
	AllowParam() (chain Args)

	// Consumes the next argument from the command-line as a parameter,
	// giving it the specified name.
	// 
	// If there are no more arguments to consume, nothing will be consumed,
	// and the named parameter will not be present in the result.
	// name: The name that will be assigned to the parameter.
	AllowParamNamed(name string) (chain Args)

	// Consumes a single flag from the command-line arguments.
	// The flag must be present, otherwise validation will fail.
	// name: The name of the flag to look for.
	//
	// The flag can only be accessed by its name or position, not by
	// any of the alternate names.
	ExpectFlag(name string, alts ...string) (chain Args)

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
	ExpectOption(name string, alts ...string) (chain Args)

	// Consumes the next argument from the command-line as a parameter.
	//
	// If there are no more arguments to consume, validation will fail.
	ExpectParam() (chain Args)

	// Consumes the next argument from the command-line as a parameter,
	// giving it the specified name.
	//
	// If there are no more arguments to consume, validation will fail.
	// name: The name that will be assigned to the parameter.
	ExpectParamNamed(name string) (chain Args)

	// Discards any remaining, unconsumed arguments, so that they will
	// not cause validation to fail.
	//
	// Alternately, can be used to force the next Allow or Expect to fail.
	Chop() (chain Args)

	// Consumes and returns any remaining, unconsumed arguments, so that
	// they will not cause validation to fail. The remaining arguments are
	// returned as a slice of strings.
	//
	// Alternately, can be used to force the next Allow or Expect to fail.
	ChopSlice() (result Args, tail []string)

	// Consumes and returns any remaining, unconsumed arguments, so that
	// they will not cause validation to fail. The remaining arguments are
	// returned as a single string, concatenated with spaces.
	//
	// Alternately, can be used to force the next Allow or Expect to fail.
	ChopString() (result Args, tail string)

	// Discards any remaining, unconsumed arguments and calls Validate.
	ChopAndValidate() (result Args, err error)

	// Called once all expectations have been specified, to parse and
	// validate the arguments.
	Validate() (chain Args, err error)

	// Gets whether the named flag was set.
	// name: The name of the flag to check.
	Flag(name string) bool

	// Checks whether the named flag was checked.
	//
	// NOTE: Does not check whether the flag was present
	// in the arguments; checks only whether it was
	// Expected or Allowed.
	// name: The name of the flag.
	HasFlag(name string) bool

	// Checks whether the named option was found.
	//
	// Use this before calling Option on an Allowed
	// (not Expected) option.
	// name: The name of the option.
	HasOption(name string) bool

	// Checks whether there is a parameter at
	// the specified index.
	// i: The 0-based index of the parameter.
	HasParamAt(i int) bool

	// Checks whether there is a parameter with
	// the specified name.
	// i: The name of the parameter.
	HasParamNamed(name string) bool

	// Gets the value of the named option.
	// name: The name of the option.
	Option(name string) string

	// Gets the value of the parameter at the specified position.
	// i: The 0-based index of the parameter.
	ParamAt(i int) string

	// Gets the value of the named parameter.
	// name: The name of the parameter.
	ParamNamed(name string) string
}
