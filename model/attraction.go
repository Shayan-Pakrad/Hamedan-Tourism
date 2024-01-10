package model

import "github.com/uptrace/bun"

type Attraction struct {
	bun.BaseModel `bun:"table:attractions"`

	ID       int64  `bun:"id,pk,autoincrement"`
	ImageURL string `bun:"image_url"`
	Title    string `bun:"title"`
	Brief    string `bun:"brief"`
	Content  string `bun:"content"`
}
