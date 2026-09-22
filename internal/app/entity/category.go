package entity

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/uptrace/bun"
)

type Category struct {
	bun.BaseModel `bun:"table:category"`

	ID        int64     `bun:"id,autoincrement"`
	GUID      uuid.UUID `bun:"guid,pk"`
	Name      string    `bun:"name"`
	CreatedAt time.Time `bun:"created_at"`
	UpdatedAt time.Time `bun:"updated_at"`
}
