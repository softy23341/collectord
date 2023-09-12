package errs

import "git.softndit.com/collector/backend/erry"

var (
	Internal = erry.NewCategory()
)

// Error2HTTPCode maps error to HTTP code
func Error2HTTPCode(err error) int {
	if Internal.Contains(err) {
		return 500
	}
	return 422
}
