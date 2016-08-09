package tapedeck

import (
	"os"
	"testing"
)

func Test_tempLibrary_Close(t *testing.T) {
	library := newTempLibrary()
	if err := library.Close(); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if _, err := os.Stat(library.Config.BasePath); os.IsExist(err) {
		t.Errorf("Expected %s to be deleted", library.Config.BasePath)
	}
}
