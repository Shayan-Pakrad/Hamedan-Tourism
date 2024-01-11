package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Event struct {
	bun.BaseModel `bun:"table:events"`

	ID    int64     `bun:"id,pk,autoincrement"`
	Title string    `bun:"title"`
	Brief string    `bun:"brief"`
	Due   time.Time `bun:"due"`
}
