package args

import "testing"

// {command} [...arguments...]
func TestCommandAndRemainder(t *testing.T) {
	result := Load([]string{"process", "something", "goes", "here"}).
		ExpectParamNamed("command")

	assertParamNamed(t, result, "command", "process")

	result, tail := result.ChopString()
	assertParamNamed(t, result, "command", "process")
	assertStringEquals(t, "something goes here", tail)
}

// Each expectation should return a new Args object.
func TestPureFunctionalChaining(t *testing.T) {
	args := Load([]string { "foo", "bar" })
	assertBoolEquals(t, false, args.HasParamNamed("fizzbuzz"))

	args2 := args.ExpectParamNamed("fizzbuzz")
	assertBoolEquals(t, false, args.HasParamNamed("fizzbuzz"))
	assertBoolEquals(t, true, args2.HasParamNamed("fizzbuzz"))
}
