package goerr_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/karrick/goerr"
)

func ExampleExitCode() {
	err := goerr.New("some error message").WithExitCode(13)

	fmt.Println(err)
	fmt.Println(err.ExitCode())
	// Output:
	// some error message
	// 13
}

func ExampleMultiLine() {
	// Lines are printed in the order in which they were added.
	err := goerr.New("cannot do thing").
		WithLineAfterOptions("line 2").
		WithLinesAfterOptions([]string{"line 3", "line 4"}).
		WithLineBetweenMessageAndOption("line 1").
		WithLineAfterOptions("line 5").
		WithLineBeforeMessage("line 0")

	fmt.Println(err)
	// Output:
	// line 0
	// cannot do thing
	// line 1
	// line 2
	// line 3
	// line 4
	// line 5
}

func ExampleWithOptionComments() {
	args := []string{"zero", "one", "--two", "three"}

	// NOTE: Lines are added in the order given for each section:
	//
	// 1. lines before the error message;
	// 2. error message line;
	// 3. lines between error message and any option lines;
	// 4. option lines;
	// 5. lines after options.

	err := goerr.New("This is the error message.").
		WithLineBeforeMessage("Optional lines before the error message.").
		WithLineBetweenMessageAndOption("Some lines between message and options.").
		WithOptions(args).
		WithOptionComment(1, "for this sub-command").
		WithOptionComment(2, "for this option").
		WithOptionComment(3, "cannot find this file").
		WithLineAfterOptions("Zero or more additional").
		WithLineAfterOptions("lines of information.").
		WithExitCode(13)

	fmt.Println(err)
	fmt.Println(err.ExitCode())
	// Output:
	// Optional lines before the error message.
	// This is the error message.
	// Some lines between message and options.
	// zero one --two three
	//                ^~~~~ cannot find this file
	//          ^~~~~ for this option
	//      ^~~ for this sub-command
	// Zero or more additional
	// lines of information.
	// 13
}

func parseIntegerOption(s string) (int, error) {
	n, err := strconv.Atoi(s)

	// NOTE: MaybeWrapf will return a nil error if given an error, but wrap err
	// such that other attributes may be attached to it if err is non-nil.

	return n, goerr.Wrapf(err, "cannot parse option").
		WithExitCode(2)
}

