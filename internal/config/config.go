package config

import (
	"errors"
	"fmt"
	"os"

	slibcfg "gitlab.com/teserakt/serverlib/config"
)

// API describes the configuration required for the API application
type API struct {
	Server              ServerCfg
	DB                  DBCfg
	ES                  ESCfg
	C2Endpoint          string
	C2Certificate       string
	OpencensusSampleAll bool
	OpencensusAddress   string
}

// ServerCfg holds configuration for api server
type ServerCfg struct {
	GRPCAddr     string
	GRPCCert     string
	GRPCKey      string
	HTTPAddr     string
	HTTPGRPCAddr string
	HTTPCert     string
	HTTPKey      string
}

// DBCfg holds configuration for databases
type DBCfg struct {
	Logging          bool
	Type             slibcfg.DBType
	File             string
	Host             string
	Database         string
	Username         string
	Password         string
	Passphrase       string
	Schema           string
	SecureConnection slibcfg.DBSecureConnectionType
}

// ESCfg holds configuration for elasticsearch
type ESCfg struct {
	URLs           []string
	LoggingEnabled bool
	LoggingIndex   string
}

// Config validation errors
var (
	ErrDBFilepathRequired      = errors.New("database file path is required")
	ErrC2EndpointRequired      = errors.New("c2 endpoint is required")
	ErrC2CertificateRequired   = errors.New("c2 certificate is required")
	ErrC2CertificatePath       = errors.New("c2 certificate can't be read")
	ErrGRPCListenAddrRequired  = errors.New("grpc listen address is required")
	ErrHTTPListenAddrRequired  = errors.New("http listen address is required")
	ErrNoPassphrase            = errors.New("no database passphrase supplied")
	ErrNoDBAddr                = errors.New("no database address supplied")
	ErrNoDatabase              = errors.New("no database name supplied")
	ErrUnsupportedDBType       = errors.New("unknown or unsupported database type")
	ErrNoDBFile                = errors.New("no database file supplied")
	ErrNoUsername              = errors.New("no username supplied")
	ErrNoPassword              = errors.New("no password supplied")
	ErrInvalidSecureConnection = errors.New("invalid secure connection mode")
	ErrNoSchema                = errors.New("no schema supplied")
	ErrGRPCCertRequired        = errors.New("grpc certificate path is required")
	ErrGRPCKeyRequired         = errors.New("grpc key path is required")
	ErrHTTPCertRequired        = errors.New("http certificate path is required")
	ErrHTTPKeyRequired         = errors.New("http key path is required")
	ErrHTTPGRPCAddrRequired    = errors.New("http-grpc address is required")
	ErrESLoggingIndexRequired  = errors.New("es-logging-index is required")
	ErrESUrlRequired           = errors.New("at least one es-urls is required")
)

// NewAPI creates a new configuration struct for the C2AE api
func NewAPI() *API {
	return &API{}
}

// ViperCfgFields returns the list of configuration bound's fields to be loaded by viper
func (c *API) ViperCfgFields() []slibcfg.ViperCfgField {
	return []slibcfg.ViperCfgField{
		{&c.Server.GRPCAddr, "listen-grpc", slibcfg.ViperString, "localhost:5556", "C2AE_GRPC_LISTEN_ADDR"},
		{&c.Server.GRPCCert, "grpc-cert", slibcfg.ViperRelativePath, "", "C2AE_GRPC_CERT"},
		{&c.Server.GRPCKey, "grpc-key", slibcfg.ViperRelativePath, "", "C2AE_GRPC_KEY"},
		{&c.Server.HTTPAddr, "listen-http", slibcfg.ViperString, "localhost:8886", "C2AE_HTTP_LISTEN_ADDR"},
		{&c.Server.HTTPGRPCAddr, "http-grpc-addr", slibcfg.ViperString, "localhost:5556", "C2AE_HTTP_GRPC_ADDR"},
		{&c.Server.HTTPCert, "http-cert", slibcfg.ViperRelativePath, "", "C2AE_HTTP_CERT"},
		{&c.Server.HTTPKey, "http-key", slibcfg.ViperRelativePath, "", "C2AE_HTTP_KEY"},

		{&c.DB.Logging, "db-logging", slibcfg.ViperBool, false, ""},
		{&c.DB.Type, "db-type", slibcfg.ViperDBType, "sqlite3", "C2AE_DB_TYPE"},
		{&c.DB.File, "db-file", slibcfg.ViperString, "", "C2AE_DB_PATH"},
		{&c.DB.Host, "db-host", slibcfg.ViperString, "", ""},
		{&c.DB.Database, "db-database", slibcfg.ViperString, "", ""},
		{&c.DB.Schema, "db-schema", slibcfg.ViperString, "", ""},
		{&c.DB.Username, "db-username", slibcfg.ViperString, "", "C2AE_DB_USERNAME"},
		{&c.DB.Password, "db-password", slibcfg.ViperString, "", "C2AE_DB_PASSWORD"},
		{&c.DB.Passphrase, "db-encryption-passphrase", slibcfg.ViperString, "", "C2AE_DB_ENCRYPTION_PASSPHRASE"},
		{&c.DB.SecureConnection, "db-secure-connection", slibcfg.ViperDBSecureConnection, slibcfg.DBSecureConnectionEnabled, "E4C2AE_DB_SECURE_CONNECTION"},

		{&c.C2Endpoint, "c2-host-port", slibcfg.ViperString, "localhost:5555", "C2AE_C2_ENDPOINT"},
		{&c.C2Certificate, "c2-cert", slibcfg.ViperRelativePath, "", "C2AE_C2CERT_PATH"},

		{&c.OpencensusSampleAll, "oc-sample-all", slibcfg.ViperBool, true, ""},
		{&c.OpencensusAddress, "oc-agent-addr", slibcfg.ViperString, "localhost:55678", "C2AE_OC_ENDPOINT"},

		{&c.ES.URLs, "es-urls", slibcfg.ViperStringSlice, nil, ""},
		{&c.ES.LoggingEnabled, "es-logging-enable", slibcfg.ViperBool, false, ""},
		{&c.ES.LoggingIndex, "es-logging-index", slibcfg.ViperString, "logs", ""},
	}
}

