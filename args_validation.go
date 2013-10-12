package args

import (
	"fmt"
	"strings"
)

// Discards any remaining, unconsumed arguments, so that they will
// not cause validation to fail.
//
// Alternately, can be used to force the next Allow or Expect to fail. 
func (chain argv) Chop() Args {
	for i, _ := range chain.consumed {
		chain.consumed[i] = true
	}
	return chain
}

// Consumes and returns any remaining, unconsumed arguments, so that
// they will not cause validation to fail. The remaining arguments are
// returned as a slice of strings.
//
// Alternately, can be used to force the next Allow or Expect to fail.
func (chain argv) ChopSlice() (result Args, tail []string) {
	count := 0
	for _, consumed := range chain.consumed {
		if !consumed {
			count++
		}
	}

	tail = make([]string, 0, count)

	for i, consumed := range chain.consumed {
		if !consumed {
			chain.consumed[i] = true
			tail = append(tail, chain.args[i])
		}
	}

	result = chain
	return
}

// Consumes and returns any remaining, unconsumed arguments, so that
// they will not cause validation to fail. The remaining arguments are
// returned as a single string, concatenated with spaces.
//
// Alternately, can be used to force the next Allow or Expect to fail.
func (chain argv) ChopString() (result Args, tail string) {
	result, sTail := chain.ChopSlice()
	tail = strings.Join(sTail, " ")
	return
}

// Discards any remaining, unconsumed arguments and calls Validate. 
func (final argv) ChopAndValidate() (result Args, err error) {
	result, err = final.Chop().Validate()
	return
}

// Called once all argvs have been specified, to parse and
// validate the arguments.
//
// Validation will fail if:
// 1. There are unconsumed arguments remaining.
//    Call Chop() to consume any remaining arguments.
// 2. Any Expected arguments were not found.
//    Allowed arguments will not cause validation errors when missing. 
func (final argv) Validate() (result Args, err error) {
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

// Error type to represent failed argument validation.
//
// Errors: The list of individual validation errors
// that were encountered. 
type ArgsError struct {
	Errors []error
}

// Display string for ArgsError.
//
// Displays the list of validation errors. 
func (argsError ArgsError) Error() string {
	return fmt.Sprintf("Validation failed: %v", argsError.Errors)
}
