package args

import (
	. "testing"
)

func argsEqual(a1, a2 []string) bool {
	if len(a1) != len(a2) {
		return false
	}

	for i, v := range(a1) {
		if v != a2[i] {
			return false
		}
	}

	return true
}

func TestParam(t *T) {
	args, param, ok := Param([]string{"Foo", "Bar", "-F", "--Foobar"})
	if !argsEqual(args, []string{"Bar", "-F", "--Foobar"}) {
		t.Errorf("The first argument should have been consumed.")
	}
	if param != "Foo" {
		t.Errorf("The first argument should have been returned.")
	}
	if !ok {
		t.Errorf("Should have returned true.")
	}

	args, param, ok = Param(args)
	if !argsEqual(args, []string{"-F", "--Foobar"}) {
		t.Errorf("The first argument should have been consumed (again).")
	}
	if param != "Bar" {
		t.Errorf("The first argument should have been returned.")
	}
	if !ok {
		t.Errorf("Should have returned true.")
	}

	args, param, ok = Param(args)
	if !argsEqual(args, []string{"--Foobar"}) {
		t.Errorf("The first argument should have been consumed (again).")
	}
	if param != "-F" {
		t.Errorf("The first argument should have been returned.")
	}
	if !ok {
		t.Errorf("Should have returned true.")
	}

	args, param, ok = Param(args)
	if !argsEqual(args, []string{}) {
		t.Errorf("The final argument should have been consumed.")
	}
	if param != "--Foobar" {
		t.Errorf("The final argument should have been returned.")
	}
	if !ok {
		t.Errorf("Should have returned true.")
	}

	args, param, ok = Param(args)
	if !argsEqual(args, []string{}) {
		t.Errorf("Arguments should be unchanged.")
	}
	if ok {
		t.Errorf("Should have returned false.")
	}
}

func TestFlag(t *T) {
	args, flag := Flag([]string{"Foo", "--bar", "-f", "what", "--what"}, "bar", "f")
	if !argsEqual(args, []string{"Foo", "-f", "what", "--what"}) {
		t.Errorf("Should have consumed the first matching flag.")
	}
	if !flag {
		t.Errorf("Should have found the flag.")
	}

	args, flag = Flag(args, "w", "what")
	if !argsEqual(args, []string{"Foo", "-f", "what"}) {
		t.Errorf("Should have consumed the flag (with leading hyphens).")
	}
	if !flag {
		t.Errorf("Should have found the flag.")
	}

	args, flag = Flag(args, "w", "what")
	if !argsEqual(args, []string{"Foo", "-f", "what"}) {
		t.Errorf("Should not have consumed anything.")
	}
	if flag {
		t.Errorf("Should not have found the flag.")
	}

	args, flag = Flag([]string{}, "foo")
	if !argsEqual(args, []string{}) {
		t.Errorf("Should not have...invented new arguments?")
	}
	if flag {
		t.Errorf("Should not have found the flag.")
	}
}

func TestOption(t *T) {
	args := []string{"Foo", "--bar", "fizz", "-b", "que", "-w"}

	args, value, ok, err := Option(args, "b", "bar")
	if !argsEqual(args, []string{"Foo", "-b", "que", "-w"}) {
		t.Errorf("Should have consumed the first matching option and its value.")
	}
	if value != "fizz" {
		t.Errorf("Should have gotten the value immediately after the first matching option.")
	}
	if !ok {
		t.Errorf("Should have found the option.")
	}
	if err != nil {
		t.Errorf("Should not have returned an error.")
	}

	args, value, ok, err = Option(args, "Foo", "f")
	if !argsEqual(args, []string{"Foo", "-b", "que", "-w"}) {
		t.Errorf("Should not have consumed anything.")
	}
	if ok {
		t.Errorf("Should not have found the option.")
	}
	if err != nil {
		t.Errorf("Should not have returned an error.")
	}

	args, value, ok, err = Option(args, "w")
	if !argsEqual(args, []string{"Foo", "-b", "que", "-w"}) {
		t.Errorf("Should not have consumed anything.")
	}
	if !ok {
		t.Errorf("Should have found the option.")
	}
	if e, ok := err.(OptionMissingValue); !ok {
		t.Errorf("Should have returned an OptionMissingValue error.")
	} else if string(e) != "-w" {
		t.Errorf("Should have returned an error with the name of the argument found.")
	}
}

func TestOptionInt(t *T) {
	args, val, ok, err := OptionInt([]string{"foo", "42"}, "foo")
	if !argsEqual(args, []string{"foo", "42"}) {
		t.Errorf("Should not have consumed anything.")
	}
	if ok {
		t.Errorf("Should not have found the option.")
	}
	if err != nil {
		t.Errorf("Should not have returned an error.")
	}

	args, val, ok, err = OptionInt([]string{"--foo", "42"}, "foo")
	if !argsEqual(args, []string{}) {
		t.Errorf("Should have consumed the option and its value.")
	}
	if val != 42 {
		t.Errorf("Should have returned the parsed value.")
	}
	if !ok {
		t.Errorf("Should have found the option.")
	}
	if err != nil {
		t.Errorf("Should not have returned an error.")
	}

	args, val, ok, err = OptionInt([]string{"--foo", "-98765"}, "foo")
	if !argsEqual(args, []string{}) {
		t.Errorf("Should have consumed the option and its value.")
	}
	if val != -98765 {
		t.Errorf("Should have returned the parsed value.")
	}
	if !ok {
		t.Errorf("Should have found the option.")
	}
	if err != nil {
		t.Errorf("Should not have returned an error.")
	}

	args, val, ok, err = OptionInt([]string{"--foo", "what"}, "foo")
	if !argsEqual(args, []string{"--foo", "what"}) {
		t.Errorf("Should not have consumed anything.")
	}
	if !ok {
		t.Errorf("Should have found the option.")
	}
	if err == nil {
		t.Errorf("Should have returned an error.")
	}
}