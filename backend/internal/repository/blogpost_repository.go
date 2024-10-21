package repository

import (
	"context"
	"backend/internal/database"
	"backend/internal/models"
	"time"
)

type BlogpostRepository struct {
	db *database.DB
}

func NewBlogpostRepository(db *database.DB) *BlogpostRepository {
	return &BlogpostRepository{db: db}
}

func (r *BlogpostRepository) GetBlogposts(ctx context.Context) ([]*models.Blogpost, error) {
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

func (r *BlogpostRepository) GetBlogpostBySlug(ctx context.Context, slug string) (*models.Blogpost, error) {
	query := `SELECT 
				post_id
				, title
				, post
				, created
				, slug
				FROM blogposts 
				WHERE slug like (%$1%)`

	var b models.Blogpost
	err := r.db.QueryRowContext(ctx, query, slug).Scan(&b.ID, &b.Title, &b.Post, &b.Slug, &b.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BlogpostRepository) CreateBlogpost(ctx context.Context, post *models.Blogpost) error {
	post.CreatedAt = time.Now()

	query := `INSERT INTO blogposts (title, post, created_date, summary, slug)
				VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, post.Title, post.Post, post.CreatedAt, post.Summary, post.Slug).Scan(&post.ID)
	return err
}
