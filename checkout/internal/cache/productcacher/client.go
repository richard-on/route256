package productcacher

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/cache"
)

// GetProductCache is a client for productService that uses cache.
type GetProductCache struct {
	// directClient is a default client for productService.
	directClient domain.ProductLister
	// cachedClient is a cache.
	cachedClient cache.Cacher[uint32, model.ProductInfo]
}

// NewCachedProductServiceClient returns a new instance of GetProductCache.
func NewCachedProductServiceClient(directClient domain.ProductLister,
	cachedClient cache.Cacher[uint32, model.ProductInfo]) *GetProductCache {
	return &GetProductCache{
		directClient: directClient,
		cachedClient: cachedClient,
	}
}

// GetProduct calls productService.GetProduct for a given product and returns domain.ProductInfo.
func (c *GetProductCache) GetProduct(ctx context.Context, sku uint32) (model.ProductInfo, error) {
	info, ok := c.getFromCache(sku)
	if !ok {
		return c.request(ctx, sku)
	}

	return info, nil
}

// request calls productService.GetProduct for a given product,
// saves a response in cache and returns domain.ProductInfo.
func (c *GetProductCache) request(ctx context.Context, sku uint32) (model.ProductInfo, error) {
	resp, err := c.directClient.GetProduct(ctx, sku)
	if err != nil {
		return model.ProductInfo{}, err
	}

	go c.cachedClient.Set(sku, resp)

	return resp, nil
}

// getFromCache returns a value from cache.
func (c *GetProductCache) getFromCache(sku uint32) (model.ProductInfo, bool) {
	resp, ok := c.cachedClient.Get(sku)
	if !ok {
		return model.ProductInfo{}, false
	}

	return resp, true
}
