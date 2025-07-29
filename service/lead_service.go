package service

import (
	"context"

	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/repository"
)

type LeadService struct {
	leadRepo *repository.LeadRepository
}

func NewLeadService(leadRepo *repository.LeadRepository) *LeadService {
	return &LeadService{
		leadRepo: leadRepo,
	}
}

func (s *LeadService) CreateLead(ctx context.Context, lead *models.CreateLeadRequest) (*models.BusinessLead, error) {
	return s.leadRepo.CreateLead(ctx, lead)
}

func (s *LeadService) GetLeads(ctx context.Context, filters models.LeadFilters) ([]models.BusinessLead, int, error) {
	return s.leadRepo.GetLeads(ctx, filters)
}

func (s *LeadService) GetLeadByID(ctx context.Context, id int) (*models.BusinessLead, error) {
	return s.leadRepo.GetLeadByID(ctx, id)
}

func (s *LeadService) UpdateLead(ctx context.Context, id int, updates *models.UpdateLeadRequest) error {
	return s.leadRepo.UpdateLead(ctx, id, updates)
}

func (s *LeadService) DeleteLead(ctx context.Context, id int) error {
	return s.leadRepo.DeleteLead(ctx, id)
}

func (s *LeadService) GetLeadStats(ctx context.Context) (*models.LeadStats, error) {
	return s.leadRepo.GetLeadStats(ctx)
}
