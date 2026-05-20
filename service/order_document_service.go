package service

import (
	"context"
	"fmt"

	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/repository"
	"github.com/aman4411/protacc-backend/utils"
)

type OrderDocumentService struct {
	docRepo   *repository.OrderDocumentRepository
	orderRepo *repository.OrderRepository
}

func NewOrderDocumentService(docRepo *repository.OrderDocumentRepository, orderRepo *repository.OrderRepository) *OrderDocumentService {
	return &OrderDocumentService{
		docRepo:   docRepo,
		orderRepo: orderRepo,
	}
}

func (s *OrderDocumentService) ListDocuments(ctx context.Context, orderID int, userID, role string) ([]models.OrderDocument, error) {
	if err := s.ensureOrderAccess(ctx, orderID, userID, role); err != nil {
		return nil, err
	}
	return s.docRepo.GetByOrderID(ctx, orderID)
}

func (s *OrderDocumentService) AddUserDocument(ctx context.Context, orderID int, userID string, req *models.AddOrderDocumentRequest) (*models.OrderDocument, error) {
	ownerID, err := s.docRepo.GetOrderOwnerID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if ownerID != userID {
		return nil, fmt.Errorf("access denied")
	}

	status, err := s.docRepo.GetOrderStatus(ctx, orderID)
	if err != nil {
		return nil, err
	}

	allowedStatuses := map[models.OrderStatus]bool{
		models.OrderStatusDocumentsRequired: true,
		models.OrderStatusDocumentsReceived: true,
		models.OrderStatusProcessing:        true,
	}
	if !allowedStatuses[status] {
		return nil, fmt.Errorf("documents cannot be submitted at the current order stage")
	}

	normalizedURL, err := utils.NormalizeGoogleDriveURL(req.DriveURL)
	if err != nil {
		return nil, err
	}

	var notes *string
	if req.Notes != "" {
		notes = &req.Notes
	}

	doc := &models.OrderDocument{
		OrderID:      orderID,
		UploadedBy:   userID,
		DocumentType: models.OrderDocumentTypeUserUpload,
		Title:        req.Title,
		DriveURL:     normalizedURL,
		Notes:        notes,
	}

	if err := s.docRepo.Create(ctx, doc); err != nil {
		return nil, err
	}

	if status == models.OrderStatusDocumentsRequired {
		_ = s.orderRepo.UpdateOrderStatus(ctx, orderID, models.OrderStatusDocumentsReceived,
			"Customer submitted documents via Google Drive", userID)
	}

	return doc, nil
}

func (s *OrderDocumentService) AddAdminDocument(ctx context.Context, orderID int, adminID string, req *models.AddOrderDocumentRequest) (*models.OrderDocument, error) {
	status, err := s.docRepo.GetOrderStatus(ctx, orderID)
	if err != nil {
		return nil, err
	}

	allowedStatuses := map[models.OrderStatus]bool{
		models.OrderStatusDocumentsReceived: true,
		models.OrderStatusInProgress:        true,
		models.OrderStatusPendingFinalPayment: true,
		models.OrderStatusFullPaymentReceived: true,
		models.OrderStatusCompleted:           true,
		models.OrderStatusProcessing:          true,
	}
	if !allowedStatuses[status] {
		return nil, fmt.Errorf("deliverables cannot be shared at the current order stage")
	}

	normalizedURL, err := utils.NormalizeGoogleDriveURL(req.DriveURL)
	if err != nil {
		return nil, err
	}

	var notes *string
	if req.Notes != "" {
		notes = &req.Notes
	}

	doc := &models.OrderDocument{
		OrderID:      orderID,
		UploadedBy:   adminID,
		DocumentType: models.OrderDocumentTypeAdminDelivery,
		Title:        req.Title,
		DriveURL:     normalizedURL,
		Notes:        notes,
	}

	if err := s.docRepo.Create(ctx, doc); err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *OrderDocumentService) EnsureOrderAccess(ctx context.Context, orderID int, userID, role string) error {
	return s.ensureOrderAccess(ctx, orderID, userID, role)
}

func (s *OrderDocumentService) ensureOrderAccess(ctx context.Context, orderID int, userID, role string) error {
	if role == "admin" {
		_, err := s.docRepo.GetOrderOwnerID(ctx, orderID)
		return err
	}
	ownerID, err := s.docRepo.GetOrderOwnerID(ctx, orderID)
	if err != nil {
		return err
	}
	if ownerID != userID {
		return fmt.Errorf("access denied")
	}
	return nil
}

func (s *OrderDocumentService) EnrichDocuments(links []models.OrderDocument) []map[string]interface{} {
	result := make([]map[string]interface{}, len(links))
	for i, doc := range links {
		entry := map[string]interface{}{
			"id":            doc.ID,
			"order_id":      doc.OrderID,
			"uploaded_by":   doc.UploadedBy,
			"document_type": doc.DocumentType,
			"title":         doc.Title,
			"drive_url":     doc.DriveURL,
			"notes":         doc.Notes,
			"created_at":    doc.CreatedAt,
			"updated_at":    doc.UpdatedAt,
			"uploader":      doc.Uploader,
		}
		if parsed, err := utils.ParseGoogleDriveURL(doc.DriveURL); err == nil {
			entry["drive"] = parsed
		}
		result[i] = entry
	}
	return result
}
