package args

import (
	"fmt"
	"strconv"
)

// Removes and returns the first argument, regardless of its form. Returns a
// bool indicating whether there was any argument to consume.
func Param(args []string) ([]string, string, bool) {
	if len(args) > 0 {
		return args[1:], args[0], true
	} else {
		return args, "", false
	}
}

// Looks for a single argument of the form "--flag" or "-f". Removes only that
// argument, and returns a bool indicating whether it was found.
func Flag(args []string, name string, names...string) ([]string, bool) {
	lookFor := make([]string, 0, len(names)+1)

	for _, name := range(append(names, name)) {
		if len(name) == 0 {
			continue
		}

		if len(name) == 1 {
			lookFor = append(lookFor, fmt.Sprintf("-%s", name))
		} else {
			lookFor = append(lookFor, fmt.Sprintf("--%s", name))
		}
	}

	for i, arg := range(args) {
		for _, name := range(lookFor) {
			if arg == name {
				return append(args[0:i], args[i+1:]...), true
			}
		}
	}

	return args, false
}

// Looks for arguments of the form "--name value" or "-n value". If the option
// name is found with no argument following it, an error will be returned.
func Option(args []string, name string, names...string) ([]string, string, bool, error) {
	lookFor := make([]string, 0, len(names)+1)

	for _, name := range(append(names, name)) {
		if len(name) == 0 {
			continue
		}

		if len(name) == 1 {
			lookFor = append(lookFor, fmt.Sprintf("-%s", name))
		} else {
			lookFor = append(lookFor, fmt.Sprintf("--%s", name))
		}
	}

	for i, arg := range(args) {
		for _, name := range(lookFor) {
			if arg == name {
				if i == len(args) - 1 {
					return args, "", true, OptionMissingValue(name)
				} else {
					val := args[i+1]
					return append(args[0:i], args[i+2:]...), val, true, nil
				}
			}
		}
	}

	return args, "", false, nil
}

// Returns the value of an option (as above) as an integer. Returns an error
// if it could not be parsed.
func OptionInt(args []string, name string, names...string) ([]string, int, bool, error) {
	args2, val, ok, err := Option(args, name, names...)
	if !ok || err != nil {
		return args, 0, ok, err
	}

	i, err := strconv.ParseInt(val, 0, 0)
	if err != nil {
		return args, 0, ok, err
	}

	return args2, int(i), ok, nil
}

type OptionMissingValue string

func (e OptionMissingValue) Error() string {
	return fmt.Sprintf("Option '%s' has no value", string(e))
}