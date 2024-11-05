package goerr

type exitCoder interface{ ExitCode() int }

type temporaryer interface{ Temporary() bool }

type unwrapper interface{ Unwrap() error }

// ExitCode returns the result of invoking the ExitCode method for err or the
// first wrapped error recursing until an error does not implement Unwrap or
// the err is nil.
func ExitCode(err error) int {
	exitCode, _ := unwrapExitCode(err)
	return exitCode
}

// Temporary returns the result of invoking the Temporary method for err or
// the first wrapped error recursing until an error does not implement Unwrap
// or the err is nil.
func Temporary(err error) bool {
	isTemporary, _ := unwrapTemporary(err)
	return isTemporary
}

// unwrapExitCode returns the exit code from err or the first unwrapped error
// that implements the ExitCode method. If err and none of its unwrapped
// values implement ExitCode, this returns 0.
func unwrapExitCode(err error) (int, bool) {
	for {
		switch tv := err.(type) {
		case nil:
			// When nil, return the default value.
			return 0, false
		case *Error:
			if tv == nil {
				// When nil, return the default value.
				return 0, false
			}
			return tv.exitCode, tv.isExitCodeSet
		case exitCoder:
			// When err implements ExitCode then return it.
			return tv.ExitCode(), true
		case unwrapper:
			// When error implements Unwrap, then recurse.
			err = tv.Unwrap()
		default:
			// When none of the above, return the default value.
			return 0, false
		}
	}
}

// unwrapTempoary returns whether err is temporary, or the result of invoking
// Temporary method of the first unwrapped error it unwraps.  If err and none
// of its unwrapped values implement Temporary, this returns false.
func unwrapTemporary(err error) (bool, bool) {
	for {
		switch tv := err.(type) {
		case nil:
			// When nil, return the default value.
			return false, false
		case *Error:
			if tv == nil {
				// When nil, return the default value.
				return false, false
			}
			return tv.temporary, tv.isTemporarySet
		case temporaryer:
			// When err implements ExitCode then return it.
			return tv.Temporary(), true
		case unwrapper:
			// When error implements Unwrap, then recurse.
			err = tv.Unwrap()
		default:
			// When none of the above, return the default value.
			return false, false
		}
	}
}
