package service

import (
	"context"

	"github.com/elangreza/e-commerce/shop/internal/entity"

	"github.com/elangreza/e-commerce/gen"
)

//go:generate mockgen -source=shop_service.go -destination=mock/mock_shop_service.go -package=mock
//go:generate mockgen -package=mock -destination=mock/mock_deps.go github.com/elangreza/e-commerce/gen WarehouseServiceClient

type (
	ShopRepo interface {
		GetShopByIDs(ctx context.Context, IDs ...int64) ([]entity.Shop, error)
	}

	ShopService struct {
		repo            ShopRepo
		warehouseClient gen.WarehouseServiceClient
		gen.UnimplementedShopServiceServer
	}
)

func NewShopService(repo ShopRepo, warehouseClient gen.WarehouseServiceClient) *ShopService {
	return &ShopService{
		repo:            repo,
		warehouseClient: warehouseClient,
	}
}

func (s *ShopService) GetShops(ctx context.Context, req *gen.GetShopsRequest) (*gen.ShopList, error) {
	shops, err := s.repo.GetShopByIDs(ctx, req.Ids...)
	if err != nil {
		return nil, err
	}

	res := []*gen.Shop{}
	for _, shop := range shops {
		sh := &gen.Shop{
			Id:         shop.ID,
			Name:       shop.Name,
			Warehouses: []*gen.Warehouse{},
		}

		if req.WithWarehouses {
			var err error
			wRes, err := s.warehouseClient.GetWarehouseByShopID(ctx, &gen.GetWarehouseByShopIDRequest{
				ShopId: shop.ID,
			})
			if err != nil {
				return nil, err
			}
			sh.Warehouses = wRes.Warehouses
		}

		res = append(res, sh)
	}
	return &gen.ShopList{
		Shops: res,
	}, nil
}
