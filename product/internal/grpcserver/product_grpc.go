package grpcserver

import (
	"context"

	pb "github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/product/params"
)

type (
	productService interface {
		ListProducts(ctx context.Context, req params.ListProductsRequest) (*params.ListProductsResponse, error)
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
	reqParams := params.ListProductsRequest{
		Search: req.GetSearch(),
		Page:   req.GetPage(),
		Limit:  req.GetLimit(),
		SortBy: req.GetSortBy(),
	}

	if err := reqParams.Validate(); err != nil {
		return nil, err
	}

	products, err := s.productService.ListProducts(ctx, reqParams)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return &pb.Product{
		Id:          response.Product.ID,
		Name:        response.Product.Name,
		Description: response.Product.Description,
		Price:       response.Product.Price,
		Picture:     response.Product.Picture,
	}, nil
}
