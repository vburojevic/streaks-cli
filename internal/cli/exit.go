package cli

import (
	"errors"
	"fmt"
)

const (
	ExitCodeUsage            = 2
	ExitCodeAppMissing       = 10
	ExitCodeShortcutsMissing = 11
	ExitCodeWrappersMissing  = 12
	ExitCodeActionFailed     = 13
)

type ExitError struct {
	Code int
	Err  error
}

func (e ExitError) Error() string {
	return e.Err.Error()
}

func exitError(code int, err error) error {
	return ExitError{Code: code, Err: err}
}

func exitCodeFromError(err error) (int, error) {
	var ee ExitError
	if errors.As(err, &ee) {
		return ee.Code, ee.Err
	}
	return 0, err
}

func formatExitError(err error) error {
	if err == nil {
		return nil
	}
	code, inner := exitCodeFromError(err)
	if code == 0 {
		return err
	}
	return fmt.Errorf("%w", inner)
}
