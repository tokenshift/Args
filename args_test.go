package args

import (
	"testing"
)

/* Assertions */

/* Test that the specified flag has the expected value.
 *
 * t: Test context.
 * args: The Expectation object being tested.
 * name: The name of the flag.
 * expected: The expected value of the flag. */
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

/* Test that the specified option is not present.
 *
 * t: Test context.
 * args: The Expectation object being tested.
 * name: The name of the option. */
func assertNoOption(t *testing.T, args Expectation, name string) {
	if args.HasOption(name) {
		t.Errorf("Option named '%v' was not expected.", name)
	}
}

/* Test that the specified option has the expected value.
 *
 * t: Test context.
 * args: The Expectation object being tested.
 * name: The name of the option.
 * expected: The expected value of the option. */
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

/* Test that the specified parameter has the expected value.
 *
 * t: Test context.
 * args: The Expectation object being tested.
 * index: 0-based index of the positional parameter..
 * expected: The expected value of the parameter. */
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

/* Test that the specified parameter has the expected value.
 *
 * t: Test context.
 * args: The Expectation object being tested.
 * name: The name of the parameter.
 * expected: The expected value of the parameter. */
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


/* Tests */

/* Allowances */

/* Flags */

func TestOptionalFlag(t *testing.T) {
	result, err := Args([]string{"--test"}).
		AllowFlag("test").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertFlag(t, result, "test", true)
}

func TestMissingOptionalFlag(t *testing.T) {
	result, err := Args([]string{}).
		AllowFlag("test").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertFlag(t, result, "test", false)
}

/* Options */

func TestMissingAllowedOption(t *testing.T) {
	result, err := Args([]string{}).
		AllowOption("key").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertNoOption(t, result, "key")
}

func TestSingleAllowedOption(t *testing.T) {
	result, err := Args([]string{"--key", "value"}).
		AllowOption("key").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertOption(t, result, "key", "value")
}

func TestSingleAllowedAltOption(t *testing.T) {
	result, err := Args([]string{"--id", "value"}).
		AllowOption("key", "id").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertOption(t, result, "key", "value")
}

func TestSingleAllowedShortAltOption(t *testing.T) {
	result, err := Args([]string{"-k", "value"}).
		AllowOption("key", "id", "k").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertOption(t, result, "key", "value")
}

/* Parameters */

func TestMissingAllowedParameter(t *testing.T) {
	result, err := Args([]string{}).
		AllowParam().
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	if result.HasParamAt(0) {
		t.Error("Did not expect a parameter.")
	}
}

func TestMissingNamedAllowedParameter(t *testing.T) {
	result, err := Args([]string{}).
		AllowNamedParam("test").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	if result.HasParamAt(0) {
		t.Error("Did not expect a parameter at position 0.")
	}

	if result.HasParamNamed("test") {
		t.Error("Did not expect a parameter named 'test'.")
	}
}

func TestSingleAllowedParameter(t *testing.T) {
	result, err := Args([]string{"test"}).
		AllowParam().
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertParamAt(t, result, 0, "test")
}

func TestSingleAllowedNamedParameter(t *testing.T) {
	result, err := Args([]string{"test"}).
		AllowNamedParam("key").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertParamNamed(t, result, "key", "test")
}

/* Expectations */

/* Flags */

func TestSingleFlag(t *testing.T) {
	result, err := Args([]string{"--test"}).
		ExpectFlag("test").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertFlag(t, result, "test", true)
}

func TestSingleMissingFlag(t *testing.T) {
	_, err := Args([]string{}).
		ExpectFlag("test").
		Validate()

	if err == nil {
		t.Error("Validation should have failed; flag was missing.")
		return
	}
}

/* Options */

func TestMissingOption(t *testing.T) {
	_, err := Args([]string{"--key"}).
		ExpectOption("key").
		Validate()

	if err == nil {
		t.Error("Validation should have failed; option 'key' was missing.")
		return
	}
}

func TestSingleOption(t *testing.T) {
	result, err := Args([]string{"--key", "value"}).
		ExpectOption("key").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertOption(t, result, "key", "value")
}

func TestSingleAltOption(t *testing.T) {
	result, err := Args([]string{"--id", "value"}).
		ExpectOption("key", "id").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertOption(t, result, "key", "value")
}

func TestSingleShortAltOption(t *testing.T) {
	result, err := Args([]string{"-k", "value"}).
		ExpectOption("key", "id", "k").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertOption(t, result, "key", "value")
}

/* Parameters */

func TestMissingParameter(t *testing.T) {
	_, err := Args([]string{}).
		ExpectParam().
		Validate()

	if err == nil {
		t.Error("Validation should have failed; missing a positional parameter.")
		return
	}
}

func TestSingleParameter(t *testing.T) {
	result, err := Args([]string{"test"}).
		ExpectParam().
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertParamAt(t, result, 0, "test")
}

func TestSingleNamedParameter(t *testing.T) {
	result, err := Args([]string{"test"}).
		ExpectNamedParam("key").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertParamNamed(t, result, "key", "test")
}

/* Chop */

func TestChoppingArguments(t *testing.T) {
	result, err := Args([]string{
		"do",
		"something",
		"--carefully",
		"--and",
		"slowly",
		"-v",
		"(that means verbose)"}).
		ExpectParam().
		Chop().
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertParamAt(t, result, 0, "do")
}

func TestChopAndValidate(t *testing.T) {
	result, err := Args([]string{
		"do",
		"something",
		"--carefully",
		"--and",
		"slowly",
		"-v",
		"(that means verbose)"}).
		ExpectParam().
		ChopAndValidate()

	if err != nil {
		t.Error(err)
		return
	}

	assertParamAt(t, result, 0, "do")
}

/* Additional Tests */

func TestMultipleArguments(t *testing.T) {
	result, err := Args([]string{
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
