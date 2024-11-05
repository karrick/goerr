package goerr_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/karrick/goerr"
)

type dummyExitCoder struct{ code int }

func (dec dummyExitCoder) Error() string {
	return fmt.Sprintf("returns exit code: %d", dec.code)
}

func (dec dummyExitCoder) ExitCode() int { return dec.code }

type dummyTemporaryer struct{ temporary bool }

func (dec dummyTemporaryer) Error() string {
	return fmt.Sprintf("returns temporary: %t", dec.temporary)
}

func (dec dummyTemporaryer) Temporary() bool { return dec.temporary }

type dummyUnwrapper struct{ err error }

func (dec dummyUnwrapper) Error() string {
	return fmt.Sprintf("unwraps err: %v", dec.err)
}

func (dec dummyUnwrapper) Unwrap() error { return dec.err }

func TestExitCode(t *testing.T) {
	t.Run("err nil", func(t *testing.T) {
		var err error

		if got, want := goerr.ExitCode(err), 0; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
	})

	t.Run("err *Error nil", func(t *testing.T) {
		var err *goerr.Error

		if got, want := goerr.ExitCode(err), 0; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
	})

	t.Run("err *Error sans exit code", func(t *testing.T) {
		err := goerr.New("some error")

		if got, want := goerr.ExitCode(err), 0; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
	})

	t.Run("err *Error with exit code", func(t *testing.T) {
		err := goerr.New("some error").WithExitCode(42)

		if got, want := goerr.ExitCode(err), 42; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
	})

	t.Run("err exitCoderer", func(t *testing.T) {
		err := &dummyExitCoder{code: 42}

		if got, want := goerr.ExitCode(err), 42; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
	})

	t.Run("err unwrapper nil", func(t *testing.T) {
		err := &dummyUnwrapper{}

		if got, want := goerr.ExitCode(err), 0; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
	})

	t.Run("err unwrapper exitCoderer", func(t *testing.T) {
		err := &dummyUnwrapper{err: &dummyExitCoder{code: 42}}

		if got, want := goerr.ExitCode(err), 42; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
	})

	t.Run("err default", func(t *testing.T) {
		err := errors.New("no exit code no unwrap")

		if got, want := goerr.ExitCode(err), 0; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
	})
}

func TestTemporary(t *testing.T) {
	t.Run("err nil", func(t *testing.T) {
		var err error

		if got, want := goerr.Temporary(err), false; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
	})

	t.Run("err *Error nil", func(t *testing.T) {
		var err *goerr.Error

		if got, want := goerr.Temporary(err), false; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
	})

	t.Run("err *Error sans temporary", func(t *testing.T) {
		err := goerr.New("some error")

		if got, want := goerr.Temporary(err), false; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
	})

	t.Run("err *Error with temporary false", func(t *testing.T) {
		err := goerr.New("some error").WithTemporary(false)

		if got, want := goerr.Temporary(err), false; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
	})

	t.Run("err *Error with temporary true", func(t *testing.T) {
		err := goerr.New("some error").WithTemporary(true)

		if got, want := goerr.Temporary(err), true; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
	})

	t.Run("err temporaryer false", func(t *testing.T) {
		err := &dummyTemporaryer{temporary: false}

		if got, want := goerr.Temporary(err), false; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
	})

	t.Run("err temporaryer true", func(t *testing.T) {
		err := &dummyTemporaryer{temporary: true}

		if got, want := goerr.Temporary(err), true; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
	})

	t.Run("err unwrapper nil", func(t *testing.T) {
		err := &dummyUnwrapper{}

		if got, want := goerr.Temporary(err), false; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
	})

	t.Run("err unwrapper temporaryer", func(t *testing.T) {
		err := &dummyUnwrapper{err: &dummyTemporaryer{temporary: true}}

		if got, want := goerr.Temporary(err), true; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
	})

	t.Run("err default", func(t *testing.T) {
		err := fmt.Errorf("no exit code no unwrap")

		if got, want := goerr.Temporary(err), false; got != want {
			t.Errorf("GOT: %v; WANT: %v", got, want)
		}
	})
}
