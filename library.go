package tapedeck

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
	"path/filepath"
	"strconv"
)

func Open(c Config) (*Library, error) {
	c.SetDefaults()

	db, err := gorm.Open("sqlite3", c.DBPath)
	l := &Library{db, c}
	if err != nil {
		return l, err
	}

	models := []interface{}{
		&Artist{},
		&Album{},
		&Genre{},
		&Song{},
	}

	if err := db.AutoMigrate(models...).Error; err != nil {
		return l, err
	}

	return l, nil
}

type Library struct {
	db     *gorm.DB
	Config Config
}

func (l *Library) Close() error {
	return l.db.Close()
}

func (l *Library) SaveSong(song *Song) error {
	if song.ArtistID != 0 {
		song.Artist = Artist{}
		if err := l.db.Where(&Artist{ID: song.ArtistID}).First(&song.Artist).Error; err != nil {
			return err
		}
	}

	if song.AlbumID != 0 {
		song.Album = Album{}
		if err := l.db.Where(&Album{ID: song.AlbumID}).First(&song.Album).Error; err != nil {
			return err
		}
	}

	if song.GenreID != 0 {
		song.Genre = Genre{}
		if err := l.db.Where(&Genre{ID: song.GenreID}).First(&song.Genre).Error; err != nil {
			return err
		}
	}

	oldFilepath := song.Filepath
	newFilepath := l.getFilepathForSong(*song)

	if err := syncSongToFile(song); err != nil {
		return err
	}

	if err := MoveFile(oldFilepath, newFilepath, l.Config.FileMode); err != nil {
		return err
	}
	song.Filepath = newFilepath

	if err := l.db.Save(&song).Error; err != nil {
		return err
	}

	return ReducePath(filepath.Dir(oldFilepath), l.Config.BasePath)
}

func (l *Library) RemoveSong(song *Song) error {
	var albumCount, artistCount int

	albumQuery := Song{AlbumID: song.AlbumID}
	artistQuery := Song{ArtistID: song.ArtistID}

	if err := l.db.Model(&albumQuery).Where(albumQuery).Count(&albumCount).Error; err != nil {
		return err
	}

	if err := l.db.Model(&artistQuery).Where(artistQuery).Count(&artistCount).Error; err != nil {
		return err
	}

	oldFilepath := song.Filepath
	tx := l.db.Begin()

	if err := tx.Delete(song).Error; err != nil {
		return err
	}
	song.ID = 0

	if albumCount <= 1 {
		if err := tx.Delete(&Album{ID: song.AlbumID}).Error; err != nil {
			tx.Rollback()
			return err
		}
		song.AlbumID = 0
	}

	if artistCount <= 1 {
		if err := tx.Delete(&Artist{ID: song.ArtistID}).Error; err != nil {
			tx.Rollback()
			return err
		}
		song.ArtistID = 0
	}

	if err := os.Remove(oldFilepath); err != nil {
		return nil
	}
	song.Filepath = ""

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return ReducePath(filepath.Dir(oldFilepath), l.Config.BasePath)
}

func (l *Library) AddSongFromFile(filepath string) (*Song, error) {
	song, err := newSongFromFile(filepath)
	if err != nil {
		return &song, err
	}

	if existingSong := l.getSongByFingerprint(song.Fingerprint); existingSong != nil {
		return existingSong, nil
	}

	newFilepath := l.getFilepathForSong(song)
	if err := CopyFile(filepath, newFilepath, l.Config.FileMode); err != nil {
		return &song, err
	}
	song.Filepath = newFilepath

	return &song, l.db.Save(&song).Error
}

func (l *Library) getSongByFingerprint(fp string) *Song {
	song := Song{Fingerprint: fp}
	if err := l.db.Where(song).First(&song).Error; err != nil {
		return nil
	} else {
		return &song
	}
}

func (l *Library) getFilepathForSong(song Song) string {
	filename := os.Expand(l.Config.FilepathTemplate, func(key string) string {
		switch key {
		case "SongName":
			return song.Name
		case "ArtistName":
			return song.Artist.Name
		case "AlbumName":
			return song.Album.Name
		case "Genre":
			return song.Genre.Name
		case "TrackNo":
			return strconv.Itoa(song.Track)
		case "Year":
			return strconv.Itoa(song.Year)
		case "FileExt":
			return filepath.Ext(song.Filepath)
		default:
			return "$" + key
		}
	})

	return filepath.Join(l.Config.BasePath, filename)
}
