package service

import (
	"context"

	"github.com/elangreza/e-commerce/shop/internal/entity"

	"github.com/elangreza/e-commerce/gen"
)

type (
	ShopRepo interface {
		GetShopByIDs(ctx context.Context, IDs ...int64) ([]entity.Shop, error)
	}

	warehouseClient interface {
		GetWarehouseByShopID(ctx context.Context, shopID int64) (*gen.GetWarehouseByShopIDResponse, error)
	}

	ShopService struct {
		repo            ShopRepo
		warehouseClient warehouseClient
		gen.UnimplementedShopServiceServer
	}
)

func NewShopService(repo ShopRepo, warehouseClient warehouseClient) *ShopService {
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
			wRes, err := s.warehouseClient.GetWarehouseByShopID(ctx, shop.ID)
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
