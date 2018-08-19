package churn

import "github.com/pkg/errors"

// Sentinal errors
var (
	ErrNameTaken    = errors.New("name is already in use")
	ErrPortNotExist = errors.New("port does not exist")
)

// IsNameTaken returns true if the given error derives from
// a name being already used
func IsNameTaken(err error) bool {
	return errors.Cause(err) == ErrNameTaken
}

// IsPortNotExist returns true if the given error derives from
// a port not existing
func IsPortNotExist(err error) bool {
	return errors.Cause(err) == ErrPortNotExist
}
