package goerr

import (
	"fmt"
	"sort"
	"strings"
)

// Error holds contextual information about the error, including an optional
// exit code, whether the error is considered temporary, a wrapped error, and
// options and option comments for displaying an error string.
type Error struct {
	optionComments           []optionComment
	options                  []string
	beforeMessage            []string
	betweenMessageAndOptions []string
	afterOptions             []string
	err                      error
	msg                      string
	exitCode                 int
	isExitCodeSet            bool
	temporary                bool
	isTemporarySet           bool
}

type optionComment struct {
	comment string
	index   int
}

type optionCommentSlice []optionComment

func (s optionCommentSlice) Len() int           { return len(s) }
func (x optionCommentSlice) Less(i, j int) bool { return x[i].index > x[j].index }
func (x optionCommentSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// MaybeWrap returns a new Error that wraps err, or returns nil when err is
// nil.
func MaybeWrap(err error) *Error {
	if err == nil {
		return nil
	}
	return &Error{err: err}
}

// New returns a new Error with a formatted message.
func New(f string, a ...any) *Error {
	return &Error{msg: fmt.Sprintf(f, a...)}
}

// Error returns an error message suitable for display.
func (e Error) Error() string {
	return strings.Join(e.ErrorLines(), "\n")
}

// ErrorLines returns error message lines suitable for display.
func (e Error) ErrorLines() []string {
	lines := append([]string(nil), e.beforeMessage...)

	if e.msg != "" {
		if e.err != nil {
			lines = append(lines, e.msg+": "+e.err.Error())
		} else {
			lines = append(lines, e.msg)
		}
	} else {
		if e.err != nil {
			lines = append(lines, e.err.Error())
		} else {
			lines = append(lines, "error without message or wrapped error") // upstream bug
		}
	}

	lines = append(lines, e.betweenMessageAndOptions...)

	// Append option comment lines.
	lines = append(lines, optionLines(e.options, e.optionComments...)...)

	// Append additional lines.
	lines = append(lines, e.afterOptions...)

	return lines
}

// ExitCode returns the exit code stored in this instance, or, if nothing
// stored in this instance, the result of invoking ExitCode on the possibly
// wrapped error, recursing until either a wrapped error implements ExitCode
// method, does not implement Unwrap, or nil error.
func (e Error) ExitCode() int {
	if e.isExitCodeSet {
		return e.exitCode
	}
	return ExitCode(e.err)
}

// Temporary returns the exit code stored in this instance, or, if nothing
// stored in this instance, the result of invoking Temporary on the possibly
// wrapped error, recursing until either a wrapped error implements Temporary
// method, does not implement Unwrap, or nil error.
func (e Error) Temporary() bool {
	if e.isTemporarySet {
		return e.temporary
	}
	return Temporary(e.err)
}

// Unwrap returns the encapsulated error, or nil.
func (e Error) Unwrap() error {
	return e.err
}

// WithExitCode stores code as the value to be returned by the ExitCode
// method.
func (e *Error) WithExitCode(code int) *Error {
	if e == nil {
		return nil
	}
	e.isExitCodeSet = true
	e.exitCode = code
	return e
}

// WithLineAfterOptions appends line to the list of lines to include after any
// option lines in the error message.
func (e *Error) WithLineAfterOptions(line string) *Error {
	if e == nil {
		return nil
	}
	e.afterOptions = append(e.afterOptions, line)
	return e
}

// WithLinesAfterOptions appends lines to the list of lines to include after
// any option lines in the error message.
func (e *Error) WithLinesAfterOptions(lines []string) *Error {
	if e == nil {
		return nil
	}
	e.afterOptions = append(e.afterOptions, lines...)
	return e
}

// WithLineBeforeMessage appends line to the list of lines to include before any
// option lines in the error message.
func (e *Error) WithLineBeforeMessage(line string) *Error {
	if e == nil {
		return nil
	}
	e.beforeMessage = append(e.beforeMessage, line)
	return e
}

// WithLinesBeforeMessage appends lines to the list of lines to include before
// any option lines in the error message.
func (e *Error) WithLinesBeforeMessage(lines []string) *Error {
	if e == nil {
		return nil
	}
	e.beforeMessage = append(e.beforeMessage, lines...)
	return e
}

// WithLineBetweenMessageAndOption appends line to the list of lines to
// include between message and any option lines.
func (e *Error) WithLineBetweenMessageAndOption(line string) *Error {
	if e == nil {
		return nil
	}
	e.betweenMessageAndOptions = append(e.betweenMessageAndOptions, line)
	return e
}

// WithLinesBetweenMessageAndOption appends lines to the list of lines to
// include between message and any option lines.
func (e *Error) WithLinesBetweenMessageAndOption(lines []string) *Error {
	if e == nil {
		return nil
	}
	e.betweenMessageAndOptions = append(e.betweenMessageAndOptions, lines...)
	return e
}

// WithMessage stores a formatted message for the error.
func (e *Error) WithMessage(f string, a ...any) *Error {
	if e == nil {
		return nil
	}
	e.msg = fmt.Sprintf(f, a...)
	return e
}

// WithOptionComment causes an additional error message line to be printed
// that underlines the option indexed by index, with comment.
func (e *Error) WithOptionComment(index int, comment string) *Error {
	if e == nil {
		return nil
	}
	e.optionComments = append(e.optionComments, optionComment{
		comment: comment,
		index:   index,
	})
	return e
}

// WithOptions stores the options to be printed when printing the error
// message.
func (e *Error) WithOptions(options []string) *Error {
	if e == nil {
		return nil
	}
	e.options = options
	return e
}

// WithTemporary stores temporary as the value to be returned by the Temporary
// method.
func (e *Error) WithTemporary(temporary bool) *Error {
	if e == nil {
		return nil
	}
	e.isTemporarySet = true
	e.temporary = temporary
	return e
}

// WithWrap stores err as the value to be returned by the Unwrap method.
func (e *Error) WithWrap(err error) *Error {
	if e == nil {
		return nil
	}
	e.err = err
	return e
}

func optionLines(opts []string, ocs ...optionComment) []string {
	// zero one --two three
	//                ^~~~~ cannot find this file
	//          ^~~~~ for this option
	//      ^~~ for this sub-command

	optCount := len(opts)
	if optCount == 0 {
		return nil
	}

	lines := make([]string, 0, 1+len(ocs))
	lines = append(lines, strings.Join(opts, " "))

	if len(ocs) == 0 {
		return lines
	}

	indices := []int{0} // index of first opt is 0

	var length int

	for _, opt := range opts {
		length += len(opt) + 1
		indices = append(indices, length)
	}

	sort.Sort(optionCommentSlice(ocs))

	for _, oc := range ocs {
		if oc.index < 0 || oc.index >= optCount {
			prefix := strings.Repeat(" ", length)
			lines = append(lines, prefix+"^ "+oc.comment)
			continue
		}

		prefix := strings.Repeat(" ", indices[oc.index]) + "^"

		if oc.index == optCount-1 {
			prefix += strings.Repeat("~", (length-indices[oc.index])-2)
		} else {
			prefix += strings.Repeat("~", (indices[oc.index+1] - indices[oc.index] - 2))
		}

		lines = append(lines, prefix+" "+oc.comment)
	}

	return lines
}
