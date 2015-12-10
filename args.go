package args

import (
	"fmt"
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
func Flag(args []string, names...string) ([]string, bool) {
	for _, name := range(names) {
		if len(name) == 0 {
			continue
		}

		var lookFor string
		if len(name) == 1 {
			lookFor = fmt.Sprintf("-%s", name)
		} else {
			lookFor = fmt.Sprintf("--%s", name)
		}

		for i, arg := range(args) {
			if arg == lookFor {
				return append(args[0:i], args[i+1:]...), true
			}
		}
	}

	return args, false
}

// Looks for arguments of the form "--name value" or "-n value". If the option
// name is found with no argument following it, an error will be returned.
func Option(args []string, names...string) ([]string, string, bool, error) {
	lookFor := make([]string, 0, len(names))

	for _, name := range(names) {
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

type OptionMissingValue string

func (e OptionMissingValue) Error() string {
	return fmt.Sprintf("Option '%s' has no value", string(e))
}