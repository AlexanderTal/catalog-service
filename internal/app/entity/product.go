package entity

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:product"`

	ID           int64     `bun:"id,autoincrement"`
	GUID         uuid.UUID `bun:"guid,pk"`
	Name         string    `bun:"name"`
	Description  *string   `bun:"description"`
	Price        int64     `bun:"price"`
	CategoryGUID uuid.UUID `bun:"category_guid"`
	CreatedAt    time.Time `bun:"created_at"`
	UpdatedAt    time.Time `bun:"updated_at"`
}

////////////////////////////////////////////////////////////////////////////////
///// HTTP REQUEST & RESPONSE //////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

type RequestProductCreate struct {
	Name         string    `json:"name"`
	Description  *string   `json:"description"`
	Price        int64     `json:"price"`
	CategoryGUID uuid.UUID `json:"category_guid"`
}

func (r RequestProductCreate) Validate() error {
	if r.Name == "" {
		return ErrIncorrectParameters
	}
	if r.Price <= 0 {
		return ErrIncorrectParameters
	}
	if r.CategoryGUID.IsNil() {
		return ErrIncorrectParameters
	}
	return nil
}

type RequestProductUpdate struct {
	Name         string    `json:"name"`
	Description  *string   `json:"description"`
	Price        int64     `json:"price"`
	CategoryGUID uuid.UUID `json:"category_guid"`
}

func (r RequestProductUpdate) Validate() error {
	if r.Price < 0 {
		return ErrIncorrectParameters
	}
	return nil
}

type RequestProductList struct {
	CategoryGUID *uuid.UUID `json:"category_guid"`
}

type ResponseProductCreate struct {
	GUID         uuid.UUID `json:"guid"`
	Name         string    `json:"name"`
	Description  *string   `json:"description"`
	Price        int64     `json:"price"`
	CategoryGUID uuid.UUID `json:"category_guid"`
	CreatedAt    time.Time `json:"created_at"`
}

type ResponseProductUpdate struct {
	GUID         uuid.UUID `json:"guid"`
	Name         string    `json:"name"`
	Description  *string   `json:"description"`
	Price        int64     `json:"price"`
	CategoryGUID uuid.UUID `json:"category_guid"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ResponseProductList struct {
	Data []ResponseProductListItem `json:"data"`
}

type ResponseProductListItem struct {
	GUID         uuid.UUID `json:"guid"`
	Name         string    `json:"name"`
	Description  *string   `json:"description"`
	Price        int64     `json:"price"`
	CategoryGUID uuid.UUID `json:"category_guid"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
