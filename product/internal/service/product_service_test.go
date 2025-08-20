package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/elangreza/e-commerce/product/internal/entity"
	"github.com/elangreza/e-commerce/product/internal/mockjson"
	"github.com/elangreza/e-commerce/product/internal/service/mock"
	params "github.com/elangreza/e-commerce/product/params"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func Test_productService_ListProducts(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// table test driven approach
	tests := []struct {
		name     string
		req      params.PaginationRequest
		want     *params.ListProductsResponse
		wantErr  bool
		funcMock func(m *mock.MockproductRepo)
	}{
		{
			name: "List Products Success",
			req: params.PaginationRequest{
				Search: "test",
				Page:   1,
				Limit:  10,
				SortBy: "name",
			},
			want: &params.ListProductsResponse{
				Products: []params.ProductResponse{
					{
						ID:          uuid.Nil.String(),
						Name:        "Test Product",
						Description: "This is a test product",
						Price:       99.99,
						Picture:     "http://example.com/image.jpg",
					},
				},
				Total:      1,
				TotalPages: 1,
			},
			wantErr: false,
			funcMock: func(m *mock.MockproductRepo) {
				reqParams := entity.ListProductRequest{
					Search: "test",
					Page:   1,
					Limit:  10,
					SortBy: "name",
				}
				m.EXPECT().ListProducts(gomock.Any(), reqParams).Return(
					[]entity.Product{
						{
							ID:          uuid.Nil,
							Name:        "Test Product",
							Description: "This is a test product",
							Price:       99.99,
							Picture:     "http://example.com/image.jpg",
						},
					},
					nil,
				).Times(1)

				m.EXPECT().TotalProducts(gomock.Any(), reqParams).Return(int64(1), nil).Times(1)

			},
		},
		{
			name: "Error when getting products",
			req: params.PaginationRequest{
				Search: "error",
				Page:   1,
				Limit:  10,
				SortBy: "name"},
			want:    nil,
			wantErr: true,
			funcMock: func(m *mock.MockproductRepo) {
				reqParams := entity.ListProductRequest{
					Search: "error",
					Page:   1,
					Limit:  10,
					SortBy: "name",
				}
				m.EXPECT().ListProducts(gomock.Any(), reqParams).Return(nil, errors.New("test err")).Times(1)
			},
		},
		{
			name: "Error when getting total products",
			req: params.PaginationRequest{
				Search: "error",
				Page:   1,
				Limit:  10,
				SortBy: "name"},
			want:    nil,
			wantErr: true,
			funcMock: func(m *mock.MockproductRepo) {
				reqParams := entity.ListProductRequest{
					Search: "error",
					Page:   1,
					Limit:  10,
					SortBy: "name",
				}
				m.EXPECT().ListProducts(gomock.Any(), reqParams).Return(
					[]entity.Product{
						{
							ID:          uuid.Nil,
							Name:        "Test Product",
							Description: "This is a test product",
							Price:       99.99,
							Picture:     "http://example.com/image.jpg",
						},
					},
					nil,
				).Times(1)
				m.EXPECT().TotalProducts(gomock.Any(), reqParams).Return(int64(0), errors.New("test err")).Times(1)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := mock.NewMockproductRepo(mockCtrl)
			service := NewProductService(mockRepo)

			test.funcMock(mockRepo)

			got, err := service.ListProducts(context.Background(), test.req)

			if (err != nil) != test.wantErr {
				t.Errorf("ListProducts() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListProducts() got = %v, want %v", got, test.want)
			}
		})
	}

}

func Test_productService_GetProduct(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// table test driven approach
	tests := []struct {
		name     string
		req      params.GetProductRequest
		want     *params.GetProductResponse
		wantErr  bool
		funcMock func(m *mock.MockproductRepo)
	}{
		{
			name: "Get Product Success",
			req: params.GetProductRequest{
				ProductID: uuid.Nil.String(),
			},
			want: &params.GetProductResponse{
				Product: &params.ProductResponse{
					ID:          uuid.Nil.String(),
					Name:        "Test Product",
					Description: "This is a test product",
					Price:       99.99,
					Picture:     "http://example.com/image.jpg",
				},
			},
			wantErr: false,
			funcMock: func(m *mock.MockproductRepo) {
				m.EXPECT().GetProductByID(gomock.Any(), uuid.Nil).Return(&entity.Product{
					ID:          uuid.Nil,
					Name:        "Test Product",
					Description: "This is a test product",
					Price:       99.99,
					Picture:     "http://example.com/image.jpg",
				}, nil).Times(1)
			},
		},
		{
			name: "Get Product Not Found",
			req: params.GetProductRequest{
				ProductID: uuid.Nil.String(),
			},
			want:    nil,
			wantErr: true,
			funcMock: func(m *mock.MockproductRepo) {
				m.EXPECT().GetProductByID(gomock.Any(), uuid.Nil).Return(nil, mockjson.DataNotFound).Times(1)
			},
		},
		{
			name: "Get Product Error",
			req: params.GetProductRequest{
				ProductID: uuid.Nil.String(),
			},
			want:    nil,
			wantErr: true,
			funcMock: func(m *mock.MockproductRepo) {
				m.EXPECT().GetProductByID(gomock.Any(), uuid.Nil).Return(nil, errors.New("test err")).Times(1)
			},
		},
		{
			name: "error when parsing uuid",
			req: params.GetProductRequest{
				ProductID: "test",
			},
			want:     nil,
			wantErr:  true,
			funcMock: func(m *mock.MockproductRepo) {},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := mock.NewMockproductRepo(mockCtrl)
			service := NewProductService(mockRepo)

			test.funcMock(mockRepo)

			got, err := service.GetProduct(context.Background(), test.req)

			if (err != nil) != test.wantErr {
				t.Errorf("GetProduct() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetProduct() got = %v, want %v", got, test.want)
			}
		})
	}
}
