package service

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/AlexanderTal/catalog-service/internal/app/entity"
)

type (
	Category interface {
		Create(ctx context.Context, req entity.RequestCategoryCreate) (entity.Category, error)
		GetByGUIDs(ctx context.Context, guids []uuid.UUID) ([]entity.Category, error)
		Update(ctx context.Context, guid uuid.UUID, req entity.RequestCategoryUpdate) (entity.Category, error)
		Delete(ctx context.Context, guid uuid.UUID) error
		List(ctx context.Context) ([]entity.Category, error)
	}

	Product interface {
		Create(ctx context.Context, req entity.RequestProductCreate) (entity.Product, error)
		GetByGUIDs(ctx context.Context, guids []uuid.UUID) ([]entity.Product, error)
		Update(ctx context.Context, guid uuid.UUID, req entity.RequestProductUpdate) (entity.Product, error)
		Delete(ctx context.Context, guid uuid.UUID) error
		List(ctx context.Context, req entity.RequestProductList) ([]entity.Product, error)
	}
)
