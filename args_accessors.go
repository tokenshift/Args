package args

import (
	"fmt"
)

// Checks whether the named flag was checked.
//
// NOTE: Does not check whether the flag was present
// in the arguments; checks only whether it was
// Expected or Allowed.
// name: The name of the flag. 
func (final expectation) HasFlag(name string) (present bool) {
	_, present = final.flags[name]
	return
}

// Checks whether the named option was found.
//
// Use this before calling Option on an Allowed
// (not Expected) option.
// name: The name of the option.
func (final expectation) HasOption(name string) (present bool) {
	_, present = final.options[name]
	return
}

// Checks whether there is a parameter at
// the specified index.
// i: The 0-based index of the parameter. 
func (final expectation) HasParamAt(i int) (present bool) {
	present = len(final.parameters) > i
	return
}

// Checks whether there is a parameter with
// the specified name.
// i: The name of the parameter. 
func (final expectation) HasParamNamed(name string) (present bool) {
	_, present = final.namedParameters[name]
	return
}

// Gets whether the named flag was set.
// name: The name of the flag to check. 
func (final expectation) Flag(name string) (value bool, err error) {
	value, ok := final.flags[name]

	if !ok {
		err = fmt.Errorf("You must explicitly Expect or Allow the flag '%v'.", name)
	}

	return
}

// Gets the value of the named option.
// name: The name of the option. 
func (final expectation) Option(name string) (value string, err error) {
	value, present := final.options[name]

	if !present {
		err = fmt.Errorf("Option '%v' was not found.", name)
	}

	return
}

// Gets the value of the parameter at the specified position.
// i: The 0-based index of the parameter. 
func (final expectation) ParamAt(index int) (value string, err error) {
	if len(final.parameters) > index {
		value = final.parameters[index]
	} else {
		err = fmt.Errorf("No parameter present at index %v.", index)
	}

	return
}

// Gets the value of the named parameter.
// name: The name of the parameter. 
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
