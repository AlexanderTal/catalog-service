package repository

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/AlexanderTal/catalog-service/internal/app/entity"
)

type (
	// Создаем интерфейс для управления транзакциями
	Transactional interface {
		InsideTx(ctx context.Context, fn func(ctx context.Context) error) error
	}

	Category interface {
		Transactional // Встраиваем
		Create(ctx context.Context, category entity.Category) error
		GetByGUIDs(ctx context.Context, guids []uuid.UUID) ([]entity.Category, error)
		Update(ctx context.Context, category entity.Category) error
		Delete(ctx context.Context, guid uuid.UUID) error
		List(ctx context.Context, name *string) ([]entity.Category, error)
	}

	Product interface {
		Transactional // Встраиваем
		Create(ctx context.Context, product entity.Product) error
		GetByGUIDs(ctx context.Context, guids []uuid.UUID) ([]entity.Product, error)
		Update(ctx context.Context, product entity.Product) error
		Delete(ctx context.Context, guid uuid.UUID) error
		List(ctx context.Context, name *string, categoryGUID *uuid.UUID, minPrice, maxPrice *int64) ([]entity.Product, error)
	}
)
