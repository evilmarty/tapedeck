package tapedeck

import (
	"path/filepath"
)

const (
	DefaultBasePath = "media"
	DefaultDBPath   = "library.db"
	DefaultFileMode = 0755
)

var (
	DefaultFilepathTemplate = filepath.Join("$ArtistName", "$AlbumName", "$TrackNo - $SongName$FileExt")
)

type Config struct {
	DBPath           string
	BasePath         string
	FilepathTemplate string
	FileMode         uint
}

func (c *Config) SetDefaults() {
	if c.DBPath == "" {
		c.DBPath = DefaultDBPath
	}

	if c.BasePath == "" {
		c.BasePath = DefaultBasePath
	}

	if c.FilepathTemplate == "" {
		c.FilepathTemplate = DefaultFilepathTemplate
	}

	if c.FileMode == 0 {
		c.FileMode = DefaultFileMode
	}
}
