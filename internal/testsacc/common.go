package testsacc

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testCheckFileExists(filename string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		filename = filepath.Clean(filename)
		_, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		return nil
	}
}

func testCheckFileNotExists(filename string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		filename = filepath.Clean(filename)
		// Check if file exists
		if _, err := os.Stat(filename); err == nil {
			return fmt.Errorf("file %s exists", filename)
		}

		return nil
	}
}

// This is a helper function that attempts to remove file creating during unit test.
func deleteFile(filename string, t *testing.T) func() {
	return func() {
		if _, err := os.Stat(filename); err == nil {
			err := os.Remove(filename)
			if err != nil {
				t.Errorf("Failed to delete file: %s", err)
			}
		}
	}
}
