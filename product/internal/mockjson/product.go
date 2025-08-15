package mockjson

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"sort"

	"github.com/elangreza/e-commerce/product/internal/entity"
	"github.com/elangreza/e-commerce/product/params"
	"github.com/google/uuid"
)

type ProductMock struct {
	products []entity.Product
}

func LoadProductJson() (*ProductMock, error) {
	data, err := os.ReadFile("internal/mockjson/products.json")
	if err != nil {
		return nil, err
	}

	var products []entity.Product
	if err = json.Unmarshal(data, &products); err != nil {
		return nil, err
	}

	productMock := &ProductMock{
		products: products,
	}

	return productMock, nil
}

func (pm *ProductMock) ListProducts(ctx context.Context, req params.ListProductsRequest) ([]entity.Product, error) {

	var filteredProducts = pm.filteredProducts(pm.products, req.Search)

	if req.SortBy == "" {
		sort.Slice(filteredProducts, func(i, j int) bool {
			switch req.SortBy {
			case "name":
				return filteredProducts[i].Name < filteredProducts[j].Name
			case "price":
				return filteredProducts[i].Price < filteredProducts[j].Price
			case "updated_at":
				return filteredProducts[i].UpdatedAt < filteredProducts[j].UpdatedAt
			default:
				return false // Default case if no valid sort option is provided
			}
		})
	}

	start := (req.Page - 1) * req.Limit
	end := min(start+req.Limit, int64(len(filteredProducts)))

	return filteredProducts[start:end], nil
}

func (pm *ProductMock) TotalProducts(ctx context.Context, req params.ListProductsRequest) (int64, error) {

	filteredProducts := pm.filteredProducts(pm.products, req.Search)

	return int64(len(filteredProducts)), nil
}

func (pm *ProductMock) filteredProducts(products []entity.Product, search string) []entity.Product {
	if search == "" {
		return products
	}

	var filtered []entity.Product
	for _, product := range products {
		if product.Name == search || product.Description == search {
			filtered = append(filtered, product)
		}
	}
	return filtered
}

func (pm *ProductMock) GetProductByID(ctx context.Context, productID uuid.UUID) (*entity.Product, error) {

	for _, product := range pm.products {
		if product.ID == productID {
			return &product, nil
		}
	}
	return nil, sql.ErrNoRows
}
