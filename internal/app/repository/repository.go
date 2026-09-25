package repository

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/AlexanderTal/catalog-service/internal/app/entity"
)

type (
	Category interface {
		Create(ctx context.Context, category entity.Category) error
		GetByGUIDs(ctx context.Context, guids []uuid.UUID) ([]entity.Category, error)
		Update(ctx context.Context, category entity.Category) error
		Delete(ctx context.Context, guid uuid.UUID) error
		List(ctx context.Context, name *string) ([]entity.Category, error)
	}

	Product interface {
		Create(ctx context.Context, product entity.Product) error
		GetByGUIDs(ctx context.Context, guids []uuid.UUID) ([]entity.Product, error)
		Update(ctx context.Context, product entity.Product) error
		Delete(ctx context.Context, guid uuid.UUID) error
		List(ctx context.Context, name *string, categoryGUID *uuid.UUID, minPrice, maxPrice *int64) ([]entity.Product, error)
	}
)
