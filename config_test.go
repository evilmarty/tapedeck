package tapedeck

import (
	"testing"
)

func Test_Config_SetDefaults(t *testing.T) {
	c := Config{}
	c.SetDefaults()

	if c.DBPath != DefaultDBPath {
		t.Errorf("Expected %s, received %s", DefaultDBPath, c.DBPath)
	}

	if c.BasePath != DefaultBasePath {
		t.Errorf("Expected %s, received %s", DefaultBasePath, c.BasePath)
	}

	if c.FilepathTemplate != DefaultFilepathTemplate {
		t.Errorf("Expected %s, received %s", DefaultFilepathTemplate, c.FilepathTemplate)
	}

	if c.FileMode != DefaultFileMode {
		t.Errorf("Expected %d, received %d", DefaultFileMode, c.FileMode)
	}
}
