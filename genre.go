package tapedeck

import (
	"time"
)

type Genre struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name string

	Songs []Song
}

func (g Genre) String() string {
	return g.Name
}
