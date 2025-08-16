package service

import (
	"fmt"
	"net"
	"strings"

	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/repository"
	"github.com/go-playground/validator/v10"
)

type ContactService interface {
	CreateContact(contact *models.CreateContactRequest, ipAddress, userAgent string) (*models.ContactMessage, error)
	GetContacts(filters models.ContactFilters) ([]models.ContactMessage, int, error)
	GetContactByID(id int) (*models.ContactMessage, error)
	UpdateContactStatus(id int, updates models.UpdateContactRequest) error
	DeleteContact(id int) error
	GetContactStats() (*models.ContactStats, error)
}

type contactService struct {
	repo      repository.ContactRepository
	validator *validator.Validate
}

func NewContactService(repo repository.ContactRepository) ContactService {
	return &contactService{
		repo:      repo,
		validator: validator.New(),
	}
}

func (s *contactService) CreateContact(contact *models.CreateContactRequest, ipAddress, userAgent string) (*models.ContactMessage, error) {
	// Validate the contact request
	if err := s.validator.Struct(contact); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Sanitize inputs
	contact.Name = strings.TrimSpace(contact.Name)
	contact.Email = strings.TrimSpace(strings.ToLower(contact.Email))
	contact.Phone = strings.TrimSpace(contact.Phone)
	contact.Subject = strings.TrimSpace(contact.Subject)
	contact.Message = strings.TrimSpace(contact.Message)

	if contact.Company != nil {
		trimmed := strings.TrimSpace(*contact.Company)
		if trimmed == "" {
			contact.Company = nil
		} else {
			contact.Company = &trimmed
		}
	}

	if contact.ServiceInterest != nil {
		trimmed := strings.TrimSpace(*contact.ServiceInterest)
		if trimmed == "" {
			contact.ServiceInterest = nil
		} else {
			contact.ServiceInterest = &trimmed
		}
	}

	// Validate IP address format
	if ipAddress != "" {
		if net.ParseIP(ipAddress) == nil {
			ipAddress = "" // Invalid IP, set to empty
		}
	}

	// Sanitize user agent
	if len(userAgent) > 500 {
		userAgent = userAgent[:500] // Truncate if too long
	}

	// Create the contact message
	result, err := s.repo.CreateContact(contact, ipAddress, userAgent)
	if err != nil {
		return nil, fmt.Errorf("failed to create contact message: %w", err)
	}

	return result, nil
}

func (s *contactService) GetContacts(filters models.ContactFilters) ([]models.ContactMessage, int, error) {
	// Validate and set defaults for pagination
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.Limit <= 0 || filters.Limit > 100 {
		filters.Limit = 10
	}

	// Sanitize search term
	if filters.Search != "" {
		filters.Search = strings.TrimSpace(filters.Search)
		if len(filters.Search) > 100 {
			filters.Search = filters.Search[:100]
		}
	}

	// Validate status filter
	if filters.Status != "" {
		validStatuses := map[string]bool{"new": true, "replied": true, "resolved": true}
		if !validStatuses[filters.Status] {
			return nil, 0, fmt.Errorf("invalid status filter: %s", filters.Status)
		}
	}

	return s.repo.GetContacts(filters)
}

func (s *contactService) GetContactByID(id int) (*models.ContactMessage, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid contact ID")
	}

	return s.repo.GetContactByID(id)
}

func (s *contactService) UpdateContactStatus(id int, updates models.UpdateContactRequest) error {
	if id <= 0 {
		return fmt.Errorf("invalid contact ID")
	}

	// Validate the update request
	if err := s.validator.Struct(updates); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Sanitize responded_by if provided
	if updates.RespondedBy != nil {
		trimmed := strings.TrimSpace(*updates.RespondedBy)
		if trimmed == "" {
			updates.RespondedBy = nil
		} else {
			updates.RespondedBy = &trimmed
		}
	}

	return s.repo.UpdateContactStatus(id, updates)
}

func (s *contactService) DeleteContact(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid contact ID")
	}

	return s.repo.DeleteContact(id)
}

func (s *contactService) GetContactStats() (*models.ContactStats, error) {
	return s.repo.GetContactStats()
}
