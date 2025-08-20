package grpcserver

// go generate
//go:generate mockgen -source=product_grpc.go -destination=./mock/mock_product_grpc.go -package=mock

import (
	"context"
	"errors"

	pb "github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/product/params"
	"github.com/elangreza/e-commerce/product/pkg/errs"
	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type (
	productService interface {
		ListProducts(ctx context.Context, req params.PaginationRequest) (*params.ListProductsResponse, error)
		GetProduct(ctx context.Context, req params.GetProductRequest) (*params.GetProductResponse, error)
	}

	ProductServer struct {
		productService productService
		pb.UnimplementedProductServiceServer
	}
)

func NewProductServer(productService productService) *ProductServer {
	return &ProductServer{
		productService: productService,
	}
}

func (s *ProductServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	reqParams := &params.PaginationRequest{
		Search: req.GetSearch(),
		Page:   req.GetPage(),
		Limit:  req.GetLimit(),
		SortBy: req.GetSortBy(),
	}

	if err := reqParams.Validate("updated_at"); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	products, err := s.productService.ListProducts(ctx, *reqParams)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	productResponses := make([]*pb.Product, len(products.Products))
	for i, product := range products.Products {
		productResponses[i] = &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Picture:     product.Picture,
		}
	}

	return &pb.ListProductsResponse{
		Products:   productResponses,
		Total:      products.Total,
		TotalPages: products.TotalPages,
	}, nil
}

func (s *ProductServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	response, err := s.productService.GetProduct(ctx, params.GetProductRequest{
		ProductID: req.GetId(),
	})
	if err != nil {
		if errors.Is(err, errs.NotFound{}) {
			return nil, status.Error(errs.NotFound{}.GrpcCode(), err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Product{
		Id:          response.Product.ID,
		Name:        response.Product.Name,
		Description: response.Product.Description,
		Price:       response.Product.Price,
		Picture:     response.Product.Picture,
	}, nil
}
