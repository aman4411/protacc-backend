package models

import (
	"time"

	"github.com/gosimple/slug"
)

type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
)

// Post is an admin-authored blog article.
type Post struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Excerpt     string     `json:"excerpt"`
	Content     string     `json:"content"` // sanitized HTML
	CoverImage  string     `json:"cover_image"`
	Category    string     `json:"category"`
	Tags        []string   `json:"tags"`
	Status      PostStatus `json:"status"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (p *Post) GenerateSlug() {
	p.Slug = slug.Make(p.Title)
}
