package models

import "time"

type OrderDocumentType string

const (
	OrderDocumentTypeUserUpload    OrderDocumentType = "user_upload"
	OrderDocumentTypeAdminDelivery OrderDocumentType = "admin_delivery"
)

type OrderDocument struct {
	ID           int               `json:"id" db:"id"`
	OrderID      int               `json:"order_id" db:"order_id"`
	UploadedBy   string            `json:"uploaded_by" db:"uploaded_by"`
	DocumentType OrderDocumentType `json:"document_type" db:"document_type"`
	Title        string            `json:"title" db:"title"`
	DriveURL     string            `json:"drive_url" db:"drive_url"`
	Notes        *string           `json:"notes,omitempty" db:"notes"`
	CreatedAt    time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at" db:"updated_at"`
	Uploader     *UserResponse     `json:"uploader,omitempty" db:"-"`
}

type AddOrderDocumentRequest struct {
	Title    string `json:"title" validate:"required,min=2,max=255"`
	DriveURL string `json:"driveUrl" validate:"required,url"`
	Notes    string `json:"notes"`
}
