package tapedeck

import (
	"github.com/wtolson/go-taglib"
	"time"
)

type Song struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Filepath    string `gorm:"not null;unique"`
	Fingerprint string `gorm:"not null;unique_index"`
	Name        string
	Year        int
	Track       int
	Length      time.Duration
	Bitrate     int
	Samplerate  int
	Channels    int
	Comment     string

	Artist   Artist
	ArtistID uint `gorm:"index"`

	Album   Album
	AlbumID uint `gorm:"index"`

	Genre   Genre
	GenreID uint `gorm:"index"`
}

func (s Song) String() string {
	return s.Name
}

func newSongFromFile(filepath string) (Song, error) {
	s := Song{}
	file, err := taglib.Read(filepath)
	if err != nil {
		return s, newTagError(filepath, err)
	}
	defer file.Close()

	s.Filepath = filepath
	s.Fingerprint = FingerprintFile(filepath)
	s.Name = file.Title()
	s.Track = file.Track()
	s.Length = file.Length()
	s.Year = file.Year()
	s.Bitrate = file.Bitrate()
	s.Channels = file.Channels()
	s.Samplerate = file.Samplerate()
	s.Comment = file.Comment()
	s.Artist.Name = file.Artist()
	s.Album.Name = file.Album()
	s.Genre.Name = file.Genre()

	if s.Name == "" {
		s.Name = GetFilename(s.Filepath)
	}

	return s, nil
}

func syncSongToFile(s *Song) error {
	file, err := taglib.Read(s.Filepath)
	if err != nil {
		return newTagError(s.Filepath, err)
	}
	defer file.Close()

	file.SetTitle(s.Name)
	file.SetTrack(s.Track)
	file.SetYear(s.Year)
	file.SetComment(s.Comment)
	file.SetArtist(s.Artist.Name)
	file.SetAlbum(s.Album.Name)
	file.SetGenre(s.Genre.Name)

	if err := file.Save(); err != nil {
		return newTagError(s.Filepath, err)
	}

	s.Fingerprint = FingerprintFile(s.Filepath)

	return nil
}
