package pcategory

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/uptrace/bun"

	"github.com/AlexanderTal/catalog-service/internal/app/entity"
	"github.com/AlexanderTal/catalog-service/internal/app/repository"
	rcpostgres "github.com/AlexanderTal/catalog-service/internal/app/repository/conn/postgres"
)

type (
	repoPg struct {
		*_DB
	}

	_DB = rcpostgres.Client
)

func NewRepoFromPostgres(client *rcpostgres.Client) repository.Category {
	return &repoPg{_DB: client}
}

func (r *repoPg) Create(ctx context.Context, category entity.Category) error {
	_, err := r.NewInsert().Model(&category).Exec(ctx)
	return err
}

func (r *repoPg) GetByGUIDs(ctx context.Context, guids []uuid.UUID) ([]entity.Category, error) {
	var categories []entity.Category
	err := r.NewSelect().
		Model(&categories).
		Where("guid IN (?)", bun.List(guids)).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *repoPg) Update(ctx context.Context, category entity.Category) error {
	res, err := r.NewUpdate().
		Model(&category).
		ExcludeColumn("id", "created_at").
		WherePK().
		Exec(ctx)
	return rcpostgres.UpdateErr(res, err)
}

func (r *repoPg) Delete(ctx context.Context, guid uuid.UUID) error {
	_, err := r.NewDelete().
		Model((*entity.Category)(nil)).
		Where("guid = ?", guid).
		Exec(ctx)
	return rcpostgres.DeleteErr(err)
}

func (r *repoPg) List(ctx context.Context, name *string) ([]entity.Category, error) {
	var categories []entity.Category
	q := r.NewSelect().Model(&categories)

	if name != nil {
		q = q.Where("name = ?", *name)
	}

	if err := q.Scan(ctx); err != nil {
		return nil, err
	}
	return categories, nil
}
