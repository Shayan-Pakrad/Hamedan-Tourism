package model

import "github.com/uptrace/bun"

type Attraction struct {
	bun.BaseModel `bun:"table:attractions"`

	ID      int64  `bun:",pk,autoincrement"`
	Title   string `bun:"title"`
	Content string `bun:"content"`
}
