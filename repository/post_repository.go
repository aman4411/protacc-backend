package repository

import (
	"context"
	"errors"

	"github.com/aman4411/protacc-backend/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostRepository struct {
	db *pgxpool.Pool
}

func NewPostRepository(db *pgxpool.Pool) *PostRepository {
	return &PostRepository{db: db}
}

const postCols = `id, title, slug, excerpt, content, cover_image, category, tags, status, published_at, created_at, updated_at`

func scanPost(row pgx.Row, p *models.Post) error {
	return row.Scan(&p.ID, &p.Title, &p.Slug, &p.Excerpt, &p.Content, &p.CoverImage,
		&p.Category, &p.Tags, &p.Status, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt)
}

func scanPosts(rows pgx.Rows) ([]models.Post, error) {
	out := []models.Post{}
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Excerpt, &p.Content, &p.CoverImage,
			&p.Category, &p.Tags, &p.Status, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

// ListRelated returns up to `limit` published posts related to the given post —
// prioritising shared tags, then filling with the most recent others.
func (r *PostRepository) ListRelated(ctx context.Context, excludeID int, tags []string, limit int) ([]models.Post, error) {
	result := []models.Post{}
	seen := map[int]bool{excludeID: true}

	if len(tags) > 0 {
		rows, err := r.db.Query(ctx, `SELECT `+postCols+`
			FROM posts WHERE status = 'published' AND id != $1 AND tags && $2
			ORDER BY published_at DESC NULLS LAST, created_at DESC LIMIT $3`, excludeID, tags, limit)
		if err != nil {
			return nil, err
		}
		matches, err := scanPosts(rows)
		rows.Close()
		if err != nil {
			return nil, err
		}
		for _, p := range matches {
			result = append(result, p)
			seen[p.ID] = true
		}
	}

	// Fill remaining slots with the latest other published posts.
	if len(result) < limit {
		rows, err := r.db.Query(ctx, `SELECT `+postCols+`
			FROM posts WHERE status = 'published' AND id != $1
			ORDER BY published_at DESC NULLS LAST, created_at DESC LIMIT $2`, excludeID, limit+len(result)+1)
		if err != nil {
			return nil, err
		}
		latest, err := scanPosts(rows)
		rows.Close()
		if err != nil {
			return nil, err
		}
		for _, p := range latest {
			if len(result) >= limit {
				break
			}
			if !seen[p.ID] {
				result = append(result, p)
				seen[p.ID] = true
			}
		}
	}
	return result, nil
}

// ListPublished returns published posts (newest first) + total count.
// An empty category returns posts across all categories.
func (r *PostRepository) ListPublished(ctx context.Context, page, limit int, category string) ([]models.Post, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM posts
		WHERE status = 'published' AND ($1 = '' OR category = $1)`, category).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	rows, err := r.db.Query(ctx, `SELECT `+postCols+`
		FROM posts WHERE status = 'published' AND ($1 = '' OR category = $1)
		ORDER BY published_at DESC NULLS LAST, created_at DESC
		LIMIT $2 OFFSET $3`, category, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	posts, err := scanPosts(rows)
	return posts, total, err
}

// ListCategories returns the distinct categories that have at least one published post.
func (r *PostRepository) ListCategories(ctx context.Context) ([]string, error) {
	rows, err := r.db.Query(ctx, `SELECT DISTINCT category FROM posts
		WHERE status = 'published' AND category <> '' ORDER BY category`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []string{}
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}

// GetPublishedBySlug returns a single published post by slug, or nil.
func (r *PostRepository) GetPublishedBySlug(ctx context.Context, slug string) (*models.Post, error) {
	var p models.Post
	err := scanPost(r.db.QueryRow(ctx, `SELECT `+postCols+` FROM posts WHERE slug = $1 AND status = 'published'`, slug), &p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// List returns all posts (admin), newest first.
func (r *PostRepository) List(ctx context.Context) ([]models.Post, error) {
	rows, err := r.db.Query(ctx, `SELECT `+postCols+` FROM posts ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPosts(rows)
}

func (r *PostRepository) GetByID(ctx context.Context, id int) (*models.Post, error) {
	var p models.Post
	err := scanPost(r.db.QueryRow(ctx, `SELECT `+postCols+` FROM posts WHERE id = $1`, id), &p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *PostRepository) Create(ctx context.Context, p *models.Post) (*models.Post, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO posts (title, slug, excerpt, content, cover_image, category, tags, status, published_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id, created_at, updated_at`,
		p.Title, p.Slug, p.Excerpt, p.Content, p.CoverImage, p.Category, p.Tags, p.Status, p.PublishedAt,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PostRepository) Update(ctx context.Context, p *models.Post) (*models.Post, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE posts SET title=$1, slug=$2, excerpt=$3, content=$4, cover_image=$5,
			category=$6, tags=$7, status=$8, published_at=$9, updated_at=NOW()
		WHERE id=$10
		RETURNING created_at, updated_at`,
		p.Title, p.Slug, p.Excerpt, p.Content, p.CoverImage, p.Category, p.Tags, p.Status, p.PublishedAt, p.ID,
	).Scan(&p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PostRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM posts WHERE id = $1`, id)
	return err
}
