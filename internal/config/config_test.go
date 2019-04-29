package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	t.Run("Validate checks for DBFilePath", func(t *testing.T) {
		c := API{}

		err := c.Validate()
		if err != ErrDBFilepathRequired {
			t.Errorf("Expected err to be %s, got %s", ErrDBFilepathRequired, err)
		}
	})

	t.Run("Validate checks for C2Endpoint", func(t *testing.T) {
		c := API{
			DBFilepath: "/some/file",
		}

		err := c.Validate()
		if err != ErrC2EndpointRequired {
			t.Errorf("Expected err to be %s, got %s", ErrC2EndpointRequired, err)
		}
	})

	t.Run("Validate checks for C2Certificate", func(t *testing.T) {
		c := API{
			DBFilepath: "/some/file",
			C2Endpoint: "someEndpoint",
		}

		err := c.Validate()
		if err != ErrC2CertificateRequired {
			t.Errorf("Expected err to be %s, got %s", ErrC2CertificateRequired, err)
		}
	})

	t.Run("Validate checks that C2Certificate file exists and is readable", func(t *testing.T) {
		expectedPath := "/unknow/file"

		c := API{
			DBFilepath:    "/some/file",
			C2Endpoint:    "someEndpoint",
			C2Certificate: expectedPath,
		}

		err := c.Validate()
		typedErr, ok := err.(ErrC2CertificatePath)
		if !ok {
			t.Errorf("Expected err to be a ErrC2CertificatePath error, got %T", err)
		}

		if typedErr.Path != expectedPath {
			t.Errorf("Expected error path to be %s, got %s", expectedPath, typedErr.Path)
		}

		if !os.IsNotExist(typedErr.Reason) {
			t.Errorf("Expectged error reason to be a os.ErrNotExists, got %T", typedErr.Reason)
		}
	})

	t.Run("Validate on valide configuration returns no errors", func(t *testing.T) {
		tempFile, err := ioutil.TempFile(os.TempDir(), "")
		if err != nil {
			t.Fatalf("Failed to create temporary file")
		}
		defer func() {
			tempFile.Close()
			os.Remove(tempFile.Name())
		}()

		c := API{
			DBFilepath:    "/some/file",
			C2Endpoint:    "something",
			C2Certificate: tempFile.Name(),
		}

		err = c.Validate()
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}
	})
}
