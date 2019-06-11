package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	slibcfg "gitlab.com/teserakt/serverlib/config"
)

func TestConfig(t *testing.T) {
	t.Run("Validate properly checks all configuration fields", func(t *testing.T) {

		validFile, err := ioutil.TempFile("", "")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		validFile.Close()
		defer os.Remove(validFile.Name())

		testCases := []struct {
			cfg         API
			expectedErr error
		}{
			{
				cfg:         API{},
				expectedErr: ErrListenAddrRequired,
			},
			{
				cfg: API{
					Addr: "127.0.0.1:5556",
				},
				expectedErr: ErrNoPassphrase,
			},
			{
				cfg: API{
					Addr: "127.0.0.1:5556",
					DB: DBCfg{
						Passphrase: "something",
					},
				},
				expectedErr: ErrUnsupportedDBType,
			},
			{
				cfg: API{
					Addr: "127.0.0.1:5556",
					DB: DBCfg{
						Passphrase: "something",
						Type:       slibcfg.DBTypeSQLite,
					},
				},
				expectedErr: ErrNoDBFile,
			},
			{
				cfg: API{
					Addr: "127.0.0.1:5556",
					DB: DBCfg{
						Passphrase: "something",
						Type:       slibcfg.DBTypeSQLite,
						File:       "/some/file",
					},
				},
				expectedErr: ErrC2EndpointRequired,
			},
			{
				cfg: API{
					Addr: "127.0.0.1:5556",
					DB: DBCfg{
						Passphrase: "something",
						Type:       slibcfg.DBTypeSQLite,
						File:       "/some/file",
					},
					C2Endpoint: "localhost:5555",
				},
				expectedErr: ErrC2CertificateRequired,
			},
			{
				cfg: API{
					Addr: "127.0.0.1:5556",
					DB: DBCfg{
						Passphrase: "something",
						Type:       slibcfg.DBTypeSQLite,
						File:       "/some/file",
					},
					C2Endpoint:    "localhost:5555",
					C2Certificate: "/some/path",
				},
				expectedErr: ErrC2CertificatePath,
			},
			{
				cfg: API{
					Addr: "127.0.0.1:5556",
					DB: DBCfg{
						Passphrase: "something",
						Type:       slibcfg.DBTypeSQLite,
						File:       "/some/file",
					},
					C2Endpoint:    "localhost:5555",
					C2Certificate: validFile.Name(),
				},
				expectedErr: nil,
			},
			{
				cfg: API{
					Addr: "127.0.0.1:5556",
					DB: DBCfg{
						Type: slibcfg.DBTypePostgres,
					},
					C2Endpoint:    "localhost:5555",
					C2Certificate: validFile.Name(),
				},
				expectedErr: ErrNoPassphrase,
			},
			{
				cfg: API{
					Addr: "127.0.0.1:5556",
					DB: DBCfg{
						Type:       slibcfg.DBTypePostgres,
						Passphrase: "something",
					},
					C2Endpoint:    "localhost:5555",
					C2Certificate: validFile.Name(),
				},
				expectedErr: ErrNoDBAddr,
			},
			{
				cfg: API{
					Addr: "127.0.0.1:5556",
					DB: DBCfg{
						Type:       slibcfg.DBTypePostgres,
						Passphrase: "something",
						Host:       "127.0.0.1:5432",
					},
					C2Endpoint:    "localhost:5555",
					C2Certificate: validFile.Name(),
				},
				expectedErr: ErrNoDatabase,
			},
			{
				cfg: API{
					Addr: "127.0.0.1:5556",
					DB: DBCfg{
						Type:       slibcfg.DBTypePostgres,
						Passphrase: "something",
						Host:       "127.0.0.1:5432",
						Database:   "something",
					},
					C2Endpoint:    "localhost:5555",
					C2Certificate: validFile.Name(),
				},
				expectedErr: ErrNoUsername,
			},
			{
				cfg: API{
					Addr: "127.0.0.1:5556",
					DB: DBCfg{
						Type:       slibcfg.DBTypePostgres,
						Passphrase: "something",
						Host:       "127.0.0.1:5432",
						Database:   "something",
						Username:   "something",
					},
					C2Endpoint:    "localhost:5555",
					C2Certificate: validFile.Name(),
				},
				expectedErr: ErrNoPassword,
			},
			{
				cfg: API{
					Addr: "127.0.0.1:5556",
					DB: DBCfg{
						Type:       slibcfg.DBTypePostgres,
						Passphrase: "something",
						Host:       "127.0.0.1:5432",
						Database:   "something",
						Username:   "something",
						Password:   "something",
					},
					C2Endpoint:    "localhost:5555",
					C2Certificate: validFile.Name(),
				},
				expectedErr: ErrNoSchema,
			},
			{
				cfg: API{
					Addr: "127.0.0.1:5556",
					DB: DBCfg{
						Type:       slibcfg.DBTypePostgres,
						Passphrase: "something",
						Host:       "127.0.0.1:5432",
						Database:   "something",
						Username:   "something",
						Password:   "something",
						Schema:     "schema",
					},
					C2Endpoint:    "localhost:5555",
					C2Certificate: validFile.Name(),
				},
				expectedErr: ErrInvalidSecureConnection,
			},
			{
				cfg: API{
					Addr: "127.0.0.1:5556",
					DB: DBCfg{
						Type:             slibcfg.DBTypePostgres,
						Passphrase:       "something",
						Host:             "127.0.0.1:5432",
						Database:         "something",
						Username:         "something",
						Password:         "something",
						Schema:           "schema",
						SecureConnection: slibcfg.DBSecureConnectionEnabled,
					},
					C2Endpoint:    "localhost:5555",
					C2Certificate: validFile.Name(),
				},
				expectedErr: nil,
			},
			{
				cfg: API{
					Addr: "127.0.0.1:5556",
					DB: DBCfg{
						Type:             slibcfg.DBTypePostgres,
						Passphrase:       "something",
						Host:             "127.0.0.1:5432",
						Database:         "something",
						Username:         "something",
						Password:         "something",
						Schema:           "schema",
						SecureConnection: slibcfg.DBSecureConnectionInsecure,
					},
					C2Endpoint:    "localhost:5555",
					C2Certificate: validFile.Name(),
				},
				expectedErr: nil,
			},
			{
				cfg: API{
					Addr: "127.0.0.1:5556",
					DB: DBCfg{
						Type:             slibcfg.DBTypePostgres,
						Passphrase:       "something",
						Host:             "127.0.0.1:5432",
						Database:         "something",
						Username:         "something",
						Password:         "something",
						Schema:           "schema",
						SecureConnection: slibcfg.DBSecureConnectionSelfSigned,
					},
					C2Endpoint:    "localhost:5555",
					C2Certificate: validFile.Name(),
				},
				expectedErr: nil,
			},
		}

		for _, testCase := range testCases {
			err := testCase.cfg.Validate()
			if err != testCase.expectedErr {
				t.Errorf("Expected error to be %v, got %v", testCase.expectedErr, err)
			}
		}

	})
}

