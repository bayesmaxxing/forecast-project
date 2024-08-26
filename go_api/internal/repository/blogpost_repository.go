package repository

import (
	"context"
	"database/sql"
	"go_api/internal/models"
)

type BlogpostRepository struct {
	db *sql.DB
}

func NewBlogpostRepository(db *sql.DB) *BlogpostRepository {
	return &BlogpostRepository{db: db}
}

func (r *BlogpostRepository) GetBlogposts(ctx context.Context) ([]*models.Blogpost, error) {
	query := `SELECT post_id, title, post, created, slug, related_forecasts
				FROM blogposts`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var blogposts []*models.Blogpost
	for rows.Next() {
		var b models.Blogpost
		if err := rows.Scan(&b.ID, &b.Title, &b.Post, &b.CreatedAt, &b.Slug, &b.RelatedForecasts); err != nil {
			return nil, err
		}
		blogposts = append(blogposts, &b)
	}
	return blogposts, rows.Err()
}

func (r *BlogpostRepository) GetBlogpostBySlug(ctx context.Context, slug string) (*models.Blogpost, error) {
	query := `SELECT post_id, title, post, created, slug, related_forecasts
				FROM blogposts 
				WHERE slug like (%$1%)`

	var b models.Blogpost
	err := r.db.QueryRowContext(ctx, query, slug).Scan(&b.ID, &b.Title, &b.Post, &b.Slug, &b.CreatedAt, &b.RelatedForecasts)

	if err != nil {
		return nil, err
	}
	return &b, nil
}