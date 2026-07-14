package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aman4411/protacc-backend/cache"
	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/repository"
	"github.com/gosimple/slug"
)

// Cache prefix + TTLs for public (published) post reads. Any post write busts
// the whole "posts:" family. Admin reads (List/GetByID) are never cached.
const (
	postCachePrefix = "posts:"
	postCacheTTL    = 60 * time.Minute
	postCatCacheTTL = 60 * time.Minute
)

type PostService struct {
	repo  *repository.PostRepository
	cache *cache.Cache
}

func NewPostService(repo *repository.PostRepository, c *cache.Cache) *PostService {
	return &PostService{repo: repo, cache: c}
}

func (s *PostService) invalidate() {
	if s.cache != nil {
		s.cache.InvalidatePrefix(postCachePrefix)
	}
}

// listPublishedResult bundles the multi-return of ListPublished for caching.
type listPublishedResult struct {
	Posts []models.Post
	Total int
}

func (s *PostService) ListPublished(ctx context.Context, page, limit int, category string) ([]models.Post, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 9
	}
	category = strings.TrimSpace(category)
	key := fmt.Sprintf("%slist:%d:%d:%s", postCachePrefix, page, limit, category)
	res, err := cache.Load(s.cache, key, postCacheTTL, func() (listPublishedResult, error) {
		posts, total, err := s.repo.ListPublished(ctx, page, limit, category)
		return listPublishedResult{Posts: posts, Total: total}, err
	})
	if err != nil {
		return nil, 0, err
	}
	return res.Posts, res.Total, nil
}

func (s *PostService) ListCategories(ctx context.Context) ([]string, error) {
	return cache.Load(s.cache, postCachePrefix+"categories", postCatCacheTTL, func() ([]string, error) {
		return s.repo.ListCategories(ctx)
	})
}

func (s *PostService) GetPublishedBySlug(ctx context.Context, slugStr string) (*models.Post, error) {
	return cache.Load(s.cache, postCachePrefix+"slug:"+slugStr, postCacheTTL, func() (*models.Post, error) {
		return s.repo.GetPublishedBySlug(ctx, slugStr)
	})
}

func (s *PostService) List(ctx context.Context) ([]models.Post, error) {
	return s.repo.List(ctx)
}

func (s *PostService) GetByID(ctx context.Context, id int) (*models.Post, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PostService) normalize(p *models.Post) error {
	p.Title = strings.TrimSpace(p.Title)
	if p.Title == "" {
		return fmt.Errorf("title is required")
	}
	if strings.TrimSpace(p.Slug) == "" {
		p.GenerateSlug()
	} else {
		p.Slug = slug.Make(p.Slug)
	}
	if p.Status != models.PostStatusPublished && p.Status != models.PostStatusDraft {
		p.Status = models.PostStatusDraft
	}
	// Normalize tags: lowercase + trim + dedupe so tag-overlap matching is reliable.
	seen := map[string]bool{}
	cleaned := []string{}
	for _, t := range p.Tags {
		t = strings.ToLower(strings.TrimSpace(t))
		if t == "" || seen[t] {
			continue
		}
		seen[t] = true
		cleaned = append(cleaned, t)
	}
	p.Tags = cleaned
	return nil
}

// ListRelated returns posts related to the given published post by shared tags.
func (s *PostService) ListRelated(ctx context.Context, excludeID int, tags []string, limit int) ([]models.Post, error) {
	if limit < 1 || limit > 12 {
		limit = 3
	}
	key := fmt.Sprintf("%srelated:%d:%d", postCachePrefix, excludeID, limit)
	return cache.Load(s.cache, key, postCacheTTL, func() ([]models.Post, error) {
		return s.repo.ListRelated(ctx, excludeID, tags, limit)
	})
}

func (s *PostService) Create(ctx context.Context, p *models.Post) (*models.Post, error) {
	if err := s.normalize(p); err != nil {
		return nil, err
	}
	if p.Status == models.PostStatusPublished && p.PublishedAt == nil {
		now := time.Now()
		p.PublishedAt = &now
	}
	res, err := s.repo.Create(ctx, p)
	if err == nil {
		s.invalidate()
	}
	return res, err
}

func (s *PostService) Update(ctx context.Context, p *models.Post) (*models.Post, error) {
	if err := s.normalize(p); err != nil {
		return nil, err
	}
	existing, err := s.repo.GetByID(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("post not found")
	}
	// Preserve the original publish date; set it the first time it's published.
	if p.Status == models.PostStatusPublished {
		if existing.PublishedAt != nil {
			p.PublishedAt = existing.PublishedAt
		} else {
			now := time.Now()
			p.PublishedAt = &now
		}
	} else {
		p.PublishedAt = existing.PublishedAt
	}
	res, err := s.repo.Update(ctx, p)
	if err == nil {
		s.invalidate()
	}
	return res, err
}

func (s *PostService) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err == nil {
		s.invalidate()
	}
	return err
}
