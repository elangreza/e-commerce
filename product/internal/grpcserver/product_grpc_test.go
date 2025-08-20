package grpcserver

import (
	"context"
	"errors"
	"net"
	"testing"

	pb "github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/product/internal/grpcserver/mock"
	"github.com/elangreza/e-commerce/product/params"
	"github.com/elangreza/e-commerce/product/pkg/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func TestProductServer_GetProduct(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	listener := bufconn.Listen(1024 * 1024)
	productService := mock.NewMockproductService(ctrl)

	s := grpc.NewServer()
	server := NewProductServer(productService)
	pb.RegisterProductServiceServer(s, server)

	go func() {
		s.Serve(listener)
	}()
	defer s.Stop()

	type args struct {
		req *pb.GetProductRequest
	}
	tests := []struct {
		name     string
		args     args
		want     *pb.Product
		funcMock func(m *mock.MockproductService)
		wantErr  bool
		wantCode codes.Code
	}{
		{
			name: "success",
			args: args{
				req: &pb.GetProductRequest{
					Id: uuid.Nil.String(),
				},
			},
			want: &pb.Product{
				Id:          uuid.Nil.String(),
				Name:        "1",
				Description: "1",
				Picture:     "1",
				Price:       1,
			},
			funcMock: func(m *mock.MockproductService) {
				m.EXPECT().GetProduct(gomock.Any(), params.GetProductRequest{
					ProductID: uuid.Nil.String(),
				}).Return(&params.GetProductResponse{
					Product: &params.ProductResponse{
						ID:          uuid.Nil.String(),
						Name:        "1",
						Description: "1",
						Price:       1,
						Picture:     "1",
					},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "error not found",
			args: args{
				req: &pb.GetProductRequest{
					Id: uuid.Nil.String(),
				},
			},
			want: nil,
			funcMock: func(m *mock.MockproductService) {
				m.EXPECT().GetProduct(gomock.Any(), params.GetProductRequest{
					ProductID: uuid.Nil.String(),
				}).Return(nil, errs.NotFound{})
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name: "error internal server",
			args: args{
				req: &pb.GetProductRequest{
					Id: uuid.Nil.String(),
				},
			},
			want: nil,
			funcMock: func(m *mock.MockproductService) {
				m.EXPECT().GetProduct(gomock.Any(), params.GetProductRequest{
					ProductID: uuid.Nil.String(),
				}).Return(nil, errors.New("test error"))
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Set up mock before starting server
			tt.funcMock(productService)

			conn, err := grpc.DialContext(context.Background(), "bufnet",
				grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
					return listener.Dial()
				}),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			assert.NoError(t, err)
			defer conn.Close()

			client := pb.NewProductServiceClient(conn)
			res, err := client.GetProduct(context.Background(), tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, res)

				st, ok := status.FromError(err)
				if ok {
					code := st.Code()
					assert.Equal(t, code, tt.wantCode)
				}
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, res)
			assert.Equal(t, tt.want.Id, res.Id)
			assert.Equal(t, tt.want.Name, res.Name)
			assert.Equal(t, tt.want.Description, res.Description)
			assert.Equal(t, tt.want.Picture, res.Picture)
			assert.Equal(t, tt.want.Price, res.Price)
		})
	}
}

func TestProductServer_ListProducts(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	listener := bufconn.Listen(1024 * 1024)
	productService := mock.NewMockproductService(ctrl)

	s := grpc.NewServer()
	server := NewProductServer(productService)
	pb.RegisterProductServiceServer(s, server)

	go func() {
		s.Serve(listener)
	}()
	defer s.Stop()

	type args struct {
		req *pb.ListProductsRequest
	}
	tests := []struct {
		name     string
		args     args
		funcMock func(m *mock.MockproductService)
		wantErr  bool
		wantCode codes.Code
	}{
		{
			name: "success",
			args: args{
				req: &pb.ListProductsRequest{
					Search: "",
					Limit:  0,
					Page:   0,
					SortBy: "",
				},
			},
			funcMock: func(m *mock.MockproductService) {
				m.EXPECT().ListProducts(gomock.Any(), params.PaginationRequest{
					Search: "",
					Page:   1,
					Limit:  10,
					SortBy: "updated_at",
				}).Return(&params.ListProductsResponse{
					Products: []params.ProductResponse{
						{
							ID:          "1",
							Name:        "1",
							Description: "1",
							Price:       1,
							Picture:     "1",
						},
					},
					Total:      1,
					TotalPages: 10,
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				req: &pb.ListProductsRequest{
					Search: "",
					Limit:  0,
					Page:   0,
					SortBy: "",
				},
			},
			funcMock: func(m *mock.MockproductService) {
				m.EXPECT().ListProducts(gomock.Any(), params.PaginationRequest{
					Search: "",
					Page:   1,
					Limit:  10,
					SortBy: "updated_at",
				}).Return(nil, errors.New("test error"))
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
		{
			name: "error invalid sort",
			args: args{
				req: &pb.ListProductsRequest{
					Search: "",
					Limit:  0,
					Page:   0,
					SortBy: "picture",
				},
			},
			funcMock: func(m *mock.MockproductService) {},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Set up mock before starting server
			tt.funcMock(productService)

			conn, err := grpc.DialContext(context.Background(), "bufnet",
				grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
					return listener.Dial()
				}),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			assert.NoError(t, err)
			defer conn.Close()

			client := pb.NewProductServiceClient(conn)
			res, err := client.ListProducts(context.Background(), tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, res)

				st, ok := status.FromError(err)
				if ok {
					code := st.Code()
					assert.Equal(t, code, tt.wantCode)
				}
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, res)
		})
	}
}
