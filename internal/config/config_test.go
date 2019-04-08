package config

import "testing"

func TestConfig(t *testing.T) {
	t.Run("Validate checks for DBFilePath", func(t *testing.T) {
		c := API{}

		err := c.Validate()
		if err != ErrDBFilepathRequired {
			t.Errorf("Expected err to be %s, got %s", ErrDBFilepathRequired, err)
		}
	})

	t.Run("Validate on valide configuration returns no errors", func(t *testing.T) {
		c := API{
			DBFilepath: "/some/file",
		}

		err := c.Validate()
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}
	})
}
