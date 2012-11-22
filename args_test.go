package args

import (
	"testing"
)

func TestSingleFlag(t *testing.T) {
	t.Error("Test not implemented.")
}

func TestSingleMissingFlag(t *testing.T) {
	t.Error("Test not implemented.")
}

func TestSingleOption(t *testing.T) {
	t.Error("Test not implemented.")
}

func TestSingleAltOption(t *testing.T) {
	t.Error("Test not implemented.")
}

func TestSingleShortAltOption(t *testing.T) {
	t.Error("Test not implemented.")
}

func assertFlag(t *testing.T, args Expectation, name string, expected bool) {
	val, err := args.Flag(name)

	if err != nil {
		t.Error(err)
		return
	}

	if val != expected {
		if expected {
			t.Errorf("Expected flag '%v'.", name)
		} else {
			t.Errorf("Did not expect flag '%v'.", name)
		}
	}
}

func assertOption(t *testing.T, args Expectation, name string, expected string) {
	val, err := args.Option(name)

	if err != nil {
		t.Error(err)
		return
	}

	if val != expected {
		t.Errorf("Expected option named '%v' to be '%v'.", name, expected)
	}
}

func assertParamAt(t *testing.T, args Expectation, index int, expected string) {
	val, err := args.ParamAt(index)

	if err != nil {
		t.Error(err)
		return
	}

	if val != expected {
		t.Errorf("Expected parameter at position %v to be '%v'.", index, expected)
	}
}

func assertParamNamed(t *testing.T, args Expectation, name string, expected string) {
	val, err := args.ParamNamed(name)

	if err != nil {
		t.Error(err)
		return
	}

	if val != expected {
		t.Errorf("Expected parameter named '%v' to be '%v'.", name, expected)
	}
}

func TestSingleParameter(t *testing.T) {
	result, err := Args([]string {"test"}).
		ExpectParam().
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertParamAt(t, result, 0, "test")
}

func TestSingleNamedParameter(t *testing.T) {
	result, err := Args([]string {"test"}).
		ExpectNamedParam("key").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertParamNamed(t, result, "key", "test")
}

func TestMultipleArguments(t *testing.T) {
	result, err := Args([]string {
		"do",
		"something",
		"--carefully",
		"--and",
		"slowly",
		"-v",
		"(that means verbose)"}).
		ExpectParam().
		ExpectNamedParam("command").
		ExpectFlag("carefully").
		ExpectOption("and").
		ExpectOption("verbose", "v").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertParamAt(t, result, 0, "do")
	assertParamNamed(t, result, "command", "something")
	assertFlag(t, result, "carefully", true)
	assertOption(t, result, "and", "slowly")
	assertOption(t, result, "verbose", "(that means verbose)")
}

// TODO: Add negative tests.