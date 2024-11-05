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

```go
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
