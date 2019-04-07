package config

import "errors"

// API describes the configuration required for the API application
type API struct {
	DBFilepath string
}

// Config validation errors
var (
	ErrDBFilepathRequired = errors.New("database file path is required")
)

// Validate checks configuration and return errors when invalid
func (c API) Validate() error {
	if len(c.DBFilepath) == 0 {
		return ErrDBFilepathRequired
	}

	return nil
}