func TestDBCfg(t *testing.T) {
	t.Run("ConnectionString returns the proper connection string for Postgres type", func(t *testing.T) {
		expectedDatabase := "test"
		expectedHost := "some/host:port"
		expectedUsername := "username"
		expectedPassword := "password"

		cfg := DBCfg{
			Type:     slibcfg.DBTypePostgres,
			Database: expectedDatabase,
			Host:     expectedHost,
			Username: expectedUsername,
			Password: expectedPassword,
		}

		expectedConnectionString := fmt.Sprintf(
			"host=%s dbname=%s user=%s password=%s %s",
			expectedHost,
			expectedDatabase,
			expectedUsername,
			expectedPassword,
			slibcfg.PostgresSSLModeFull,
		)

		cnxStr, err := cfg.ConnectionString()

		if err != nil {
			t.Errorf("expected no error, got %s", err)
		}

		if expectedConnectionString != cnxStr {
			t.Errorf("expected connectionString to be %s, got %s", expectedConnectionString, cnxStr)
		}
	})

	t.Run("ConnectionString returns the proper connection string for SQLite type", func(t *testing.T) {
		expectedFile := "some/db/file"

		cfg := DBCfg{
			Type: slibcfg.DBTypeSQLite,
			File: expectedFile,
		}

		cnxStr, err := cfg.ConnectionString()

		if err != nil {
			t.Errorf("expected no error, got %s", err)
		}

		if expectedFile != cnxStr {
			t.Errorf("expected connectionString to be %s, got %s", expectedFile, cnxStr)
		}
	})

	t.Run("ConnectionString returns an error on unsupported DB type", func(t *testing.T) {
		cfg := DBCfg{
			Type: slibcfg.DBType("unknow"),
		}

		_, err := cfg.ConnectionString()

		if err != ErrUnsupportedDBType {
			t.Errorf("Expected err to be %s, got %s", ErrUnsupportedDBType, err)
		}
	})

}
