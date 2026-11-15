package sproduct

import (
	"context"
	"errors"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/AlexanderTal/catalog-service/internal/app/entity"
	"github.com/AlexanderTal/catalog-service/internal/app/repository/mocks"
	"github.com/AlexanderTal/catalog-service/internal/pkg/testutil"
)

type createProductSuite struct {
	suite.Suite
	srv          *srv
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

func (s *createProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.srv = &srv{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func TestCreateProductSuite(t *testing.T) {
	suite.Run(t, new(createProductSuite))
}

func (s *createProductSuite) TestCreate() {
	type args struct {
		req entity.RequestProductCreate
	}
	type want struct {
		err error
	}

	categoryGUID := uuid.Must(uuid.NewV4())

	testCases := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "success",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "New Product",
					Description:  testutil.PtrString("Good product"),
					Price:        1000,
					CategoryGUID: categoryGUID,
				},
			},
			want: want{err: nil},
			prepare: func(args args) {
				// 1. Мокаем транзакцию, чтобы она просто выполнила callback
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				// 2. Мокаем проверку на дубликат имени (ничего не нашли)
				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
					Return([]entity.Product{}, nil).
					Once()

				// 3. Мокаем проверку категории (нашли)
				s.categoryRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.req.CategoryGUID}).
					Return([]entity.Category{{GUID: categoryGUID}}, nil).
					Once()

				// 4. Мокаем сохранение в базу
				s.productRepo.EXPECT().
					Create(s.ctx, mock.MatchedBy(func(p entity.Product) bool {
						return p.Name == args.req.Name &&
							*p.Description == *args.req.Description &&
							p.Price == args.req.Price &&
							p.CategoryGUID == args.req.CategoryGUID
					})).
					Return(nil).
					Once()
			},
		},
		{
			name: "already exists",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "Existing Product",
					Price:        500,
					CategoryGUID: categoryGUID,
				},
			},
			want: want{err: entity.ErrAlreadyExists},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
					Return([]entity.Product{{Name: "Existing Product"}}, nil).
					Once()
			},
		},
		{
			name: "category not found",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "New Product",
					Price:        1000,
					CategoryGUID: categoryGUID,
				},
			},
			want: want{err: entity.ErrNotFound},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
					Return([]entity.Product{}, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.req.CategoryGUID}).
					Return([]entity.Category{}, nil). // пустой массив = не нашли
					Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.args)

			result, err := s.srv.Create(s.ctx, tc.args.req)

			if tc.want.err != nil {
				s.ErrorIs(err, tc.want.err)
				s.Empty(result.GUID)
			} else {
				s.NoError(err)
				s.NotEmpty(result.GUID)
				s.Equal(tc.args.req.Name, result.Name)
				s.Equal(tc.args.req.Description, result.Description)
				s.Equal(tc.args.req.Price, result.Price)
				s.Equal(tc.args.req.CategoryGUID, result.CategoryGUID)
			}
		})
	}
}

type getByGUIDsProductSuite struct {
	suite.Suite
	srv          *srv
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

func (s *getByGUIDsProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.srv = &srv{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func TestGetByGUIDsProductSuite(t *testing.T) {
	suite.Run(t, new(getByGUIDsProductSuite))
}

func (s *getByGUIDsProductSuite) TestGetByGUIDs() {
	productGUID := uuid.Must(uuid.NewV4())

	type args struct {
		guids []uuid.UUID
	}
	type want struct {
		products []entity.Product
	}

	testCases := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "single product found",
			args: args{
				guids: []uuid.UUID{productGUID},
			},
			want: want{
				products: []entity.Product{{GUID: productGUID, Name: "Test Product"}},
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, args.guids).
					Return([]entity.Product{{GUID: productGUID, Name: "Test Product"}}, nil).
					Once()
			},
		},
		{
			name: "not found returns empty slice",
			args: args{
				guids: []uuid.UUID{uuid.Must(uuid.NewV4())},
			},
			want: want{
				products: []entity.Product{},
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, args.guids).
					Return([]entity.Product{}, nil).
					Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.args)

			result, err := s.srv.GetByGUIDs(s.ctx, tc.args.guids)

			s.NoError(err)
			s.Equal(len(tc.want.products), len(result))
			if len(tc.want.products) > 0 {
				s.Equal(tc.want.products[0].GUID, result[0].GUID)
				s.Equal(tc.want.products[0].Name, result[0].Name)
			}
		})
	}
}

