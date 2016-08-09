package tapedeck

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func Test_newSongFromFile(t *testing.T) {
	name := "Test"
	length := time.Duration(1000000000)
	track := 1
	year := 2016
	bitrate := 48
	channels := 1
	samplerate := 44100
	comment := "Test"
	artist := "The Perfectionist"
	album := "Testing"
	genre := "Test"

	song, err := newSongFromFile("fixtures/sample.mp3")
	if err != nil {
		t.Fatalf("Error not expected: %s", err)
	}

	if song.Name != name {
		t.Errorf("Expected name to be %s, not %s", name, song.Name)
	}

	if song.Track != track {
		t.Errorf("Expected track to be %d, not %d", track, song.Track)
	}

	if song.Year != year {
		t.Errorf("Expected year to be %d, not %d", year, song.Year)
	}

	if song.Length != length {
		t.Errorf("Expected length to be %d, not %d", length, song.Length)
	}

	if song.Bitrate != bitrate {
		t.Errorf("Expected bitrate to be %d, not %d", bitrate, song.Bitrate)
	}

	if song.Channels != channels {
		t.Errorf("Expected channels to be %d, not %d", channels, song.Channels)
	}

	if song.Samplerate != samplerate {
		t.Errorf("Expected samplerate to be %d, not %d", samplerate, song.Samplerate)
	}

	if song.Comment != comment {
		t.Errorf("Expected comment to be %s, not %s", comment, song.Comment)
	}

	if song.Artist.Name != artist {
		t.Errorf("Expected artist to be %s, not %s", artist, song.Artist.Name)
	}

	if song.Album.Name != album {
		t.Errorf("Expected album to be %s, not %s", album, song.Album.Name)
	}

	if song.Genre.Name != genre {
		t.Errorf("Expected genre to be %s, not %s", genre, song.Genre.Name)
	}
}

func Test_syncSongToFile(t *testing.T) {
	file, err := ioutil.TempFile("", "song-")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	filepath := file.Name() + ".mp3"
	file.Close()
	os.Remove(file.Name())

	if err := CopyFile("fixtures/sample.mp3", filepath, 0500); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	song, err := newSongFromFile(filepath)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	song.Name = "Passing"

	if err := syncSongToFile(&song); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	updatedSong, err := newSongFromFile(filepath)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if song.Name != updatedSong.Name {
		t.Errorf("Expected %s, not %s", song.Name, updatedSong.Name)
	}
}
