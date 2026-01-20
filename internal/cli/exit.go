package cli

import "errors"

const (
	ExitCodeUsage            = 2
	ExitCodeAppMissing       = 10
	ExitCodeShortcutsMissing = 11
	ExitCodeShortcutMissing  = 12
	ExitCodeActionFailed     = 13
)

func errorCodeLabel(code int) string {
	switch code {
	case ExitCodeUsage:
		return "usage"
	case ExitCodeAppMissing:
		return "app_missing"
	case ExitCodeShortcutsMissing:
		return "shortcuts_missing"
	case ExitCodeShortcutMissing:
		return "shortcut_missing"
	case ExitCodeActionFailed:
		return "action_failed"
	default:
		return ""
	}
}

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
