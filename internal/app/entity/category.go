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

////////////////////////////////////////////////////////////////////////////////
///// HTTP REQUEST & RESPONSE //////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

type RequestCategoryCreate struct {
	Name string `json:"name"`
}

func (r RequestCategoryCreate) Validate() error {
	if r.Name == "" {
		return ErrIncorrectParameters
	}
	return nil
}

type RequestCategoryUpdate struct {
	Name string `json:"name"`
}

func (r RequestCategoryUpdate) Validate() error {
	if r.Name != "" && len(r.Name) < 2 {
		return ErrIncorrectParameters
	}
	return nil
}

type ResponseCategoryCreate struct {
	GUID      uuid.UUID `json:"guid"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type ResponseCategoryUpdate struct {
	GUID      uuid.UUID `json:"guid"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResponseCategoryList struct {
	Data []ResponseCategoryListItem `json:"data"`
}

type ResponseCategoryListItem struct {
	GUID      uuid.UUID `json:"guid"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
