# goerr

Go library to build and wrap errors.

## Description

Go library that provides an Error structure which can contain exit code,
temporary error information, and provide for multi-line contextual output in
error messages.

Errors can optionally contain the desired exit code for a program if the error
is not handled by calling code. In other words, if cannot parse a command line
option is

## Examples

### Displaying Command Line Usage Errors

This makes it easy to print multi-line error messages that show which command
line options have an error and why.

```Go
func ExampleUsage() {
	args := []string{"zero", "--one", "two", "three"}

	err := New("cannot do thing").
		WithExitCode(13).
		WithOptions(args).
		WithOptionComment(3, "cannot do this operation").
		WithOptionComment(2, "cannot find this file").
		WithOptionComment(1, "for this option")

	fmt.Println(err)
	fmt.Println(err.ExitCode())
	fmt.Println(err.Temporary())
	// Output:
	// cannot do thing
	// zero --one two three
	//                ^~~~~ cannot do this operation
	//            ^~~ cannot find this file
	//      ^~~~~ for this option
	// 13
	// false
}
```

### Wrap and Wrapf

The `Wrap` and `Wrapf` functions can be used to wrap error return values,
optionally providing additional context about the error.

```Go
func parseIntegerOption(s string) (int, error) {
	n, err := strconv.Atoi(s)

	// NOTE: Wrap and Wrapf will return a nil error if given an error, but
	// wrap err such that other attributes may be attached to it if err is
	// non-nil.

	return n, goerr.Wrapf(err, "cannot parse option").
		WithExitCode(2)
}

func ExampleMaybeWrapNonNilError() {
	value, err := parseIntegerOption("123")

	// NOTE: Because Wrap and Wrapf return nil when it is given a nil error to
	// wrap, and the nil error has exit code value of 0 and temporary value of
	// false, the output below will be 123, <nil>, and 0.

	fmt.Println(value)
	fmt.Println(err)
	fmt.Println(goerr.ExitCode(err))
	// Output:
	// 123
	// <nil>
	// 0
}

func ExampleMaybeWrapNilError() {
	value, err := parseIntegerOption("123abc")

	fmt.Println(value)
	fmt.Println(err)
	fmt.Println(goerr.ExitCode(err))
	// Output:
	// 0
	// cannot parse option: strconv.Atoi: parsing "123abc": invalid syntax
	// 2
}
```
