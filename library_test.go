package tapedeck

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func Test_Open__Migrations(t *testing.T) {
	tests := map[interface{}]string{
		&Artist{}: "artists",
		&Album{}:  "albums",
		&Genre{}:  "genres",
		&Song{}:   "songs",
	}

	library := newTempLibrary()
	defer library.Close()

	for model, tableName := range tests {
		if !library.db.HasTable(model) {
			t.Errorf("Expected to have table: %s", tableName)
		}
	}
}

func Test_Library_Close(t *testing.T) {
	library := newTempLibrary()
	if err := library.Close(); err != nil {
		t.Errorf("Did not expect error: %s", err)
	}
}

func Test_Library_AddSongFromFile(t *testing.T) {
	library := newTempLibrary()
	defer library.Close()

	song, err := library.AddSongFromFile("fixtures/sample.mp3")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if song.ID == 0 {
		t.Errorf("Expected song to be saved")
	}

	if song.ArtistID == 0 {
		t.Errorf("Expected artist to be saved")
	}

	if song.AlbumID == 0 {
		t.Errorf("Expected album to be saved")
	}

	if song.GenreID == 0 {
		t.Errorf("Expected genre to be saved")
	}

	filepath := filepath.Join(
		library.Config.BasePath,
		song.Artist.Name,
		song.Album.Name,
		fmt.Sprintf("%d - %s.mp3", song.Track, song.Name),
	)
	if song.Filepath != filepath {
		t.Errorf("Expected filepath to be %s, not %s", filepath, song.Filepath)
	}

	fingerprint := FingerprintFile("fixtures/sample.mp3")
	if song.Fingerprint != fingerprint {
		t.Errorf("Expected fingerprint to be %s, not %s", fingerprint, song.Fingerprint)
	}

	sameSong, err := library.AddSongFromFile("fixtures/sample.mp3")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if song.ID != sameSong.ID {
		t.Errorf("Expected to receive the same song when same file given")
	}
}

func Test_Library_SaveSong(t *testing.T) {
	library := newTempLibrary()
	defer library.Close()

	song, err := library.AddSongFromFile("fixtures/sample.mp3")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	oldDir := filepath.Dir(song.Filepath)

	song.Name = "Passing"

	if err := library.SaveSong(song); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	newFilepath := filepath.Join(
		library.Config.BasePath,
		song.Artist.Name,
		song.Album.Name,
		fmt.Sprintf("%d - %s.mp3", song.Track, song.Name),
	)

	updatedSong, err := newSongFromFile(newFilepath)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if song.Name != updatedSong.Name {
		t.Errorf("Expected song name to be %s, not %s", song.Name, updatedSong.Name)
	}

	if song.Filepath != updatedSong.Filepath {
		t.Errorf("Expected song path to be %s, not %s", updatedSong.Filepath, song.Filepath)
	}

	if _, err := os.Stat(oldDir); os.IsExist(err) {
		t.Errorf("Expected %s to be deleted", oldDir)
	}
}

func Test_Library_RemoveSong(t *testing.T) {
	library := newTempLibrary()
	defer library.Close()

	song, err := library.AddSongFromFile("fixtures/sample.mp3")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	oldFilepath := song.Filepath
	oldDir := filepath.Dir(oldFilepath)

	if err := library.RemoveSong(song); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if song.ID != 0 {
		t.Errorf("Expected song's ID to be zero")
	}

	if song.ArtistID != 0 {
		t.Errorf("Expected song's ArtistID to be zero")
	}

	if song.AlbumID != 0 {
		t.Errorf("Expected song's AlbumID to be zero")
	}

	if song.Filepath != "" {
		t.Errorf("Expected song's filepath to be empty")
	}

	if _, err := os.Stat(oldFilepath); !os.IsNotExist(err) {
		t.Errorf("Expected %s to be deleted", oldFilepath)
	}

	if _, err := os.Stat(oldDir); !os.IsNotExist(err) {
		t.Errorf("Expected %s to be deleted", oldDir)
	}
}

func Test_Library_Rescan(t *testing.T) {
	library := newTempLibrary()
	defer library.Close()

	if err := CopyFile("fixtures/sample.mp3", filepath.Join(library.Config.BasePath, "sample.mp3"), 0755); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if err := library.Rescan(); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	fingerprint := FingerprintFile("fixtures/sample.mp3")
	song := library.getSongByFingerprint(fingerprint)
	if song == nil {
		t.Errorf("Expected to have song")
	}
}

func Test_Library_getFilepathFromSong(t *testing.T) {
	library := newTempLibrary()
	defer library.Close()

	song, err := newSongFromFile("fixtures/sample.mp3")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	tests := map[string]string{
		"$ArtistName": song.Artist.Name,
		"$AlbumName":  song.Album.Name,
		"$SongName":   song.Name,
		"$TrackNo":    fmt.Sprintf("%d", song.Track),
		"$Genre":      song.Genre.Name,
		"$Year":       fmt.Sprintf("%d", song.Year),
		"$FileExt":    ".mp3",
	}

	for template, path := range tests {
		library.Config.FilepathTemplate = template
		expected := filepath.Join(library.Config.BasePath, path)
		if actual := library.getFilepathForSong(song); expected != actual {
			t.Errorf("Expected %s, not %s", expected, actual)
		}
	}
}
