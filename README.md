# Args.go #

Command-line argument parsing and validation for Go.

## Arguments ##

**Args.go** recognizes three types of arguments: parameters, options and flags. Flags are arguments of the form "--flag"; options can have the form "--foo" or "-f", and are immediately followed by an option value. A parameter is simply a single argument recognized by its position.

Flags and options will be consumed wherever they are found in the arguments. Parameters are consumed in the order they are included on the command line.

### Rules ###

Argument rules have two forms: `Allow...` and `Expect...`. An *allowance* identifies an argument that MAY be present, while an *expectation* identifies an argument that MUST be present. `Has...` methods are provided to check for the presence of specific arguments. A missing *allowed* argument has no effect on validation, while a missing *expected* argument will cause validation to fail.

#### Flags ####

A flag is an argument of the form "--flag" or "-f". A flag can be accessed by name; its value is a `bool` identifying whether or not the flag was present.

It is an error to attempt to retrieve the value of a flag that was neither *allowed* nor *expected*.

#### Options ####

An option is an argument of the form "--option" or "-o", followed by another argument of any form. An argument WILL NOT be matched as an option if there is no argument following it, even if it has the correct name. Options can be given a primary name and any number of alternate names, and can be accessed by their primary name. The value of an option is the argument immediately following the option.

#### Parameters ####

Parameters are arguments of any form, consumed in the order they are included on the command line. They can be assigned names for convenience, but the assigned name has no impact on whether or not a parameter will be matched. Note that arguments of the form "--name" or "-n" can also be consumed as parameters, though the leading hyphens will not be stripped; the only way a parameter will fail to be matched is if there are no unconsumed arguments remaining.


## Use ##

Call `args.Args([]string)` with a list of command-line arguments (most likely `os.Args`). The result will be an Expectation object, upon which the above argument rules can be attached. All of the `Allow...` and `Expect...` methods are chainable; an example chain invocation could be:

	result, err := args.Args(...).
		ExpectNamedParam("command").
		ExpectParam("target").
		AllowOption("verbose", "v").
		AllowOption("force", "f").
		Validate()
		
Validation will fail unless all arguments are successfully consumed. Use the `Chop()` method (also chainable) to immediately consume all remaining arguments. There is also a `ChopAndValidate()` method to combine these operations.

Arguments are parsed as the argument rules are defined, not when `Validate()` is called. As a result, you can make later rules depend on values previously encountered. For example:

	result, err := args.Args(...).
		ExpectNamedParam("command")
		
	if result.ParamNamed("command") == "process" {
		result, err = result.ExpectParam("target")
	}
	
	result, err = result.
		AllowOption("verbose", "v").
		AllowOption("force", "f").
		Validate()
		
Attempting to access an argument that was not processed will result in a runtime panic. ALWAYS either `Validate()` the result or explicitly check for each argument using the `Has...` methods before attempting to use an argument.