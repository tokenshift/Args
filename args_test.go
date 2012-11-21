package args

import (
	"testing"
)

func TestDoAddWithInit(t *testing.T) {
	result, err := Args([]string {"do", "add", "\"this is a test\"", "--init"}).
		ExpectParam().
		ExpectNamedParam("command").
		ExpectNamedParam("title").
		AllowFlag("init").
		Validate()

	if err != nil {
		t.Error(err)
	}

	if val, err := result.ParamAt(0); err != nil || val != "do" {
		t.Error("Expected 'do' as first param.")
	}

	if val, err := result.ParamAt(1); err != nil || val != "add" {
		t.Error("Expected 'add' as second param.")
	}

	if val, err := result.ParamNamed("command"); err != nil || val != "add" {
		t.Error("Expected 'add' as command.")
	}

	if val, err := result.ParamAt(2); err != nil || val != "\"this is a test\"" {
		t.Error("Expected \"this is a test\" as third param.")
	}

	if val, err := result.ParamNamed("title"); err != nil || val != "\"this is a test\"" {
		t.Error("Expected \"this is a test\" as title.")
	}

	if val, err := result.Flag("init"); err != nil || !val {
		t.Error("Expected 'init' to be true.")
	}
}