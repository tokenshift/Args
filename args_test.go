package args

import "strings"
import "testing"

// Assertions 

// Tests that two booleans are equal.
func assertBoolEquals(t *testing.T, expected bool, actual bool) bool {
	if expected != actual {
		t.Errorf("Expected %v, got %v.", expected, actual)
		return false
	}
	return true
}

// Test that the specified flag has the expected value.
//
// t: Test context.
// args: The Args object being tested.
// name: The name of the flag.
// expected: The expected value of the flag. 
func assertFlag(t *testing.T, args Args, name string, expected bool) {
	val, err:= args.Flag(name)

	if err != nil {
		t.Error(err)
	}

	if val != expected {
		if expected {
			t.Errorf("Expected flag '%v'.", name)
		} else {
			t.Errorf("Did not expect flag '%v'.", name)
		}
	}
}

// Test that the specified option is not present.
//
// t: Test context.
// args: The Args object being tested.
// name: The name of the option. 
func assertNoOption(t *testing.T, args Args, name string) {
	if args.HasOption(name) {
		t.Errorf("Option named '%v' was not expected.", name)
	}
}

// Test that the specified option has the expected value.
//
// t: Test context.
// args: The Args object being tested.
// name: The name of the option.
// expected: The expected value of the option. 
func assertOption(t *testing.T, args Args, name string, expected string) {
	val, err := args.Option(name)

	if err != nil {
		t.Error(err)
	}

	if val != expected {
		t.Errorf("Expected option named '%v' to be '%v'; it was '%v'.", name, expected, val)
	}
}

// Test that the specified parameter has the expected value.
//
// t: Test context.
// args: The Args object being tested.
// index: 0-based index of the positional parameter..
// expected: The expected value of the parameter. 
func assertParamAt(t *testing.T, args Args, index int, expected string) {
	val, err := args.ParamAt(index)

	if err != nil {
		t.Error(err)
	}

	if val != expected {
		t.Errorf("Expected parameter at position %v to be '%v'.", index, expected)
	}
}

// Test that the specified parameter has the expected value.
//
// t: Test context.
// args: The Args object being tested.
// name: The name of the parameter.
// expected: The expected value of the parameter. 
func assertParamNamed(t *testing.T, args Args, name string, expected string) {
	val, err := args.ParamNamed(name)

	if err != nil {
		t.Error(err)
	}

	if val != expected {
		t.Errorf("Expected parameter named '%v' to be '%v'.", name, expected)
	}
}

func assertStringEquals(t *testing.T, expected, actual string) bool {
	if expected != actual {
		t.Errorf("Expected '%s', got '%s'.", expected, actual)
		return false
	}
	return true
}

// Tests 

// Flags 

func TestOptionalFlag(t *testing.T) {
	result, err := Load([]string{"--test"}).
		AllowFlag("test").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertFlag(t, result, "test", true)
}

func TestMissingOptionalFlag(t *testing.T) {
	result, err := Load([]string{}).
		AllowFlag("test").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertFlag(t, result, "test", false)
}

// Options 

func TestMissingAllowedOption(t *testing.T) {
	result, err := Load([]string{}).
		AllowOption("key").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertNoOption(t, result, "key")
}

func TestSingleAllowedOption(t *testing.T) {
	result, err := Load([]string{"--key", "value"}).
		AllowOption("key").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertOption(t, result, "key", "value")
}

func TestSingleAllowedAltOption(t *testing.T) {
	result, err := Load([]string{"--id", "value"}).
		AllowOption("key", "id").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertOption(t, result, "key", "value")
}

func TestSingleAllowedShortAltOption(t *testing.T) {
	result, err := Load([]string{"-k", "value"}).
		AllowOption("key", "id", "k").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertOption(t, result, "key", "value")
}

// Parameters 

func TestMissingAllowedParameter(t *testing.T) {
	result, err := Load([]string{}).
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
	result, err := Load([]string{}).
		AllowParamNamed("test").
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
	result, err := Load([]string{"test"}).
		AllowParam().
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertParamAt(t, result, 0, "test")
}

