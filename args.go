package args

import ("fmt")

/* Entry point for command-line argument parsing.
 *
 * Constructs a new Expectation object. */
func Args(args []string) (parser Expectation) {
	exp := expectation {
		args: make([]string, len(args)),
		processors: make([]processor, len(args)),
	}

	copy(args, exp.args)

	return exp
}

/* Represents the result of parsing and validating
 * a set of command-line arguments.

 * Provides methods to access specific parameter, option
 * and flag values. */
type ArgResult interface {
	/* Gets whether the named flag was set.
	 * name: The name of the flag to check. */
	Flag(name string) (value bool, err error)

	/* Gets the value of the parameter at the specified position.
	 * i: The 0-based index of the parameter. */
	ParamAt(i int) (param string, err error)

	/* Gets the value of the named parameter.
	 * name: The name of the parameter. */
	ParamNamed(name string) (param string, err error)
}

type argResult struct {
}


/* Represents a set of parse-and-validation instructions.
 *
 * Expectation methods are intended to be chained, ala jQuery.
 * The returned object may or may not be the same object as
 * the original; do not depend on either behavior! */
type Expectation interface {
	/* Consumes a single flag from the command-line arguments, if found.
	 *
	 * name: The name of the flag to look for. */
	AllowFlag(name string) (chain Expectation)

	/* Consumes a single option and its value from the command-line arguments,
	 * if found.
	 *
	 * name: The name of the option, as it will appear in the form "--name".
	 * alt: The short/single-character name of the option, as it will appear
	 * in the form "-n".
	 * At least one of the arguments must be non-nil. */
	AllowOption(name string, alt string) (chain Expectation)

	/* Consumes the next argument from the command-line as a parameter.
	 * 
	 * If there are no more arguments to consume, nothing will be consumed. */
	AllowParam() (chain Expectation)

	/* Consumes the next argument from the command-line as a parameter,
	 * giving it the specified name.
	 * 
	 * If there are no more arguments to consume, nothing will be consumed,
	 * and the named parameter will not be present in the ArgResult.
	 * name: The name that will be assigned to the parameter. */
	AllowNamedParam(name string) (chain Expectation)

	/* Consumes a single flag from the command-line arguments.
	 * The flag must be present, otherwise validation will fail.
	 * name: The name of the flag to look for. */
	ExpectFlag(name string) (chain Expectation)

	/* Consumes a single option and its value from the command-line arguments.
	 *
	 * If the option is not found, validation will fail.
	 * name: The name of the option, as it will appear in the form "--name".
	 * alt: The short/single-character name of the option, as it will appear
	 * in the form "-n".
	 * At least one of the arguments must be non-nil. */
	ExpectOption(name string, alt string) (chain Expectation)

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

	/* Discards any remaining, unconsumed arguments, so that they will
	 * not cause validation to fail.
	 *
	 * Alternately, can be used to force the next Allow or Expect to fail. */
	Chop() (chain Expectation)

	/* Discards any remaining, unconsumed arguments and calls Validate. */
	ChopAndValidate() (result ArgResult, err error)

	/* Called once all expectations have been specified, to parse and
	 * validate the arguments. */
	Validate() (result ArgResult, err error)
}

/* Implementation of the Expectation interface. */
type expectation struct {
	args []string
	processors []processor
}

type processor func(exp expectation) (err error)

/* Consumes a single flag from the command-line arguments, if found.
 *
 * name: The name of the flag to look for. */
func (old expectation) AllowFlag(name string) (Expectation) {
	old.processors = append(old.processors, func(exp expectation) (err error) {
		return
	})
	return old
}

/* Consumes a single option and its value from the command-line arguments,
 * if found.
 *
 * name: The name of the option, as it will appear in the form "--name".
 * alt: The short/single-character name of the option, as it will appear
 * in the form "-n".
 * At least one of the arguments must be non-nil. */
func (old expectation) AllowOption(name string, alt string) (Expectation) {
	old.processors = append(old.processors, func(exp expectation) (err error) {
		return
	})
	return old
}

/* Consumes the next argument from the command-line as a parameter.
 * 
 * If there are no more arguments to consume, nothing will be consumed. */
func (old expectation) AllowParam() (Expectation) {
	old.processors = append(old.processors, func(exp expectation) (err error) {
		return
	})
	return old
}

/* Consumes the next argument from the command-line as a parameter,
 * giving it the specified name.
 * 
 * If there are no more arguments to consume, nothing will be consumed,
 * and the named parameter will not be present in the ArgResult.
 * name: The name that will be assigned to the parameter. */
func (old expectation) AllowNamedParam(name string) (Expectation) {
	old.processors = append(old.processors, func(exp expectation) (err error) {
		return
	})
	return old
}

/* Consumes a single flag from the command-line arguments.
 * The flag must be present, otherwise validation will fail.
 * name: The name of the flag to look for. */
func (old expectation) ExpectFlag(name string) (Expectation) {
	old.processors = append(old.processors, func(exp expectation) (err error) {
		return
	})
	return old
}

/* Consumes a single option and its value from the command-line arguments.
 *
 * If the option is not found, validation will fail.
 * name: The name of the option, as it will appear in the form "--name".
 * alt: The short/single-character name of the option, as it will appear
 * in the form "-n".
 * At least one of the arguments must be non-nil. */
func (old expectation) ExpectOption(name string, alt string) (Expectation) {
	old.processors = append(old.processors, func(exp expectation) (err error) {
		return
	})
	return old
}

/* Consumes the next argument from the command-line as a parameter.
 *
 * If there are no more arguments to consume, validation will fail. */
func (old expectation) ExpectParam() (Expectation) {
	old.processors = append(old.processors, func(exp expectation) (err error) {
		return
	})
	return old
}

/* Consumes the next argument from the command-line as a parameter,
 * giving it the specified name.
 *
 * If there are no more arguments to consume, validation will fail.
 * name: The name that will be assigned to the parameter. */
func (old expectation) ExpectNamedParam(name string) (Expectation) {
	old.processors = append(old.processors, func(exp expectation) (err error) {
		return
	})
	return old
}

/* Discards any remaining, unconsumed arguments, so that they will
 * not cause validation to fail.
 *
 * Alternately, can be used to force the next Allow or Expect to fail. */
func (old expectation) Chop() (Expectation) {
	old.processors = append(old.processors, func(exp expectation) (err error) {
		return
	})
	return old
}

/* Discards any remaining, unconsumed arguments and calls Validate. */
func (final expectation) ChopAndValidate() (result ArgResult, err error) {
	err = fmt.Errorf("Not yet implemented.")
	return
}

/* Called once all expectations have been specified, to parse and
 * validate the arguments. */
func (final expectation) Validate() (result ArgResult, err error) {
	err = fmt.Errorf("Not yet implemented.")
	return
}