// Validate checks configuration and return errors when invalid
func (c API) Validate() error {
	if err := c.Server.Validate(); err != nil {
		return err
	}

	if err := c.DB.Validate(); err != nil {
		return err
	}

	if len(c.C2Endpoint) == 0 {
		return ErrC2EndpointRequired
	}

	if len(c.C2Certificate) == 0 {
		return ErrC2CertificateRequired
	}

	if _, err := os.Stat(c.C2Certificate); err != nil {
		return ErrC2CertificatePath
	}

	if err := c.ES.Validate(); err != nil {
		return err
	}

	return nil
}

// Validate checks ServerCfg and returns an error if anything is invalid
func (c ServerCfg) Validate() error {
	if len(c.GRPCAddr) == 0 {
		return ErrGRPCListenAddrRequired
	}

	if len(c.GRPCCert) == 0 {
		return ErrGRPCCertRequired
	}

	if len(c.GRPCKey) == 0 {
		return ErrGRPCKeyRequired
	}

	if len(c.HTTPAddr) == 0 {
		return ErrHTTPListenAddrRequired
	}

	if len(c.HTTPGRPCAddr) == 0 {
		return ErrHTTPGRPCAddrRequired
	}

	if len(c.HTTPCert) == 0 {
		return ErrHTTPCertRequired
	}

	if len(c.HTTPKey) == 0 {
		return ErrHTTPKeyRequired
	}

	return nil
}

// Validate checks DBCfg and returns an error if anything is invalid
func (c DBCfg) Validate() error {
	if len(c.Passphrase) == 0 {
		return ErrNoPassphrase
	}

	switch c.Type {
	case slibcfg.DBTypePostgres:
		return c.validatePostgres()
	case slibcfg.DBTypeSQLite:
		return c.validateSQLite()
	default:
		return ErrUnsupportedDBType
	}
}

func (c DBCfg) validatePostgres() error {
	if len(c.Host) == 0 {
		return ErrNoDBAddr
	}

	if len(c.Database) == 0 {
		return ErrNoDatabase
	}

	if len(c.Username) == 0 {
		return ErrNoUsername
	}

	if len(c.Password) == 0 {
		return ErrNoPassword
	}

	if len(c.Schema) == 0 {
		return ErrNoSchema
	}

	if c.SecureConnection != slibcfg.DBSecureConnectionEnabled &&
		c.SecureConnection != slibcfg.DBSecureConnectionSelfSigned &&
		c.SecureConnection != slibcfg.DBSecureConnectionInsecure {
		return ErrInvalidSecureConnection
	}

	return nil
}

func (c DBCfg) validateSQLite() error {
	if len(c.File) == 0 {
		return ErrNoDBFile
	}

	return nil
}

// ConnectionString returns the string to use to establish the db connection
func (c DBCfg) ConnectionString() (string, error) {
	switch slibcfg.DBType(c.Type) {
	case slibcfg.DBTypePostgres:
		return fmt.Sprintf(
			"host=%s dbname=%s user=%s password=%s %s",
			c.Host, c.Database, c.Username, c.Password, c.SecureConnection.PostgresSSLMode(),
		), nil
	case slibcfg.DBTypeSQLite:
		return c.File, nil
	default:
		return "", ErrUnsupportedDBType
	}
}

// Log returns parts of the configuration that is safe to log (ie: no passwords)
func (c DBCfg) Log() []interface{} {
	switch slibcfg.DBType(c.Type) {
	case slibcfg.DBTypePostgres:
		return []interface{}{"type", c.Type.String(), "host", c.Host, "dbname", c.Database, "user", c.Username, "secureMode", c.SecureConnection.PostgresSSLMode()}
	case slibcfg.DBTypeSQLite:
		return []interface{}{"type", c.Type.String(), "file", c.File}
	default:
		return []interface{}{"type", "unknown"}
	}
}

// IsEnabled indicate whenever the elasticsearch is required by configuration or not
func (c ESCfg) IsEnabled() bool {
	return c.LoggingEnabled
}

// Validate checks ESCfg and returns an error if anything is invalid
func (c ESCfg) Validate() error {
	if c.LoggingEnabled {
		if len(c.LoggingIndex) == 0 {
			return ErrESLoggingIndexRequired
		}
	}

	if c.IsEnabled() && len(c.URLs) == 0 {
		return ErrESUrlRequired
	}

	return nil
}
