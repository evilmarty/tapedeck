package tapedeck

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type tempLibrary struct {
	Library
}

func (tl *tempLibrary) Close() error {
	if err := tl.Library.Close(); err != nil {
		return err
	}
	return os.RemoveAll(tl.Config.BasePath)
}

func newTempLibrary() *tempLibrary {
	dir, err := ioutil.TempDir("", "library-")
	if err != nil {
		panic(err)
	}

	config := Config{
		DBPath:   filepath.Join(dir, "library.db"),
		BasePath: dir,
	}

	library, err := Open(config)
	if err != nil {
		panic(err)
	}

	return &tempLibrary{*library}
}
