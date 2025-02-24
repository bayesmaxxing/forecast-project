package repository

import (
	"backend/internal/database"
	"backend/internal/models"
	"context"
	"time"
)

// BlogpostRepository defines the interface for blogpost data operations
type BlogpostRepository interface {
	GetBlogposts(ctx context.Context) ([]*models.Blogpost, error)
	GetBlogpostBySlug(ctx context.Context, slug string) (*models.Blogpost, error)
	CreateBlogpost(ctx context.Context, post *models.Blogpost) error
}

// PostgresBlogpostRepository implements the BlogpostRepository interface
type PostgresBlogpostRepository struct {
	db *database.DB
}

// NewBlogpostRepository creates a new PostgresBlogpostRepository instance
func NewBlogpostRepository(db *database.DB) BlogpostRepository {
	return &PostgresBlogpostRepository{db: db}
}

func (r *PostgresBlogpostRepository) GetBlogposts(ctx context.Context) ([]*models.Blogpost, error) {
	query := `SELECT 
				post_id
				, title
				, post
				, created
				, slug
				FROM blogposts`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var blogposts []*models.Blogpost
	for rows.Next() {
		var b models.Blogpost
		if err := rows.Scan(&b.ID, &b.Title, &b.Post, &b.CreatedAt, &b.Slug); err != nil {
			return nil, err
		}
		blogposts = append(blogposts, &b)
	}
	return blogposts, rows.Err()
}

func (r *PostgresBlogpostRepository) GetBlogpostBySlug(ctx context.Context, slug string) (*models.Blogpost, error) {
	query := `SELECT 
				post_id
				, title
				, post
				, created
				, slug
				FROM blogposts 
				WHERE slug like $1`
	slugPattern := "%" + slug + "%"
	var b models.Blogpost
	err := r.db.QueryRowContext(ctx, query, slugPattern).Scan(&b.ID, &b.Title, &b.Post, &b.CreatedAt, &b.Slug)

	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *PostgresBlogpostRepository) CreateBlogpost(ctx context.Context, post *models.Blogpost) error {
	post.CreatedAt = time.Now()

	query := `INSERT INTO blogposts (title, post, created_date, summary, slug)
				VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, post.Title, post.Post, post.CreatedAt, post.Summary, post.Slug).Scan(&post.ID)
	return err
}
