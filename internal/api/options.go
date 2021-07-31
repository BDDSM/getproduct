package api

import "github.com/korableg/getproduct/pkg/product/repository"

type EngineOption func(*Engine)

func WithProductRepository(pr *repository.ProductRepository) EngineOption {
	return func(e *Engine) {
		e.repository = pr
	}
}
