package service

import (
	"context"

	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/repository"
)

type CartService struct {
	repo *repository.CartRepository
}

func NewCartService(repo *repository.CartRepository) *CartService {
	return &CartService{
		repo: repo,
	}
}

func (s *CartService) AddToCart(ctx context.Context, userID string, serviceID int) error {
	return s.repo.AddToCart(ctx, userID, serviceID)
}

func (s *CartService) GetCartItems(ctx context.Context, userID string) ([]models.CartItem, error) {
	return s.repo.GetCartItems(ctx, userID)
}

func (s *CartService) RemoveFromCart(ctx context.Context, userID string, serviceID int) error {
	return s.repo.RemoveFromCart(ctx, userID, serviceID)
}

func (s *CartService) ClearCart(ctx context.Context, userID string) error {
	return s.repo.ClearCart(ctx, userID)
}
