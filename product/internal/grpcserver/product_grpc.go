package grpcserver

// go generate
//go:generate mockgen -source=product_grpc.go -destination=./mock/mock_product_grpc.go -package=mock

import (
	"context"
	"errors"
	"strings"

	"github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/product/internal/params"
	"github.com/elangreza/e-commerce/product/pkg/errs"
	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type (
	productService interface {
		ListProducts(ctx context.Context, req params.PaginationParams) (*params.ListProductsResponse, error)
		GetProduct(ctx context.Context, req params.GetProductRequest) (*params.GetProductResponse, error)
	}

	ProductServer struct {
		productService productService
		gen.UnimplementedProductServiceServer
	}
)

func NewProductServer(productService productService) *ProductServer {
	return &ProductServer{
		productService: productService,
	}
}

func (s *ProductServer) ListProducts(ctx context.Context, req *gen.ListProductsRequest) (*gen.ListProductsResponse, error) {
	paginationParams := params.PaginationParams{
		Sorts:  strings.Split(req.GetSortBy(), ","),
		Search: req.GetSearch(),
		Limit:  req.GetLimit(),
		Page:   req.GetPage(),
	}

	paginationParams.SetValidSortKey("updated_at", "name", "price")

	if err := paginationParams.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	products, err := s.productService.ListProducts(ctx, paginationParams)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	productResponses := make([]*gen.Product, len(products.Products))
	for i, product := range products.Products {
		productResponses[i] = &gen.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price.ToProto(),
			ImageUrl:    product.ImageUrl,
		}
	}

	return &gen.ListProductsResponse{
		Products:   productResponses,
		Total:      products.Total,
		TotalPages: products.TotalPages,
	}, nil
}

func (s *ProductServer) GetProduct(ctx context.Context, req *gen.GetProductRequest) (*gen.Product, error) {
	response, err := s.productService.GetProduct(ctx, params.GetProductRequest{
		ProductID: req.GetId(),
	})
	if err != nil {
		var notFoundErr errs.NotFound
		if errors.As(err, &notFoundErr) {
			return nil, status.Error(notFoundErr.GrpcCode(), err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.Product{
		Id:          response.Product.ID,
		Name:        response.Product.Name,
		Description: response.Product.Description,
		Price:       response.Product.Price.ToProto(),
		ImageUrl:    response.Product.ImageUrl,
		Stock:       response.Product.Stock,
	}, nil
}
