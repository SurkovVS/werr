package werr

import (
	"errors"
	"fmt"
	"net/http"
)

// For use in errors.As(err,AsMutch)
var AsMutch = func() **werror {
	w := &werror{}
	return &w
}()

type (
	werror struct {
		status int
		err    error
	}
)

// Init new werr
func New(err error) *werror {
	werr, ok := err.(*werror)
	if ok {
		return werr
	}
	return &werror{err: err}
}

// Error implemetation
func (werr *werror) Error() string {
	if text := http.StatusText(werr.status); text != "" {
		return fmt.Sprintf("(werr status - %d %s) ", werr.status, text) + werr.err.Error()
	}
	return werr.err.Error()
}

// Set HTTP status code
func (werr *werror) SetStatus(status int) *werror {
	werr.status = status
	return werr
}

// Get HTTP status code
func (werr *werror) Status() int {
	return werr.status
}

// Wrap error with text as <text>: <current werror>
func (werr *werror) Wrap(t string) *werror {
	werr.err = fmt.Errorf("%s: %w", t, werr.err)
	return werr
}

// Return unwrapped error
func (werr *werror) Unwrap() error {
	return errors.Unwrap(werr.err)
}