type deleteProductSuite struct {
	suite.Suite
	srv          *srv
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

func (s *deleteProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.srv = &srv{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func TestDeleteProductSuite(t *testing.T) {
	suite.Run(t, new(deleteProductSuite))
}

func (s *deleteProductSuite) TestDelete() {
	productGUID := uuid.Must(uuid.NewV4())

	type args struct {
		guid uuid.UUID
	}
	type want struct {
		err error
	}

	testCases := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "success",
			args: args{guid: productGUID},
			want: want{err: nil},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).
					Return([]entity.Product{{GUID: args.guid}}, nil).
					Once()

				s.productRepo.EXPECT().
					Delete(s.ctx, args.guid).
					Return(nil).
					Once()
			},
		},
		{
			name: "not found",
			args: args{guid: productGUID},
			want: want{err: entity.ErrNotFound},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).
					Return([]entity.Product{}, nil). // Пустой массив = не нашли
					Once()
			},
		},
		{
			name: "delete error",
			args: args{guid: productGUID},
			want: want{err: errors.New("db error")},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).
					Return([]entity.Product{{GUID: args.guid}}, nil).
					Once()

				s.productRepo.EXPECT().
					Delete(s.ctx, args.guid).
					Return(errors.New("db error")).
					Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.args)

			err := s.srv.Delete(s.ctx, tc.args.guid)

			if tc.want.err != nil {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

type listProductSuite struct {
	suite.Suite
	srv          *srv
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

func (s *listProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.srv = &srv{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func TestListProductSuite(t *testing.T) {
	suite.Run(t, new(listProductSuite))
}

func (s *listProductSuite) TestList() {
	type args struct {
		req entity.RequestProductList
	}
	type want struct {
		products []entity.Product
	}

	testCases := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "success",
			args: args{
				req: entity.RequestProductList{},
			},
			want: want{
				products: []entity.Product{{Name: "Product 1"}},
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					List(s.ctx, (*string)(nil), args.req.CategoryGUID, args.req.MinPrice, args.req.MaxPrice).
					Return([]entity.Product{{Name: "Product 1"}}, nil).
					Once()
			},
		},
		{
			name: "empty result",
			args: args{
				req: entity.RequestProductList{},
			},
			want: want{
				products: []entity.Product{},
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					List(s.ctx, (*string)(nil), args.req.CategoryGUID, args.req.MinPrice, args.req.MaxPrice).
					Return([]entity.Product{}, nil).
					Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.args)

			result, err := s.srv.List(s.ctx, tc.args.req)

			s.NoError(err)
			s.Equal(len(tc.want.products), len(result))
		})
	}
}

type updateProductSuite struct {
	suite.Suite
	srv          *srv
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

func (s *updateProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.srv = &srv{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func TestUpdateProductSuite(t *testing.T) {
	suite.Run(t, new(updateProductSuite))
}

func (s *updateProductSuite) TestUpdate() {
	productGUID := uuid.Must(uuid.NewV4())
	oldCategoryGUID := uuid.Must(uuid.NewV4())
	newCategoryGUID := uuid.Must(uuid.NewV4())

	type args struct {
		guid uuid.UUID
		req  entity.RequestProductUpdate
	}
	type want struct {
		err error
	}

	testCases := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "full update",
			args: args{
				guid: productGUID,
				req: entity.RequestProductUpdate{
					Name:         "Updated Name",
					Description:  testutil.PtrString("Updated Desc"),
					Price:        2000,
					CategoryGUID: newCategoryGUID,
				},
			},
			want: want{err: nil},
			prepare: func(args args) {
				s.productRepo.EXPECT().InsideTx(s.ctx, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }).Once()

				s.productRepo.EXPECT().GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).Return([]entity.Product{{
					GUID: args.guid, Name: "Old Name", Price: 1000, CategoryGUID: oldCategoryGUID,
				}}, nil).Once()

				s.productRepo.EXPECT().List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).Return([]entity.Product{}, nil).Once()
				s.categoryRepo.EXPECT().GetByGUIDs(s.ctx, []uuid.UUID{newCategoryGUID}).Return([]entity.Category{{GUID: newCategoryGUID}}, nil).Once()

				s.productRepo.EXPECT().Update(s.ctx, mock.MatchedBy(func(p entity.Product) bool {
					return p.Name == args.req.Name && p.Price == args.req.Price && *p.Description == *args.req.Description && p.CategoryGUID == newCategoryGUID
				})).Return(nil).Once()
			},
		},
		{
			name: "partial update - name only",
			args: args{
				guid: productGUID,
				req: entity.RequestProductUpdate{
					Name: "New Name",
				},
			},
			want: want{err: nil},
			prepare: func(args args) {
				s.productRepo.EXPECT().InsideTx(s.ctx, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }).Once()
				s.productRepo.EXPECT().GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).Return([]entity.Product{{GUID: args.guid, Name: "Old Name"}}, nil).Once()

				s.productRepo.EXPECT().List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).Return([]entity.Product{}, nil).Once()

				s.productRepo.EXPECT().Update(s.ctx, mock.MatchedBy(func(p entity.Product) bool {
					return p.Name == args.req.Name
				})).Return(nil).Once()
			},
		},
		{
			name: "not found",
			args: args{
				guid: productGUID,
				req:  entity.RequestProductUpdate{Name: "New Name"},
			},
			want: want{err: entity.ErrNotFound},
			prepare: func(args args) {
				s.productRepo.EXPECT().InsideTx(s.ctx, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }).Once()
				s.productRepo.EXPECT().GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).Return([]entity.Product{}, nil).Once()
			},
		},
		{
			name: "duplicate name",
			args: args{
				guid: productGUID,
				req:  entity.RequestProductUpdate{Name: "Existing Name"},
			},
			want: want{err: entity.ErrAlreadyExists},
			prepare: func(args args) {
				s.productRepo.EXPECT().InsideTx(s.ctx, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }).Once()
				s.productRepo.EXPECT().GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).Return([]entity.Product{{GUID: args.guid, Name: "Old Name"}}, nil).Once()

				otherGUID := uuid.Must(uuid.NewV4())
				s.productRepo.EXPECT().List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).Return([]entity.Product{{GUID: otherGUID, Name: "Existing Name"}}, nil).Once()
			},
		},
		{
			name: "category not found",
			args: args{
				guid: productGUID,
				req: entity.RequestProductUpdate{
					CategoryGUID: newCategoryGUID,
				},
			},
			want: want{err: entity.ErrNotFound},
			prepare: func(args args) {
				s.productRepo.EXPECT().InsideTx(s.ctx, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }).Once()
				s.productRepo.EXPECT().GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).Return([]entity.Product{{GUID: args.guid, Name: "Old Name"}}, nil).Once()

				s.categoryRepo.EXPECT().GetByGUIDs(s.ctx, []uuid.UUID{newCategoryGUID}).Return([]entity.Category{}, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.args)
			_, err := s.srv.Update(s.ctx, tc.args.guid, tc.args.req)
			if tc.want.err != nil {
				s.ErrorIs(err, tc.want.err)
			} else {
				s.NoError(err)
			}
		})
	}
}
