package model

import "github.com/uptrace/bun"

type Blog struct {
	bun.BaseModel `bun:"table:blogs"`

	ID      int64  `bun:"id,pk,autoincrement"`
	Title   string `bun:"title"`
	Content string `bun:"content"`
	Brief   string `bun:"brief"`
	Views   int    `bun:"views,default:0"`
}
