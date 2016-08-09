package tapedeck

import (
	"time"
)

type Album struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name string

	Artist   Artist
	ArtistID uint `gorm:"index"`

	Songs []Song
}

func (a Album) String() string {
	return a.Name
}
