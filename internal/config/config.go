package config

import (
	"errors"
	"fmt"
	"os"
)

// API describes the configuration required for the API application
type API struct {
	DBFilepath    string
	Addr          string
	C2Endpoint    string
	C2Certificate string
}

// Config validation errors
var (
	ErrDBFilepathRequired    = errors.New("database file path is required")
	ErrC2EndpointRequired    = errors.New("c2 endpoint is required")
	ErrC2CertificateRequired = errors.New("c2 certificate is required")
)

// ErrC2CertificatePath is an error returned when failing to read the certificate file at given path
// Like permissions error, or file not found...
type ErrC2CertificatePath struct {
	Path   string
	Reason error
}

func (e ErrC2CertificatePath) Error() string {
	return fmt.Sprintf("certificate file not found at %s: %s", e.Path, e.Reason)
}

// Validate checks configuration and return errors when invalid
func (c API) Validate() error {
	if len(c.DBFilepath) == 0 {
		return ErrDBFilepathRequired
	}

	if len(c.C2Endpoint) == 0 {
		return ErrC2EndpointRequired
	}

	if len(c.C2Certificate) == 0 {
		return ErrC2CertificateRequired
	}

	if _, err := os.Stat(c.C2Certificate); err != nil {
		return ErrC2CertificatePath{
			Path:   c.C2Certificate,
			Reason: err,
		}
	}

	return nil
}
