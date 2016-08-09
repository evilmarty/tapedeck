package tapedeck

import (
	"time"
)

type Artist struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name string

	Songs []Song
}

func (a Artist) String() string {
	return a.Name
}
