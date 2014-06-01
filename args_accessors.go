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
func (final argv) HasFlag(name string) bool {
	_, present := final.flags[name]
	return present
}

// Checks whether the named option was found.
//
// Use this before calling Option on an Allowed
// (not Expected) option.
// name: The name of the option.
func (final argv) HasOption(name string) bool {
	_, present := final.options[name]
	return present
}

// Checks whether there is a parameter at
// the specified index.
// i: The 0-based index of the parameter. 
func (final argv) HasParamAt(i int) bool {
	return len(final.parameters) > i
}

// Checks whether there is a parameter with
// the specified name.
// i: The name of the parameter. 
func (final argv) HasParamNamed(name string) bool {
	_, present := final.namedParameters[name]
	return present
}

// Gets whether the named flag was set.
// name: The name of the flag to check. 
func (final argv) Flag(name string) bool {
	value, ok := final.flags[name]

	if ok {
		return value
	} else {
		panic(fmt.Errorf("You must explicitly Expect or Allow the flag '%v'.", name))
	}
}

// Gets the value of the named option.
// name: The name of the option. 
func (final argv) Option(name string) string {
	value, present := final.options[name]

	if present {
		return value
	} else {
		panic(fmt.Errorf("Option '%v' was not found.", name))
	}
}

// Gets the value of the parameter at the specified position.
// i: The 0-based index of the parameter. 
func (final argv) ParamAt(index int) string {
	if len(final.parameters) > index {
		return final.parameters[index]
	} else {
		panic(fmt.Errorf("No parameter present at index %v.", index))
	}
}

// Gets the value of the named parameter.
// name: The name of the parameter. 
func (final argv) ParamNamed(name string) string {
	index, found := final.namedParameters[name]
	if found {
		return final.parameters[index]
	} else {
		panic(fmt.Errorf("No parameter present with name %v.", name))
	}
}