func TestSingleAllowedNamedParameter(t *testing.T) {
	result, err := Load([]string{"test"}).
		AllowParamNamed("key").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertParamNamed(t, result, "key", "test")
}

// Argss 

// Flags 

func TestSingleFlag(t *testing.T) {
	result, err := Load([]string{"--test"}).
		ExpectFlag("test").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertFlag(t, result, "test", true)
}

func TestSingleMissingFlag(t *testing.T) {
	_, err := Load([]string{}).
		ExpectFlag("test").
		Validate()

	if err == nil {
		t.Error("Validation should have failed; flag was missing.")
		return
	}
}

// Options 

func TestMissingOption(t *testing.T) {
	_, err := Load([]string{"--key"}).
		ExpectOption("key").
		Validate()

	if err == nil {
		t.Error("Validation should have failed; option 'key' was missing.")
		return
	}
}

func TestSingleOption(t *testing.T) {
	result, err := Load([]string{"--key", "value"}).
		ExpectOption("key").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertOption(t, result, "key", "value")
}

func TestSingleAltOption(t *testing.T) {
	result, err := Load([]string{"--id", "value"}).
		ExpectOption("key", "id").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertOption(t, result, "key", "value")
}

func TestSingleShortAltOption(t *testing.T) {
	result, err := Load([]string{"-k", "value"}).
		ExpectOption("key", "id", "k").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertOption(t, result, "key", "value")
}

// Parameters 

func TestMissingParameter(t *testing.T) {
	_, err := Load([]string{}).
		ExpectParam().
		Validate()

	if err == nil {
		t.Error("Validation should have failed; missing a positional parameter.")
		return
	}
}

func TestSingleParameter(t *testing.T) {
	result, err := Load([]string{"test"}).
		ExpectParam().
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertParamAt(t, result, 0, "test")
}

func TestSingleNamedParameter(t *testing.T) {
	result, err := Load([]string{"test"}).
		ExpectParamNamed("key").
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertParamNamed(t, result, "key", "test")
}

// Chop

func TestChoppingArguments(t *testing.T) {
	result, err := Load([]string{
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

func TestChopSlice(t *testing.T) {
	result, tail := Load([]string{
		"do",
		"something",
		"--carefully",
		"--and",
		"slowly",
		"-v",
		"(that means verbose)"}).
		ExpectParam().
		ChopSlice()

	result, err := result.
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertParamAt(t, result, 0, "do")

	if strings.Join(tail, " ") != "something --carefully --and slowly -v (that means verbose)" {
		t.Error("Remaining (chopped) arguments should have been returned as a slice.")
	}
}

func TestChopSliceWithOptionsAfterwards(t *testing.T) {
	result, tail := Load([]string{
		"add",
		"feature",
		"this",
		"is",
		"a",
		"test",
		"--priority",
		"high"}).
		ExpectParam().
		AllowFlag("init").
		AllowOption("priority", "p").
		ChopSlice()

	result, err := result.
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertParamAt(t, result, 0, "add")
	assertFlag(t, result, "init", false)
	assertOption(t, result, "priority", "high")

	if strings.Join(tail, " ") != "feature this is a test" {
		t.Errorf("Expected \"%v\", got \"%v\"",
			"feature this is a test",
			tail)
	}
}

func TestChopString(t *testing.T) {
	result, tail := Load([]string{
		"do",
		"something",
		"--carefully",
		"--and",
		"slowly",
		"-v",
		"(that means verbose)"}).
		ExpectParam().
		ChopString()

	result, err := result.
		Validate()

	if err != nil {
		t.Error(err)
		return
	}

	assertParamAt(t, result, 0, "do")

	if tail != "something --carefully --and slowly -v (that means verbose)" {
		t.Error("Remaining (chopped) arguments should have been returned as a string.")
	}
}

func TestChopAndValidate(t *testing.T) {
	result, err := Load([]string{
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

// Additional Tests 

func TestMultipleArguments(t *testing.T) {
	result, err := Load([]string{
		"do",
		"something",
		"--carefully",
		"--and",
		"slowly",
		"-v",
		"(that means verbose)"}).
		ExpectParam().
		ExpectParamNamed("command").
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