func ExampleMaybeWrapNonNilError() {
	value, err := parseIntegerOption("123")

	// NOTE: Because MaybeWrap returns nil when it is given a nil error to
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

func TestError(t *testing.T) {
	t.Run("init", func(t *testing.T) {
		t.Run("MaybeWrap", func(t *testing.T) {
			t.Run("sans error", func(t *testing.T) {
				t.Run("with message", func(t *testing.T) {
					// Wrapping a nil error returns a nil Error, however the
					// mutating methods can still be invoked on the nil Error
					// instance.
					var nilErr error

					ee := goerr.Wrapf(nilErr, "cannot configure: %v", "bad-data").
						WithExitCode(13).
						WithTemporary(true)

					if got, want := ee, (*goerr.Error)(nil); got != want {
						t.Errorf("GOT: %v; WANT: %v", got, want)
					}
					if got, want := goerr.ExitCode(ee), 0; got != want {
						t.Errorf("GOT: %v; WANT: %v", got, want)
					}
					if got, want := goerr.Temporary(ee), false; got != want {
						t.Errorf("GOT: %v; WANT: %v", got, want)
					}
				})
			})

			t.Run("with error", func(t *testing.T) {
				t.Run("sans message", func(t *testing.T) {
					err := fmt.Errorf("cannot parse int: %q", "123abc")

					err = goerr.Wrap(err).
						WithExitCode(13).
						WithTemporary(true)

					if got, want := err.Error(), "cannot parse int: \"123abc\""; got != want {
						t.Errorf("GOT: %q; WANT: %q", got, want)
					}

					ee, ok := err.(*goerr.Error)
					if !ok {
						t.Fatalf("GOT: %T; WANT: *goerr.Error", err)
					}
					if got, want := ee.Unwrap().Error(), "cannot parse int: \"123abc\""; got != want {
						t.Errorf("GOT: %q; WANT: %q", got, want)
					}
					if got, want := goerr.ExitCode(ee), 13; got != want {
						t.Errorf("GOT: %v; WANT: %v", got, want)
					}
					if got, want := goerr.Temporary(ee), true; got != want {
						t.Errorf("GOT: %v; WANT: %v", got, want)
					}
				})

				t.Run("with message", func(t *testing.T) {
					err := fmt.Errorf("cannot parse int: %q", "123abc")

					err = goerr.Wrapf(err, "cannot configure").
						WithExitCode(13).
						WithTemporary(true)

					if got, want := err.Error(), "cannot configure: cannot parse int: \"123abc\""; got != want {
						t.Errorf("GOT: %q; WANT: %q", got, want)
					}

					ee, ok := err.(*goerr.Error)
					if !ok {
						t.Fatalf("GOT: %T; WANT: *goerr.Error", err)
					}
					if got, want := ee.Unwrap().Error(), "cannot parse int: \"123abc\""; got != want {
						t.Errorf("GOT: %q; WANT: %q", got, want)
					}
					if got, want := goerr.ExitCode(ee), 13; got != want {
						t.Errorf("GOT: %v; WANT: %v", got, want)
					}
					if got, want := goerr.Temporary(ee), true; got != want {
						t.Errorf("GOT: %v; WANT: %v", got, want)
					}
				})

			})
		})

		t.Run("New", func(t *testing.T) {
			ee := goerr.New("cannot parse int: %q", "123abc")

			if got, want := ee.Error(), "cannot parse int: \"123abc\""; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}

			if got, want := ee.Unwrap(), error(nil); got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}

			// NOTE: Defaults can be overridden after the fact.

			if got, want := ee.ExitCode(), 0; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}

			ee.WithExitCode(13)

			if got, want := ee.ExitCode(), 13; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}

			ee.WithExitCode(42)

			if got, want := ee.ExitCode(), 42; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}

			ee.WithExitCode(0)

			if got, want := ee.ExitCode(), 0; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}

			if got, want := ee.Temporary(), false; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}

			ee.WithTemporary(true)

			if got, want := ee.Temporary(), true; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}

			ee.WithTemporary(false)

			if got, want := ee.Temporary(), false; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("sans error", func(t *testing.T) {
			t.Run("sans message", func(t *testing.T) {
				var ee goerr.Error

				if got, want := ee.Error(), "error without message or wrapped error"; got != want {
					t.Errorf("GOT: %q; WANT: %q", got, want)
				}
			})

			t.Run("with message", func(t *testing.T) {
				ee := goerr.New("foo: %v", "bar")

				if got, want := ee.Error(), "foo: bar"; got != want {
					t.Errorf("GOT: %q; WANT: %q", got, want)
				}
			})
		})

		t.Run("with error", func(t *testing.T) {
			t.Run("sans message", func(t *testing.T) {
				ee := goerr.Wrap(fmt.Errorf("foo: %v", "bar"))

				if got, want := ee.Error(), "foo: bar"; got != want {
					t.Errorf("GOT: %q; WANT: %q", got, want)
				}
			})

			t.Run("with message", func(t *testing.T) {
				ee := goerr.New("foo: %v", "bar")

				if got, want := ee.Error(), "foo: bar"; got != want {
					t.Errorf("GOT: %q; WANT: %q", got, want)
				}
			})
		})

		t.Run("with options sans comments", func(t *testing.T) {
			err := goerr.New("some error message").
				WithOptions([]string{"zero", "one", "--two", "three"})

			lines := err.ErrorLines()
			if got, want := len(lines), 2; got != want {
				t.Fatalf("GOT: %v; WANT: %v", got, want)
			}
			if got, want := lines[0], "some error message"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
			if got, want := lines[1], "zero one --two three"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
		})

		t.Run("with options with comments", func(t *testing.T) {
			err := goerr.New("some error message").
				WithLineBeforeMessage("optional lines before the error message").
				WithLineBetweenMessageAndOption("Some lines between message and options").
				WithOptions([]string{"zero", "one", "--two", "three"}).
				WithOptionComment(1, "cannot do sub-command").
				WithOptionComment(3, "invalid argument").
				WithOptionComment(2, "for this option").
				WithLineAfterOptions("Zero or more additional lines")

			lines := err.ErrorLines()
			if got, want := len(lines), 8; got != want {
				t.Fatalf("GOT: %v; WANT: %v", got, want)
			}
			if got, want := lines[0], "optional lines before the error message"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
			if got, want := lines[1], "some error message"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
			if got, want := lines[2], "Some lines between message and options"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
			if got, want := lines[3], "zero one --two three"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
			if got, want := lines[4], "               ^~~~~ invalid argument"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
			if got, want := lines[5], "         ^~~~~ for this option"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
			if got, want := lines[6], "     ^~~ cannot do sub-command"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
			if got, want := lines[7], "Zero or more additional lines"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
		})
	})

	t.Run("ExitCode", func(t *testing.T) {
		t.Run("sans exit code", func(t *testing.T) {
			var ee goerr.Error

			if got, want := ee.ExitCode(), 0; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})

		t.Run("with exit code", func(t *testing.T) {
			ee := new(goerr.Error)

			ee = ee.WithExitCode(13)

			if got, want := ee.ExitCode(), 13; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})
	})

	t.Run("Unwrap", func(t *testing.T) {
		t.Run("sans wrapped error", func(t *testing.T) {
			var ee goerr.Error

			if got, want := ee.Unwrap(), error(nil); got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})

		t.Run("with wrapped error", func(t *testing.T) {
			ee := goerr.Wrap(fmt.Errorf("foo: %v", "bar"))

			if got, want := ee.Unwrap().Error(), "foo: bar"; got != want {
				t.Errorf("GOT: %v; WANT: %v", got, want)
			}
		})
	})
}
